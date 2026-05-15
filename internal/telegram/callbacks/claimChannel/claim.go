package claimchannel

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		query := update.InlineQuery.Query
		from := update.InlineQuery.From

		cmd := strings.Split(query, " ")
		if len(cmd) == 0 || strings.ToLower(cmd[0]) != "claim" {
			return
		}

		if len(cmd) < 2 {
			article := buildErrorArticle(from.ID, "invalid", "🔎 ID Inválido", "⚠️ O ID deve ser um número válido.")
			answerInline(ctx, b, update, []models.InlineQueryResult{article})
			return
		}

		channelIdRaw := strings.TrimSpace(cmd[1])
		channelIdRaw = strings.Map(func(r rune) rune {
			// Remove caracteres Unicode LRI, RLI, FSI, PDI
			if r >= 0x2066 && r <= 0x2069 {
				return -1
			}
			return r
		}, channelIdRaw)

		channelID, err := strconv.ParseInt(channelIdRaw, 10, 64)
		if err != nil {
			article := buildErrorArticle(from.ID, "invalid_format", "🔎 ID Inválido", "⚠️ O ID deve ser um número válido.")
			answerInline(ctx, b, update, []models.InlineQueryResult{article})
			return
		}

		// Consulta no banco — pode já existir ou não
		channel, _ := c.ChannelService.GetChannelByID(ctx, channelID)

		// Verifica se o usuário é o criador do canal no Telegram
		admins, err := b.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
			ChatID: channelID,
		})
		if err != nil {
			article := buildErrorArticle(from.ID, "admin_error", "⚠️ Erro ao acessar canal", "🚫 O bot pode não estar no canal ou o ID está errado.")
			answerInline(ctx, b, update, []models.InlineQueryResult{article})
			return
		}

		isOwner := false
		for _, admin := range admins {
			if admin.Type == "creator" && admin.Owner.User.ID == from.ID {
				isOwner = true
				break
			}
		}
		if !isOwner {
			article := buildErrorArticle(from.ID, "not_owner", "🚫 Sem Permissão", fmt.Sprintf("⚠️ Você não é o criador/dono deste canal (%d).", channelID))
			answerInline(ctx, b, update, []models.InlineQueryResult{article})
			return
		}

		// Criador confirmado. Agora preparamos a mensagem de controle.

		ownerUser, err := c.UserService.GetUserByID(ctx, channel.OwnerID)
		if err != nil {
			article := buildErrorArticle(from.ID, "user_error", "❌ Erro ao buscar proprietário", "Erro ao buscar usuário proprietário.")
			answerInline(ctx, b, update, []models.InlineQueryResult{article})
			return
		}

		// Criar sessão temporária no Redis
		session, err := c.SessionManager.CreateClaimSession(ctx, channelID, channel.OwnerID, from.ID)
		if err != nil {
			logger.Error("BOT", "Erro ao criar cache: %v", err)
			return
		}

		vars := map[string]string{
			"channelId":   fmt.Sprint(channelID),
			"channelName": channel.Title,
			"ownerId":     fmt.Sprint(ownerUser.UserId),
			"ownerName":   ownerUser.FirstName,
			"sessionKey":  session.Key,
		}

		text, buttons := parser.GetMessage("claim-ownership-require-message", vars)

		article := &models.InlineQueryResultArticle{
			ID:          fmt.Sprintf("claim_success_%d", from.ID),
			Title:       "✅ Canal Encontrado",
			Description: fmt.Sprintf("Canal %d - Confirme a propriedade", channelID),
			InputMessageContent: &models.InputTextMessageContent{
				MessageText: text,
				ParseMode:   "HTML",
			},
			ReplyMarkup: buttons,
		}

		answerInline(ctx, b, update, []models.InlineQueryResult{article})
	}
}

func AcceptClaimHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callback := update.CallbackQuery
		from := update.CallbackQuery.From

		callbackData := callback.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			logger.Warn("BOT", "Callback invalido: %s", callbackData)
			return
		}

		getSession, err := c.SessionManager.GetChannelSession(ctx, parts[1])
		if err != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Tempo Esgotado. Faça o processo novamente!",
				ShowAlert:       true,
			})
			return
		}

		channel, err := c.ChannelService.GetChannelByID(ctx, getSession.ChannelID)
		if err != nil {
			logger.Error("BOT", "Erro ao obter info do canal: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao obter informações do canal!",
				ShowAlert:       true,
			})
			return
		}

		c.SessionManager.DeleteChannelSession(ctx, parts[1])

		err = c.ChannelService.UpdateOwnerChannel(ctx, getSession.ChannelID, getSession.OwnerID, getSession.NewOwnerID)
		if err != nil {
			logger.Error("BOT", "Erro ao transferir posse do canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Erro ao passar a posse para o novo usuário.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}
		userId := fmt.Sprintf("%d", getSession.NewOwnerID)
		channelId := fmt.Sprintf("%d", channel.ID)

		data := map[string]string{
			"channelId":   fmt.Sprintf("%d", getSession.ChannelID),
			"channelName": channel.Title,
			"miniAppUrl":  auth.GenerateMiniAppUrl(userId, channelId),
		}

		textNew, buttonNew := parser.GetMessage("success-new-paccess-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      from.ID,
			Text:        textNew,
			ReplyMarkup: buttonNew,
			ParseMode:   "HTML",
		})
	}
}

func buildErrorArticle(userID int64, idSuffix, title, message string) *models.InlineQueryResultArticle {
	return &models.InlineQueryResultArticle{
		ID:    fmt.Sprintf("error_%s_%d", idSuffix, userID),
		Title: title,
		InputMessageContent: &models.InputTextMessageContent{
			MessageText: message,
			ParseMode:   "HTML",
		},
	}
}

func answerInline(ctx context.Context, b *bot.Bot, update *models.Update, results []models.InlineQueryResult) {
	_, err := b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
		CacheTime:     0,
	})
	if err != nil {
		logger.Error("BOT", "Erro ao responder inline_query: %v", err)
	}
}
