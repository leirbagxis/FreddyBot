// file: internal/telegram/events/channelPost/newpack.go (sugestão)
package channelpost

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// Estado por canal: quando recebeu !newpack e está aguardando sticker
type newPackState struct {
	WaitingForSticker bool
	MessageID         int // id da mensagem do !newpack (pra editar depois)
}

var newPackStates sync.Map // map[int64]newPackState

var cmdNewPackRegex = regexp.MustCompile(`^(?:/newpack|!newpack)(\s|$)`)

// -------------------- template --------------------
func renderNewPackTemplate(tpl, titulo, link string) string {
	out := tpl
	out = strings.ReplaceAll(out, "$titulo", titulo)
	out = strings.ReplaceAll(out, "$link", link)

	out = strings.ReplaceAll(out, "$title", titulo)
	out = strings.ReplaceAll(out, "$url", link)
	return out
}

func shouldDisableLinkPreview(channel dbmodels.Channel) bool {
	// padrão do seu código é "permitir" quando não há config explícita. [file:7]
	if channel.DefaultCaption == nil || channel.DefaultCaption.MessagePermission == nil {
		return false
	}
	return !channel.DefaultCaption.MessagePermission.LinkPreview
}

func (mp MessageProcessor) TryHandleNewPack(ctx context.Context, channel dbmodels.Channel, post models.Message) (handled bool, err error) {
	channelID := post.Chat.ID

	// 1) Se recebeu comando !newpack
	if post.Text != "" && cmdNewPackRegex.MatchString(strings.TrimSpace(post.Text)) {
		// Marca modo ativo (reaproveita seu mecanismo)
		mp.mediaGroupManager.SetNewPackActive(channelID, true) // [file:3]
		newPackStates.Store(channelID, newPackState{
			WaitingForSticker: true,
			MessageID:         post.ID,
		})

		editCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
		defer cancel()

		// Similar ao JS: "Envie-me um sticker do seu pack..."
		_, err := mp.bot.EditMessageText(editCtx, &bot.EditMessageTextParams{
			ChatID:    channelID,
			MessageID: post.ID,
			Text:      "Envie-me um sticker do seu pack...",
		})
		if err != nil {
			// Se falhar, ainda assim mantém o estado? Aqui preferi desativar pra não travar.
			newPackStates.Delete(channelID)
			mp.mediaGroupManager.SetNewPackActive(channelID, false)
			return true, err
		}
		return true, nil
	}

	// 2) Se chegou sticker e o canal está aguardando
	if post.Sticker != nil {
		v, ok := newPackStates.Load(channelID)
		if !ok {
			return false, nil
		}
		state := v.(newPackState)
		if !state.WaitingForSticker {
			return false, nil
		}

		// Validar se sticker tem set_name (sticker de pack público)
		setName := post.Sticker.SetName
		if strings.TrimSpace(setName) == "" {
			sendCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()

			_, _ = mp.bot.SendMessage(sendCtx, &bot.SendMessageParams{
				ChatID: channelID,
				Text:   "Sticker não faz parte de um pack público.",
			})
			newPackStates.Delete(channelID)
			mp.mediaGroupManager.SetNewPackActive(channelID, false)
			return true, nil
		}

		// Buscar infos do pack
		getCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		pack, err := mp.bot.GetStickerSet(getCtx, &bot.GetStickerSetParams{
			Name: setName,
		})
		if err != nil {
			newPackStates.Delete(channelID)
			mp.mediaGroupManager.SetNewPackActive(channelID, false)
			return true, err
		}

		packURL := fmt.Sprintf("https://t.me/addstickers/%s", setName)

		title := strings.TrimSpace(pack.Title)
		if title == "" {
			title = "Meu Pack"
		}
		link := fmt.Sprintf("https://t.me/addstickers/%s", setName)

		tpl := channel.NewPackCaption
		if strings.TrimSpace(tpl) == "" {
			tpl = "$titulo\n$link"
		}

		caption := renderNewPackTemplate(tpl, title, link)

		// aplicar link preview permission (como processors.go faz no texto). [file:7]
		disableLP := shouldDisableLinkPreview(channel)

		editCtx, cancel2 := context.WithTimeout(ctx, 10*time.Second)
		defer cancel2()

		// Edita o texto da mensagem original do comando
		editParams := bot.EditMessageTextParams{
			ChatID:    channelID,
			MessageID: state.MessageID,
			Text:      detectParseMode(caption),
			ParseMode: models.ParseModeHTML,
		}

		if disableLP {
			val := true
			editParams.LinkPreviewOptions = &models.LinkPreviewOptions{
				IsDisabled: &val,
			}
		}

		_, err = mp.bot.EditMessageText(editCtx, &editParams)

		if err != nil {
			newPackStates.Delete(channelID)
			mp.mediaGroupManager.SetNewPackActive(channelID, false)
			return true, err
		}

		// Adiciona botão com link do pack (equivalente ao inline keyboard do JS)
		kb := models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: title, URL: packURL},
				},
			},
		}

		_, err = mp.bot.EditMessageReplyMarkup(editCtx, &bot.EditMessageReplyMarkupParams{
			ChatID:      channelID,
			MessageID:   state.MessageID,
			ReplyMarkup: kb,
		})
		if err != nil {
			// Não precisa falhar o fluxo inteiro; o texto já foi editado.
			log.Printf("falha ao editar reply markup newpack: %v", err)
		}

		_, err = mp.bot.EditMessageReplyMarkup(editCtx, &bot.EditMessageReplyMarkupParams{
			ChatID:      channelID,
			MessageID:   post.ID,
			ReplyMarkup: kb,
		})

		if err != nil {
			// Não precisa falhar o fluxo inteiro; o texto já foi editado.
			log.Printf("falha ao editar reply markup newpack: %v", err)
		}

		// (Opcional) enviar sticker separador, semelhante ao JS
		if channel.Separator != nil && channel.Separator.SeparatorID != "" {
			sepCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			_, _ = mp.bot.SendSticker(sepCtx, &bot.SendStickerParams{
				ChatID:  channelID,
				Sticker: &models.InputFileString{channel.Separator.SeparatorID},
			})
		}

		// Finaliza estado
		newPackStates.Delete(channelID)
		mp.mediaGroupManager.SetNewPackActive(channelID, false)
		return true, nil
	}

	return false, nil
}

// helper mínimo (ou reaproveite seu utils.RemoveHTMLTags/escape já existente)
func escapeHTML(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;")
	return r.Replace(s)
}
