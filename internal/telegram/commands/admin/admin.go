package admin

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
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
		})
	}
}
