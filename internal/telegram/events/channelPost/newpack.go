package channelpost

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/mymmrac/telego"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// Estado por canal: quando recebeu !newpack e está aguardando sticker
type newPackState struct {
	WaitingForSticker bool
	MessageID         int
}

var newPackStates sync.Map // map[int64]newPackState
var cmdNewPackRegex = regexp.MustCompile(`^!newpack`)

func TryHandleNewPackTelego(ctx context.Context, b *telego.Bot, channel dbmodels.Channel, post telego.Message) (handled bool, err error) {
	channelID := post.Chat.ID
	mgm := GetMediaGroupManagerTelego()

	if post.Text != "" && cmdNewPackRegex.MatchString(strings.TrimSpace(post.Text)) {
		mgm.SetNewPackActive(channelID, true)
		newPackStates.Store(channelID, newPackState{
			WaitingForSticker: true,
			MessageID:         post.MessageID,
		})

		_, err := b.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: channelID},
			MessageID: post.MessageID,
			Text:      "Envie-me um sticker do seu pack...",
		})
		if err != nil {
			newPackStates.Delete(channelID)
			mgm.SetNewPackActive(channelID, false)
			return true, err
		}
		return true, nil
	}

	if post.Sticker != nil {
		v, ok := newPackStates.Load(channelID)
		if !ok {
			return false, nil
		}
		state := v.(newPackState)
		if !state.WaitingForSticker {
			return false, nil
		}

		setName := post.Sticker.SetName
		if strings.TrimSpace(setName) == "" {
			_, _ = b.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: channelID},
				Text:   "Sticker não faz parte de um pack público.",
			})
			newPackStates.Delete(channelID)
			mgm.SetNewPackActive(channelID, false)
			return true, nil
		}

		pack, err := b.GetStickerSet(context.Background(), &telego.GetStickerSetParams{
			Name: setName,
		})
		if err != nil {
			newPackStates.Delete(channelID)
			mgm.SetNewPackActive(channelID, false)
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
		
		pm := GetPermissionManager()
		perms := pm.CheckPermissions(&channel, MessageTypeText)
		disableLP := !perms.CanUseLinkPreview

		editParams := telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: channelID},
			MessageID: state.MessageID,
			Text:      DetectParseMode(caption),
			ParseMode: telego.ModeHTML,
		}

		if disableLP {
			editParams.LinkPreviewOptions = &telego.LinkPreviewOptions{
				IsDisabled: true,
			}
		}

		_, err = b.EditMessageText(context.Background(), &editParams)

		if err != nil {
			newPackStates.Delete(channelID)
			mgm.SetNewPackActive(channelID, false)
			return true, err
		}

		kb := telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: title, URL: packURL},
				},
			},
		}

		_, err = b.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
			ChatID:      telego.ChatID{ID: channelID},
			MessageID:   state.MessageID,
			ReplyMarkup: &kb,
		})
		if err != nil {
			logger.Error("BOT", "falha ao editar reply markup newpack: %v", err)
		}

		_, err = b.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
			ChatID:      telego.ChatID{ID: channelID},
			MessageID:   post.MessageID,
			ReplyMarkup: &kb,
		})

		if err != nil {
			logger.Error("BOT", "falha ao editar reply markup newpack: %v", err)
		}

		if channel.Separator != nil && channel.Separator.SeparatorID != "" {
			_, _ = b.SendSticker(context.Background(), &telego.SendStickerParams{
				ChatID:  telego.ChatID{ID: channelID},
				Sticker: telego.InputFile{FileID: channel.Separator.SeparatorID},
			})
		}

		newPackStates.Delete(channelID)
		mgm.SetNewPackActive(channelID, false)
		return true, nil
	}

	return false, nil
}

func renderNewPackTemplate(tpl, title, link string) string {
	res := strings.ReplaceAll(tpl, "$titulo", title)
	res = strings.ReplaceAll(res, "$name", title)
	res = strings.ReplaceAll(res, "$link", link)
	return res
}
