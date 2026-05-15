package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func CheckMaintenceMiddleware(c *container.AppContainer) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, upt *models.Update) {
			if upt.ChannelPost != nil {
				next(ctx, b, upt)
				return
			}

			// 1. Get maintenance status first
			maintenance, err := c.ServerService.GetMaintenance(ctx)
			if err != nil {
				logger.Error("MID", "Erro ao buscar status de manutenção: %v", err)
				next(ctx, b, upt)
				return
			}

			// 2. If not in maintenance, allow everything
			if !maintenance {
				next(ctx, b, upt)
				return
			}

			// 3. If in maintenance, check if user is admin or owner
			userID := getUpdateUserID(upt)
			if userID != 0 {
				// Owner sempre passa
				if userID == config.OwnerID {
					next(ctx, b, upt)
					return
				}

				// Admins também passam
				user, err := c.UserService.GetUserByID(ctx, userID)
				if err == nil && user != nil && user.IsAdmin {
					next(ctx, b, upt)
					return
				}
			}

			// 4. Send maintenance response
			sendMaintenceResponse(ctx, b, upt, userID)
		}
	}
}


func sendMaintenceResponse(ctx context.Context, b *bot.Bot, upt *models.Update, userID int64) {
	switch {
	case upt.Message != nil:
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      "⚙️ <b>Manutenção em andamento</b>\n\nO bot está temporariamente indisponível para melhorias no sistema.\n\n⏳ <i>Voltaremos em breve.</i>\n\nObrigado pela sua paciência 💙",
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: upt.Message.ID,
			},
		})

	case upt.CallbackQuery != nil:
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: upt.CallbackQuery.ID,
			Text:            "🚧 O bot está em manutenção no momento.",
			ShowAlert:       true,
			CacheTime:       0,
		})

	case upt.InlineQuery != nil:
		_, _ = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: upt.InlineQuery.ID,
			Results: []models.InlineQueryResult{
				&models.InlineQueryResultArticle{
					ID:    "maintenance",
					Title: "⚙️ Manutenção em andamento",
					InputMessageContent: &models.InputTextMessageContent{
						MessageText: "O bot está em manutenção no momento. Tente novamente mais tarde.",
					},
				},
			},
			CacheTime: 0,
		})

	case upt.ChannelPost != nil:
		// Se quiser bloquear post de canal silenciosamente, não faz nada.
		return
	}
}
