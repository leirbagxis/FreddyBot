package coupon

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/telegram/modules"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskCouponHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		sessionKey := parts[1]
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

		c.CacheService.SetAwaitingCoupon(ctx, userID, sessionKey)
		text, _ := parser.GetMessage("coupon-required-message", map[string]string{
			"firstName": fmt.Sprintf("%s", update.CallbackQuery.From.FirstName),
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      text,
			ParseMode: "HTML",
		})
	}
}

func CheckCouponHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		from := update.Message.From
		couponText := update.Message.Text

		sessionKey, err := c.CacheService.GetAwaitingCoupon(ctx, from.ID)
		if err != nil {
			return
		}

		coupon, err := c.CouponService.ValidateCoupon(ctx, couponText, from.ID)
		if coupon == nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    from.ID,
				Text:      fmt.Sprintf("⚠️ Não foi possível aplicar o cupom.\n<i>%v</i>", err),
				ParseMode: models.ParseModeHTML,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		payment, err := c.PaymentService.GetPaymentWithPayload(ctx, sessionKey)
		if err != nil {
			fmt.Errorf("Erro ao pegar pagamento: %s || %w", sessionKey, err)
		}

		err = c.CouponService.ApplyCouponToPayment(ctx, coupon, payment, from.ID)
		fmt.Println(err)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: from.ID,
				Text:   fmt.Sprintf("❌ Não foi possível aplicar o cupom: %v", err),
			})
			return
		}

		modules.SendChannelActivationPayment(ctx, b, update, c, sessionKey)

	}
}
