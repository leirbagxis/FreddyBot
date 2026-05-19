package admin

import (
	"context"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	userModes "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetAllUsersHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		const chunkSize = 50
		offset := 0
		bot := ctx.Bot()

		for {
			users, total, err := app.UserService.GetAllUsersPaginated(context.Background(), chunkSize, offset)
			if err != nil {
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID: update.Message.Chat.ChatID(),
					Text:   "Erro ao buscar usuários.",
				})
				return nil
			}

			if len(users) == 0 {
				if offset == 0 {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID: update.Message.Chat.ChatID(),
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

			_, err = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      sb.String(),
				ParseMode: telego.ModeHTML,
			})
			if err != nil {
				break
			}

			offset += chunkSize
			if int64(offset) >= total {
				break
			}
		}
		return nil
	}
}

func GetAllChannelsHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		const chunkSize = 50
		offset := 0
		bot := ctx.Bot()

		for {
			channels, total, err := app.ChannelService.GetAllChannelsPaginated(context.Background(), chunkSize, offset)
			if err != nil {
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID: update.Message.Chat.ChatID(),
					Text:   "Erro ao buscar canais.",
				})
				return nil
			}

			if len(channels) == 0 {
				if offset == 0 {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID: update.Message.Chat.ChatID(),
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

			_, err = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:             update.Message.Chat.ChatID(),
				Text:               sb.String(),
				ParseMode:          telego.ModeHTML,
				LinkPreviewOptions: &telego.LinkPreviewOptions{IsDisabled: true},
			})
			if err != nil {
				break
			}

			offset += chunkSize
			if int64(offset) >= total {
				break
			}
		}
		return nil
	}
}

func GetBackUpHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		dbPath := config.DatabaseFile
		
		file, err := os.Open(dbPath)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Erro ao abrir banco: %v", err),
			})
			return nil
		}
		defer file.Close()

		timestamp := time.Now().Format("2006-01-02-15-04-05")

		params := &telego.SendDocumentParams{
			ChatID: update.Message.Chat.ChatID(),
			Document: telego.InputFile{
				File: file,
			},
			Caption: fmt.Sprintf("🗂️ Backup gerado em %s", timestamp),
		}

		_, err = bot.SendDocument(context.Background(), params)
		if err != nil {
			logger.Error("ADMIN", "Erro ao enviar backup: %v", err)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Erro ao enviar backup: %v", err),
			})
		}
		return nil
	}
}

func GetInfoChannelHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		channelIDStr := strings.TrimSpace(update.Message.Text[len("/info"):])
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return nil
		}

		channel, err := app.ChannelService.GetChannelByID(context.Background(), channelID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return nil
		}

		owner, err := app.UserService.GetUserByID(context.Background(), channel.OwnerID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Dono não encontrado!",
			})
			return nil
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
		})
		return nil
	}
}

func RegisterTransferHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		input := strings.TrimSpace(update.Message.Text[len("/transfer"):])
		parts := strings.Fields(input)
		if len(parts) < 2 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Uso: /transfer <channelId> <newOwnerId>",
			})
			return nil
		}

		channelID, _ := strconv.ParseInt(parts[0], 10, 64)
		newOwnerID, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := app.ChannelService.GetChannelByID(context.Background(), channelID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return nil
		}

		tgUser, err := bot.GetChat(context.Background(), &telego.GetChatParams{ChatID: telego.ChatID{ID: newOwnerID}})
		if err != nil || tgUser.FirstName == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ ID de usuário inválido: %d", newOwnerID),
			})
			return nil
		}

		err = app.ChannelService.UpdateOwnerChannel(context.Background(), channelID, channel.OwnerID, newOwnerID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Erro ao transferir canal",
			})
			return nil
		}

		msg := fmt.Sprintf(
			"✅ <b>Transferência realizada com sucesso!</b>\n<b>Canal:</b> %s\n<b>ID:</b> %d\n<b>Novo Dono:</b> %s (%d)\n\n🔗 <a href=\"%s\">Abrir painel do canal</a>",
			html.EscapeString(channel.Title),
			channelID,
			html.EscapeString(tgUser.FirstName),
			newOwnerID,
			auth.GenerateMiniAppUrl(parts[1], parts[0]),
		)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
		})
		return nil
	}
}

func GetInfoUserHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		userIDStr := strings.TrimSpace(update.Message.Text[len("/user"):])
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return nil
		}

		user, err := app.UserService.GetUserByID(context.Background(), userID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Usuário não encontrado!: %v", err),
			})
			return nil
		}

		channels, _ := app.ChannelService.GetUserChannels(context.Background(), user.UserId)
		header := fmt.Sprintf("👤 <b><a href='tg://user?id=%d'>%s</a></b> (<code>%d</code>)\n📦 Canais: <b>%d</b>\n\n",
			user.UserId,
			html.EscapeString(user.FirstName),
			user.UserId,
			len(channels),
		)

		if len(channels) == 0 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      header + "Usuário ainda não possui canais.",
				ParseMode: telego.ModeHTML,
			})
			return nil
		}

		const chunkSize = 20
		for i := 0; i < len(channels); i += chunkSize {
			chunk := channels[i:minInt(i+chunkSize, len(channels))]
			var lines []string
			for _, c := range chunk {
				lines = append(lines, fmt.Sprintf("<a href='%s'>%s</a> - <code>%d</code>",
					c.InviteURL,
					html.EscapeString(c.Title),
					c.ID,
				))
			}

			msg := header + strings.Join(lines, "\n")
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      msg,
				ParseMode: telego.ModeHTML,
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
		}
		return nil
	}
}

func RemoveChannelHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		channelIDStr := strings.TrimSpace(update.Message.Text[len("/remove"):])
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return nil
		}

		channel, err := app.ChannelService.GetChannelByID(context.Background(), channelID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Canal não encontrado!: %v", err),
			})
			return nil
		}

		if err = app.ChannelService.DisconnectChannel(context.Background(), channel.OwnerID, channelID); err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ Não foi possivel deletar o canal: %v", err),
			})
			return nil
		}

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      "✅ Canal excluído com sucesso!",
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func NoticeCommandHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		if update.Message == nil {
			return nil
		}

		text := update.Message.Text
		if text == "" {
			return nil
		}

		lines := strings.Split(text, "\n")
		if len(lines) <= 1 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ A mensagem de aviso está vazia.",
			})
			return nil
		}

		noticeText := strings.TrimSpace(strings.Join(lines[1:], "\n"))

		users, err := app.UserService.GetAllUsers(context.Background())
		if err != nil || len(users) == 0 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Nenhum usuário encontrado.",
			})
			return nil
		}

		var failedUsers []string
		sentCount := 0

		for _, user := range users {
			_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: user.UserId},
				Text:      noticeText,
				ParseMode: telego.ModeHTML,
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      finalMsg.String(),
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func NoticeChannelsHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		user, _ := bot.GetMe(context.Background())

		data := map[string]string{
			"botUsername": user.Username,
		}

		text, kb := parser.GetMessageTelego("publi", data)

		channels, err := app.ChannelService.GetAllChannels(context.Background())
		if err != nil || len(channels) == 0 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Nenhum canal encontrado.",
			})
			return nil
		}

		var failedChannels []string
		sentCount := 0

		for _, ch := range channels {
			_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:      telego.ChatID{ID: ch.ID},
				Text:        text,
				ParseMode:   telego.ModeHTML,
				ReplyMarkup: kb,
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      resultMsg.String(),
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func SendMessageToIdHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		if update.Message == nil || update.Message.Text == "" {
			return nil
		}

		lines := strings.Split(update.Message.Text, "\n")
		if len(lines) < 2 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Uso inválido. Envie no formato:\n/send <id>\n<mensagem>",
			})
			return nil
		}

		idStr := strings.TrimSpace(strings.TrimPrefix(lines[0], "/send"))
		targetID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   fmt.Sprintf("❌ ID inválido: %v", err),
			})
			return nil
		}

		message := strings.Join(lines[1:], "\n")
		message = strings.TrimSpace(message)

		if message == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Mensagem vazia.",
			})
			return nil
		}

		_, err = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: targetID},
			Text:      message,
			ParseMode: telego.ModeHTML,
		})
		if err != nil {
			logger.Error("ADMIN", "Erro ao enviar mensagem para %d: %v", targetID, err)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      fmt.Sprintf("❌ Falha ao enviar para <code>%d</code>: %v", targetID, err),
				ParseMode: telego.ModeHTML,
			})
			return nil
		}

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      fmt.Sprintf("✅ Mensagem enviada para <code>%d</code> com sucesso.", targetID),
			ParseMode: telego.ModeHTML,
		})
		return nil
	}
}

func AddChannelCommandHandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		botInfo, _ := bot.GetMe(context.Background())

		msgText := strings.TrimSpace(update.Message.Text)
		args := strings.SplitN(msgText, " ", 3)
		if len(args) < 3 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Uso correto: /add <channel_id> <owner_id>",
			})
			return nil
		}

		channelIDStr := args[1]
		ownerIDStr := args[2]
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		ownerID, err2 := strconv.ParseInt(ownerIDStr, 10, 64)
		if err != nil || err2 != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ IDs inválidos. Certifique-se de que ambos são numéricos.",
			})
			return nil
		}

		existingChannel, _ := c.ChannelService.GetChannelByID(context.Background(), channelID)
		if existingChannel != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Canal já existe no banco de dados.",
			})
			return nil
		}

		channelInfo, err := bot.GetChat(context.Background(), &telego.GetChatParams{ChatID: telego.ChatID{ID: channelID}})
		if err != nil {
			logger.Error("ADMIN", "Erro ao buscar canal: %v", err)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{ChatID: update.Message.Chat.ChatID(), Text: "❌ Erro ao buscar informações do canal."})
			return nil
		}

		ownerInfo, err := bot.GetChat(context.Background(), &telego.GetChatParams{ChatID: telego.ChatID{ID: ownerID}})
		if err != nil {
			logger.Error("ADMIN", "Erro ao buscar usuário: %v", err)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{ChatID: update.Message.Chat.ChatID(), Text: "❌ Erro ao buscar informações do usuário."})
			return nil
		}

		_ = c.UserService.UpsertUser(context.Background(), &userModes.User{
			UserId:    ownerID,
			FirstName: utils.RemoveHTMLTags(ownerInfo.FirstName),
		})

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

		channel, err := c.ChannelService.CreateChannelWithDefaults(context.Background(), channelID, channelInfo.Title, inviteURL, newPackCaption, defaultCaption, ownerID)
		if err != nil {
			logger.Error("ADMIN", "Erro ao criar canal: %v", err)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{ChatID: update.Message.Chat.ChatID(), Text: "❌ Erro ao salvar canal."})
			return nil
		}

		miniApp := auth.GenerateMiniAppUrl(fmt.Sprintf("%d", ownerID), fmt.Sprintf("%d", channelID))
		msg := fmt.Sprintf("✅ Canal salvo com sucesso - (%s - %d)\n\n%s", channel.Title, channel.ID, miniApp)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
		})
		return nil
	}
}

func NoticeUsersReplyHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		if update.Message == nil {
			return nil
		}

		if update.Message.ReplyToMessage == nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Responda a uma message para enviar o aviso aos usuários.",
			})
			return nil
		}

		noticeText := update.Message.ReplyToMessage.Text
		if noticeText == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ A mensagem respondida está vazia.",
			})
			return nil
		}

		users, err := app.UserService.GetAllUsers(context.Background())
		if err != nil || len(users) == 0 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Nenhum usuário encontrado.",
			})
			return nil
		}

		var failedUsers []string
		sentCount := 0

		for _, user := range users {
			_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: user.UserId},
				Text:      noticeText,
				ParseMode: telego.ModeHTML,
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      finalMsg.String(),
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func NoticeChannelsReplyHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		if update.Message == nil {
			return nil
		}

		if update.Message.ReplyToMessage == nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Responda a uma mensagem para enviar o aviso aos canais.",
			})
			return nil
		}

		noticeText := update.Message.ReplyToMessage.Text
		if noticeText == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ A mensagem respondida está vazia.",
			})
			return nil
		}

		channels, err := app.ChannelService.GetAllChannels(context.Background())
		if err != nil || len(channels) == 0 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Nenhum canal encontrado.",
			})
			return nil
		}

		var failedChannels []string
		sentCount := 0

		for _, ch := range channels {
			_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: ch.ID},
				Text:      noticeText,
				ParseMode: telego.ModeHTML,
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      resultMsg.String(),
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func ToggleMaintenceHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
		if update.Message == nil {
			return nil
		}

		maintenance, err := app.ServerService.ToggleMaintenance(context.Background())
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Erro ao alterar o modo de manutenção.",
			})
			return nil
		}

		var msg string
		if maintenance {
			msg = "⚠️ <b>Modo de manutenção ativado</b>\n\nO bot pode ficar temporariamente indisponível."
		} else {
			msg = "✅ <b>Modo de manutenção desativado</b>\n\nO bot voltou a funcionar normalmente."
		}

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return nil
	}
}

func SetAdminHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, upt telego.Update) error {
		bot := ctx.Bot()
		if upt.Message == nil {
			return nil
		}

		args := strings.Fields(upt.Message.Text)
		if len(args) < 2 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    upt.Message.Chat.ChatID(),
				Text:      "❌ Uso correto:\n<code>/setadmin [userID]</code>",
				ParseMode: telego.ModeHTML,
			})
			return nil
		}

		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: upt.Message.Chat.ChatID(),
				Text:   "❌ userID inválido.",
			})
			return nil
		}

		isAdmin, err := app.UserService.UpdateUserAdmin(context.Background(), userID)
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: upt.Message.Chat.ChatID(),
				Text:   "❌ Erro ao alterar status de admin.",
			})
			return nil
		}

		var msg string
		if isAdmin {
			msg = fmt.Sprintf("✅ Usuário <code>%d</code> agora é administrador.", userID)
		} else {
			msg = fmt.Sprintf("⚠️ Usuário <code>%d</code> não é mais administrador.", userID)
		}

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    upt.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
		})
		return nil
	}
}

func AdminHelpHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		bot := ctx.Bot()
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

		kb := &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{
						Text: "📊 Abrir Dashboard Admin",
						WebApp: &telego.WebAppInfo{
							URL: fmt.Sprintf("%s/admin/dash", config.WebAppURL),
						},
					},
				},
			},
		}

		params := &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msg,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}

		_, _ = bot.SendMessage(context.Background(), params)
		return nil
	}
}

func CheckBotAdminHandlerTelego(app *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		const targetBotID = 5986082367
		const targetBotUser = "@XavolaBot"
		bot := ctx.Bot()

		channels, err := app.ChannelService.GetAllChannels(context.Background())
		if err != nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ Erro ao buscar canais do banco.",
			})
			return nil
		}

		statusMsg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      fmt.Sprintf("🔍 Verificando <b>%d</b> canais... Isso pode levar um tempo.", len(channels)),
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})

		var foundChannels []string
		var mu sync.Mutex
		count := 0

		chQueue := make(chan *userModes.Channel, len(channels))
		for i := range channels {
			chQueue <- &channels[i]
		}
		close(chQueue)

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
					member, err := bot.GetChatMember(context.Background(), &telego.GetChatMemberParams{
						ChatID: telego.ChatID{ID: ch.ID},
						UserID: targetBotID,
					})

					if err == nil {
						status := member.MemberStatus()
						if status == telego.MemberStatusAdministrator || status == telego.MemberStatusCreator {
							owner, _ := app.UserService.GetUserByID(context.Background(), ch.OwnerID)
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
				}
			}()
		}
		wg.Wait()

		if count == 0 {
			_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
				ChatID:    update.Message.Chat.ChatID(),
				MessageID: statusMsg.MessageID,
				Text:      fmt.Sprintf("✅ O bot %s não foi encontrado como admin em nenhum canal.", targetBotUser),
			})
			return nil
		}

		header := fmt.Sprintf("🤖 <b>Bot %s encontrado em %d canais:</b>\n\n", targetBotUser, count)

		var sb strings.Builder
		sb.WriteString(header)
		first := true
		for i, info := range foundChannels {
			if sb.Len()+len(info) > 3800 {
				if first {
					_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
						ChatID:             update.Message.Chat.ChatID(),
						MessageID:          statusMsg.MessageID,
						Text:               sb.String(),
						ParseMode:          telego.ModeHTML,
						LinkPreviewOptions: &telego.LinkPreviewOptions{IsDisabled: true},
					})
					first = false
				} else {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID:             update.Message.Chat.ChatID(),
						Text:               sb.String(),
						ParseMode:          telego.ModeHTML,
						LinkPreviewOptions: &telego.LinkPreviewOptions{IsDisabled: true},
					})
				}
				sb.Reset()
			}
			sb.WriteString(info + "────────────────────\n")
			if i == len(foundChannels)-1 {
				if first {
					_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
						ChatID:             update.Message.Chat.ChatID(),
						MessageID:          statusMsg.MessageID,
						Text:               sb.String(),
						ParseMode:          telego.ModeHTML,
						LinkPreviewOptions: &telego.LinkPreviewOptions{IsDisabled: true},
					})
				} else {
					_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
						ChatID:             update.Message.Chat.ChatID(),
						Text:               sb.String(),
						ParseMode:          telego.ModeHTML,
						LinkPreviewOptions: &telego.LinkPreviewOptions{IsDisabled: true},
					})
				}
			}
		}
		return nil
	}
}
