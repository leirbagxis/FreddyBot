package admin

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	userModes "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func GetAllUsersHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		users, err := app.UserRepo.GetAllUSers(ctx)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erro ao buscar usuários.",
			})
			return
		}
		const chunkSize = 50
		total := len(users)

		for i := 0; i < total; i += chunkSize {
			end := i + chunkSize
			if end > total {
				end = total
			}

			chunk := users[i:end]
			msg := fmt.Sprintf("👥 Total de Usuários: <b>%d</b>\n<blockquote>Página %d</blockquote>\n",
				total, (i/chunkSize)+1)

			for _, u := range chunk {
				msg += fmt.Sprintf("<i>%d - %s</i>\n", u.UserId, u.FirstName)
			}

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      msg,
				ParseMode: models.ParseModeHTML,
			})

			if err != nil {
				break
			}
		}

	}
}

func GetAllChannelsHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		channels, err := app.ChannelRepo.GetAllChannels(ctx)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erro ao buscar canais.",
			})
			return
		}
		const chunkSize = 50
		total := len(channels)
		val := true

		for i := 0; i < total; i += chunkSize {
			end := i + chunkSize
			if end > total {
				end = total
			}

			chunk := channels[i:end]
			msg := fmt.Sprintf("📦 Total de Canais: <b>%d</b>\n<blockquote>Página %d</blockquote>\n",
				total, (i/chunkSize)+1)

			for _, c := range chunk {
				userID := fmt.Sprintf("%d", c.OwnerID)
				channelID := fmt.Sprintf("%d", c.ID)
				link := auth.GenerateMiniAppUrl(userID, channelID)
				msg += fmt.Sprintf(`<i>%d -</i> <a href="%s">%s</a>`+"\n", c.ID, link, c.Title)
			}

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      msg,
				ParseMode: models.ParseModeHTML,
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: &val,
				},
			})

			if err != nil {
				break
			}
		}
	}
}

func GetBackUpHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		dbPath := config.DatabaseFile
		fileData, err := os.ReadFile(dbPath)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Erro ao ler banco: %v", err),
			})
			return
		}

		timestamp := time.Now().Format("2006-01-02-15-04-05")
		fileName := fmt.Sprintf("caption-backup-%s.db", strings.ReplaceAll(timestamp, ":", "-"))

		params := &bot.SendDocumentParams{
			ChatID: update.Message.Chat.ID,
			Document: &models.InputFileUpload{
				Filename: fileName,
				Data:     bytes.NewReader(fileData),
			},
			Caption: fmt.Sprintf("🗂️ Backup gerado em %s", timestamp),
		}

		_, err = b.SendDocument(ctx, params)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Erro ao enviar backup: %v", err),
			})
		}
	}
}

func GetInfoChannelHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		channelIDStr := strings.TrimSpace(update.Message.Text[len("/info"):])
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return
		}

		channel, err := app.ChannelRepo.GetChannelByID(ctx, channelID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return
		}

		owner, err := app.UserRepo.GetUserById(ctx, channel.OwnerID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Dono não encontrado!",
			})
			return
		}

		ownerID := fmt.Sprintf("%d", config.OwnerID)
		msg := fmt.Sprintf(
			"ID: <code>%d</code>\nCanal: %s\nLink: %s\nDono: %s (%d)\nPainel: %s",
			channel.ID,
			html.EscapeString(channel.Title),
			channel.InviteURL,
			html.EscapeString(owner.FirstName),
			owner.UserId,
			auth.GenerateMiniAppUrl(ownerID, channelIDStr),
		)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
		})
	}
}

func RegisterTransferHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		input := strings.TrimSpace(update.Message.Text[len("/transfer"):])
		parts := strings.Fields(input)
		if len(parts) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Uso: /transfer <channelId> <newOwnerId>",
			})
			return
		}

		channelID, _ := strconv.ParseInt(parts[0], 10, 64)
		newOwnerID, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := app.ChannelRepo.GetChannelByID(ctx, channelID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return
		}

		tgUser, err := b.GetChat(ctx, &bot.GetChatParams{ChatID: newOwnerID})
		if err != nil || tgUser.FirstName == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ ID de usuário inválido: %d", newOwnerID),
			})
			return
		}

		err = app.ChannelRepo.UpdateOwnerChannel(ctx, channelID, channel.OwnerID, newOwnerID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Erro ao transferir canal",
			})
			return
		}

		msg := fmt.Sprintf(
			"✅ <b>Transferência realizada com sucesso!</b>\n<b>Canal:</b> %s\n<b>ID:</b> %d\n<b>Novo Dono:</b> %s (%d)\n\n🔗 <a href=\"%s\">Abrir painel do canal</a>",
			html.EscapeString(channel.Title),
			channelID,
			html.EscapeString(tgUser.FirstName),
			newOwnerID,
			auth.GenerateMiniAppUrl(parts[1], parts[0]),
		)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
		})

	}
}

func GetInfoUserHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userIDStr := strings.TrimSpace(update.Message.Text[len("/user"):])
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return
		}

		user, err := app.UserRepo.GetUserById(ctx, userID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Usuário não encontrado!: %v", err),
			})
			return
		}

		channels, _ := app.ChannelRepo.GetAllChannelsByUserID(ctx, user.UserId)
		header := fmt.Sprintf("👤 <b>%s</b> (<code>%d</code>)\n📦 Canais: <b>%d</b>\n\n",
			html.EscapeString(user.FirstName),
			user.UserId,
			len(channels),
		)

		if len(channels) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      header + "Usuário ainda não possui canais.",
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		const chunkSize = 20
		for i := 0; i < len(channels); i += chunkSize {
			chunk := channels[i:min(i+chunkSize, len(channels))]
			var lines []string
			for _, c := range chunk {
				lines = append(lines, fmt.Sprintf("<b>%d</b> - %s", c.ID, html.EscapeString(c.Title)))
			}

			msg := header + strings.Join(lines, "\n")
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      msg,
				ParseMode: models.ParseModeHTML,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
		}
	}
}

func RemoveChannelHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		channelIDStr := strings.TrimSpace(update.Message.Text[len("/remove"):])
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return
		}

		channel, err := app.ChannelRepo.GetChannelByID(ctx, channelID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return
		}

		if err = app.ChannelRepo.DeleteChannelWithRelations(ctx, channel.OwnerID, channelID); err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Não foi possivel deletar o canal: %v", err),
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "✅ Canal excluído com sucesso!",
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func NoticeCommandHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		text := update.Message.Text
		if text == "" {
			return
		}

		lines := strings.Split(text, "\n")
		if len(lines) <= 1 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ A mensagem de aviso está vazia.",
			})
			return
		}

		noticeText := strings.TrimSpace(strings.Join(lines[1:], "\n"))

		users, err := app.UserRepo.GetAllUSers(ctx)
		if err != nil || len(users) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Nenhum usuário encontrado.",
			})
			return
		}

		var failedUsers []string
		sentCount := 0

		for _, user := range users {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    user.UserId,
				Text:      noticeText,
				ParseMode: models.ParseModeHTML,
			})
			if err != nil {
				log.Printf("Erro ao enviar aviso para %d - %s: %v", user.UserId, user.FirstName, err)
				failedUsers = append(failedUsers, fmt.Sprintf("<code>%d</code> - %s", user.UserId, html.EscapeString(user.FirstName)))
			} else {
				sentCount++
			}
		}

		var finalMsg strings.Builder
		finalMsg.WriteString(fmt.Sprintf("📨 Aviso enviado para <b>%d</b> usuários.\n", sentCount))

		if len(failedUsers) > 0 {
			finalMsg.WriteString(fmt.Sprintf("\n❌ Falhou para %d usuários:\n", len(failedUsers)))
			finalMsg.WriteString(strings.Join(failedUsers, "\n"))
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      finalMsg.String(),
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func NoticeChannelsHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const noticeText = `📢 <b>Atualização Importante</b>

Estamos realizando melhorias no painel de gerenciamento de legendas. Caso encontre algum erro ou comportamento inesperado, entre em contato com o suporte.

Obrigado por usar o FreddyBot!`

		channels, err := app.ChannelRepo.GetAllChannels(ctx)
		if err != nil || len(channels) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Nenhum canal encontrado.",
			})
			return
		}

		var failedChannels []string
		sentCount := 0

		for _, ch := range channels {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    ch.ID,
				Text:      noticeText,
				ParseMode: models.ParseModeHTML,
			})
			if err != nil {
				log.Printf("Erro ao enviar aviso para canal %d - %s: %v", ch.ID, ch.Title, err)
				failedChannels = append(failedChannels, fmt.Sprintf("<code>%d</code> - %s", ch.ID, html.EscapeString(ch.Title)))
			} else {
				sentCount++
			}
		}

		var resultMsg strings.Builder
		resultMsg.WriteString(fmt.Sprintf("📨 Aviso enviado para <b>%d</b> canais.\n", sentCount))

		if len(failedChannels) > 0 {
			resultMsg.WriteString(fmt.Sprintf("\n❌ Falhou para %d canais:\n", len(failedChannels)))
			resultMsg.WriteString(strings.Join(failedChannels, "\n"))
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      resultMsg.String(),
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func SendMessageToIdHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.Text == "" {
			return
		}

		lines := strings.Split(update.Message.Text, "\n")
		if len(lines) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Uso inválido. Envie no formato:\n/send <id>\n<mensagem>",
			})
			return
		}

		idStr := strings.TrimSpace(strings.TrimPrefix(lines[0], "/send"))
		targetID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return
		}

		message := strings.Join(lines[1:], "\n")
		message = strings.TrimSpace(message)

		if message == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Mensagem vazia.",
			})
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    targetID,
			Text:      message,
			ParseMode: models.ParseModeHTML,
		})
		if err != nil {
			log.Printf("Erro ao enviar mensagem para %d: %v", targetID, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      fmt.Sprintf("❌ Falha ao enviar para <code>%d</code>: %v", targetID, err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("✅ Mensagem enviada para <code>%d</code> com sucesso.", targetID),
			ParseMode: models.ParseModeHTML,
		})
	}
}

func AddChannelCommandHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		botInfo, _ := b.GetMe(ctx)

		msgText := strings.TrimSpace(update.Message.Text)
		args := strings.SplitN(msgText, " ", 3)
		if len(args) < 3 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Uso correto: /add <channel_id> <owner_id>",
			})
			return
		}

		channelIDStr := args[1]
		ownerIDStr := args[2]
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		ownerID, err2 := strconv.ParseInt(ownerIDStr, 10, 64)
		if err != nil || err2 != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ IDs inválidos. Certifique-se de que ambos são numéricos.",
			})
			return
		}

		// Verifica se canal já existe
		existingChannel, _ := c.ChannelRepo.GetChannelByID(ctx, channelID)
		fmt.Println(existingChannel)
		if existingChannel != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Canal já existe no banco de dados.",
			})
			return
		}

		// Pega informações do canal e do dono
		channelInfo, err := b.GetChat(ctx, &bot.GetChatParams{ChatID: channelID})
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "❌ Erro ao buscar informações do canal."})
			return
		}

		ownerInfo, err := b.GetChat(ctx, &bot.GetChatParams{ChatID: ownerID})
		if err != nil {
			log.Printf("Erro ao buscar usuário: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "❌ Erro ao buscar informações do usuário."})
			return
		}

		// Cria usuário caso não exista
		_ = c.UserRepo.UpsertUser(ctx, &userModes.User{
			UserId:    ownerID,
			FirstName: utils.RemoveHTMLTags(ownerInfo.FirstName),
		})

		// Gera caption
		newPackCaption := fmt.Sprintf(`╔═━──━═༻✧༺═━──━═╗

        𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
        (｡•́︿•̀｡)っ✧.*ೃ༄
        ˗ˏˋ [$name]($link) ⁠⋆｡˚ ☁︎
             彡♡ ₊˚

⋆｡˚ ❀ @%s ☽⁺₊

╚═━──━═༻✧༺═━──━═╝`, botInfo.Username)

		defaultCaption := fmt.Sprintf("➽ 𝐛𝐲 @%s", botInfo.Username)
		inviteURL := channelInfo.InviteLink
		if channelInfo.Username != "" {
			inviteURL = fmt.Sprintf("t.me/%s", channelInfo.Username)
		}

		// Cria canal
		channel, err := c.ChannelRepo.CreateChannelWithDefaults(ctx, channelID, channelInfo.Title, inviteURL, newPackCaption, defaultCaption, ownerID)
		if err != nil {
			log.Printf("Erro ao criar canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "❌ Erro ao salvar canal."})
			return
		}

		miniApp := auth.GenerateMiniAppUrl(fmt.Sprintf("%d", ownerID), fmt.Sprintf("%d", channelID))
		msg := fmt.Sprintf("✅ Canal salvo com sucesso - (%s - %d)\n\n%s", channel.Title, channel.ID, miniApp)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
		})
	}
}
