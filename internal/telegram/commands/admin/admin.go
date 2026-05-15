package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	userModes "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func GetAllUsersHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const chunkSize = 50
		offset := 0

		for {
			users, total, err := app.UserService.GetAllUsersPaginated(ctx, chunkSize, offset)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Erro ao buscar usuários.",
				})
				return
			}

			if len(users) == 0 {
				if offset == 0 {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Nenhum usuário encontrado.",
					})
				}
				break
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("👥 Total de Usuários: <b>%d</b>\n<blockquote>Página %d</blockquote>\n",
				total, (offset/chunkSize)+1))

			for _, u := range users {
				sb.WriteString(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> - %d\n", u.UserId, u.FirstName, u.UserId))
			}

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      sb.String(),
				ParseMode: models.ParseModeHTML,
			})
			if err != nil {
				break
			}

			offset += chunkSize
			if int64(offset) >= total {
				break
			}
		}
	}
}

func GetAllChannelsHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const chunkSize = 50
		offset := 0
		val := true

		for {
			channels, total, err := app.ChannelService.GetAllChannelsPaginated(ctx, chunkSize, offset)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Erro ao buscar canais.",
				})
				return
			}

			if len(channels) == 0 {
				if offset == 0 {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Nenhum canal encontrado.",
					})
				}
				break
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("📦 Total de Canais: <b>%d</b>\n<blockquote>Página %d</blockquote>\n",
				total, (offset/chunkSize)+1))

			for _, c := range channels {
				sb.WriteString(fmt.Sprintf(`<a href='%s'>%s</a> - <code>%d</code>`+"\n", c.InviteURL, c.Title, c.ID))
			}

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      sb.String(),
				ParseMode: models.ParseModeHTML,
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: &val,
				},
			})
			if err != nil {
				break
			}

			offset += chunkSize
			if int64(offset) >= total {
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

		channel, err := app.ChannelService.GetChannelByID(ctx, channelID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return
		}

		owner, err := app.UserService.GetUserByID(ctx, channel.OwnerID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Dono não encontrado!",
			})
			return
		}

		ownerID := fmt.Sprintf("%d", config.OwnerID)
		msg := fmt.Sprintf(
			"ID: <code>%d</code>\nCanal: %s\nLink: %s\nDono: <a href='tg://user?id=%d'>%s</a> (<code>%d</code>)\nPainel: %s",
			channel.ID,
			html.EscapeString(channel.Title),
			channel.InviteURL,
			owner.UserId,
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

		channel, err := app.ChannelService.GetChannelByID(ctx, channelID)
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

		err = app.ChannelService.UpdateOwnerChannel(ctx, channelID, channel.OwnerID, newOwnerID)
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

		user, err := app.UserService.GetUserByID(ctx, userID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Usuário não encontrado!: %v", err),
			})
			return
		}

		channels, _ := app.ChannelService.GetUserChannels(ctx, user.UserId)
		header := fmt.Sprintf("👤 <b><a href='tg://user?id=%d'>%s</a></b> (<code>%d</code>)\n📦 Canais: <b>%d</b>\n\n",
			user.UserId,
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
				lines = append(lines, fmt.Sprintf("<a href='%s'>%s</a> - <code>%d</code>",
					c.InviteURL,
					html.EscapeString(c.Title),
					c.ID,
				))
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

		channel, err := app.ChannelService.GetChannelByID(ctx, channelID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return
		}

		if err = app.ChannelService.DisconnectChannel(ctx, channel.OwnerID, channelID); err != nil {

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

		users, err := app.UserService.GetAllUsers(ctx)
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
				logger.Error("ADMIN", "Erro ao enviar aviso para %d - %s: %v", user.UserId, user.FirstName, err)
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
		user, _ := b.GetMe(ctx)

		data := map[string]string{
			"botUsername": user.Username,
		}

		text, button := parser.GetMessage("publi", data)

		channels, err := app.ChannelService.GetAllChannels(ctx)
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
				ChatID:      ch.ID,
				Text:        text,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: button,
			})
			if err != nil {
				logger.Error("ADMIN", "Erro ao enviar aviso para canal %d - %s: %v", ch.ID, ch.Title, err)
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
			logger.Error("ADMIN", "Erro ao enviar mensagem para %d: %v", targetID, err)
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
		existingChannel, _ := c.ChannelService.GetChannelByID(ctx, channelID)
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
			logger.Error("ADMIN", "Erro ao buscar canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "❌ Erro ao buscar informações do canal."})
			return
		}

		ownerInfo, err := b.GetChat(ctx, &bot.GetChatParams{ChatID: ownerID})
		if err != nil {
			logger.Error("ADMIN", "Erro ao buscar usuário: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "❌ Erro ao buscar informações do usuário."})
			return
		}

		// Cria usuário caso não exista
		_ = c.UserService.UpsertUser(ctx, &userModes.User{
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
		channel, err := c.ChannelService.CreateChannelWithDefaults(ctx, channelID, channelInfo.Title, inviteURL, newPackCaption, defaultCaption, ownerID)
		if err != nil {
			logger.Error("ADMIN", "Erro ao criar canal: %v", err)
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

func NoticeUsersReplyHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		if update.Message.ReplyToMessage == nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Responda a uma mensagem para enviar o aviso aos usuários.",
			})
			return
		}

		noticeText := update.Message.ReplyToMessage.Text
		if noticeText == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ A mensagem respondida está vazia.",
			})
			return
		}

		users, err := app.UserService.GetAllUsers(ctx)
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
				logger.Error("ADMIN", "Erro ao enviar aviso para %d - %s: %v", user.UserId, user.FirstName, err)
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

func NoticeChannelsReplyHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		// Verifica se é uma resposta
		if update.Message.ReplyToMessage == nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Responda a uma mensagem para enviar o aviso aos canais.",
			})
			return
		}

		// Pega o texto da mensagem respondida
		noticeText := update.Message.ReplyToMessage.Text
		if noticeText == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ A mensagem respondida está vazia.",
			})
			return
		}

		channels, err := app.ChannelService.GetAllChannels(ctx)
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
				logger.Error("ADMIN", "Erro ao enviar aviso para canal %d - %s: %v", ch.ID, ch.Title, err)
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

func ToggleMaintenceHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		maintenance, err := app.ServerService.ToggleMaintenance(ctx)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Erro ao alterar o modo de manutenção.",
			})
			return
		}

		var msg string

		if maintenance {
			msg = "⚠️ <b>Modo de manutenção ativado</b>\n\nO bot pode ficar temporariamente indisponível."
		} else {
			msg = "✅ <b>Modo de manutenção desativado</b>\n\nO bot voltou a funcionar normalmente."
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func SetAdminHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, upt *models.Update) {
		if upt.Message == nil {
			return
		}

		args := strings.Fields(upt.Message.Text)

		if len(args) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    upt.Message.Chat.ID,
				Text:      "❌ Uso correto:\n<code>/setadmin [userID]</code>",
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upt.Message.Chat.ID,
				Text:   "❌ userID inválido.",
			})
			return
		}

		isAdmin, err := app.UserService.UpdateUserAdmin(ctx, userID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upt.Message.Chat.ID,
				Text:   "❌ Erro ao alterar status de admin.",
			})
			return
		}

		var msg string

		if isAdmin {
			msg = fmt.Sprintf("✅ Usuário <code>%d</code> agora é administrador.", userID)
		} else {
			msg = fmt.Sprintf("⚠️ Usuário <code>%d</code> não é mais administrador.", userID)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    upt.Message.Chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
		})
	}
}

func LogRemoji(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		sla, _ := json.Marshal(update)
		logger.Info("DEBUG", "Update payload: %s", string(sla))
	}
}

func CheckBotAdminHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const targetBotID = 5986082367
		const targetBotUser = "@XavolaBot"

		channels, err := app.ChannelService.GetAllChannels(ctx)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Erro ao buscar canais do banco.",
			})
			return
		}

		statusMsg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("🔍 Verificando <b>%d</b> canais... Isso pode levar um tempo.", len(channels)),
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})

		var foundChannels []string
		var mu sync.Mutex
		count := 0

		// Canal para processar os canais em paralelo
		type result struct {
			info string
			err  error
		}

		chQueue := make(chan *userModes.Channel, len(channels))
		for i := range channels {
			chQueue <- &channels[i]
		}
		close(chQueue)

		// Usar um WaitGroup para gerenciar as goroutines (limite de 10 concorrentes)
		var wg sync.WaitGroup
		numWorkers := 10
		if len(channels) < numWorkers {
			numWorkers = len(channels)
		}

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ch := range chQueue {
					member, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
						ChatID: ch.ID,
						UserID: targetBotID,
					})

					if err == nil && (member.Type == "administrator" || member.Type == "creator") {
						owner, _ := app.UserService.GetUserByID(ctx, ch.OwnerID)
						ownerName := "Desconhecido"
						if owner != nil {
							ownerName = owner.FirstName
						}

						info := fmt.Sprintf(
							"<b>Canal:</b> %s\n<b>ID:</b> <code>%d</code>\n<b>Link:</b> %s\n<b>Dono:</b> <a href='tg://user?id=%d'>%s</a> (<code>%d</code>)\n",
							html.EscapeString(ch.Title),
							ch.ID,
							ch.InviteURL,
							ch.OwnerID,
							html.EscapeString(ownerName),
							ch.OwnerID,
						)

						mu.Lock()
						foundChannels = append(foundChannels, info)
						count++
						mu.Unlock()
					}
				}
			}()
		}
		wg.Wait()

		if count == 0 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: statusMsg.ID,
				Text:      fmt.Sprintf("✅ O bot %s não foi encontrado como admin em nenhum canal.", targetBotUser),
			})
			return
		}

		header := fmt.Sprintf("🤖 <b>Bot %s encontrado em %d canais:</b>\n\n", targetBotUser, count)

		// Enviar em blocos para evitar limite de caracteres do Telegram
		var sb strings.Builder
		sb.WriteString(header)
		first := true
		for i, info := range foundChannels {
			if sb.Len()+len(info) > 3800 {
				if first {
					b.EditMessageText(ctx, &bot.EditMessageTextParams{
						ChatID:             update.Message.Chat.ID,
						MessageID:          statusMsg.ID,
						Text:               sb.String(),
						ParseMode:          models.ParseModeHTML,
						LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: bot.True()},
					})
					first = false
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:             update.Message.Chat.ID,
						Text:               sb.String(),
						ParseMode:          models.ParseModeHTML,
						LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: bot.True()},
					})
				}
				sb.Reset()
			}
			sb.WriteString(info + "────────────────────\n")
			if i == len(foundChannels)-1 {
				if first {
					b.EditMessageText(ctx, &bot.EditMessageTextParams{
						ChatID:             update.Message.Chat.ID,
						MessageID:          statusMsg.ID,
						Text:               sb.String(),
						ParseMode:          models.ParseModeHTML,
						LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: bot.True()},
					})
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:             update.Message.Chat.ID,
						Text:               sb.String(),
						ParseMode:          models.ParseModeHTML,
						LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: bot.True()},
					})
				}
			}
		}
	}
}

func AdminHelpHandler(app *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		msg := `👨‍💻 <b>Painel de Administração</b>

<b>Comandos de Listagem:</b>
/users - Lista todos os usuários
/channels - Lista todos os canais
/user [id] - Informações detalhadas do usuário
/info [id] - Informações detalhadas do canal

<b>Comandos de Mensagem:</b>
/notice [msg] - Envia aviso para todos os usuários (primeira linha é comando)
/publi - Envia mensagem de publicidade padrão para todos os canais
/send [id]\n[msg] - Envia mensagem privada para um ID específico
/allusers (reply) - Envia a mensagem respondida para todos os usuários
/allchannels (reply) - Envia a mensagem respondida para todos os canais

<b>Comandos de Gerenciamento:</b>
/add [canalID] [donoID] - Adiciona canal e dono manualmente
/remove [id] - Remove um canal do sistema
/transfer [canalID] [novoDonoID] - Transfere posse de um canal
/setadmin [id] - Alterna status de administrador de um usuário
/maintence - Ativa/Desativa modo de manutenção
/backup - Gera backup do banco de dados

<b>Utilidades:</b>
/checkbot - Verifica se o XavolaBot é admin nos canais
/emoji - Log de update (debug)`

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "📊 Abrir Dashboard Admin",
						WebApp: &models.WebAppInfo{
							URL: fmt.Sprintf("%s/admin/dash", config.WebAppURL),
						},
					},
				},
			},
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        msg,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
	}
}
