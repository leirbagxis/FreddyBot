package modules

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func CancelPayment(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callback := update.CallbackQuery

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		sessionKey := parts[1]

		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		payment, err := c.PaymentService.GetPaymentWithPayload(ctx, sessionKey)
		if err != nil || payment.Status == "canceled" {
			if payment.Status == "canceled" {
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            "⌛ Este pagamento esta cancelado. Faça o processo novamente!",
					ShowAlert:       true,
				})
				return
			}

			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Pagamento Inexistente. Faça o processo novamente!",
				ShowAlert:       true,
			})
			return
		}

		c.PaymentService.CancelPayment(ctx, payment.UserID, sessionKey)
		c.SessionManager.DeleteChannelSession(ctx, sessionKey)
		c.CacheService.DeleteAwaitingCoupon(ctx, payment.UserID)

		// Responder ao callback
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "❌ Este pagamento esta cancelado com sucesso!",
		})

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "❌ Este pagamento esta cancelado com sucesso!",
			ShowAlert:       true,
		})
	}
}
