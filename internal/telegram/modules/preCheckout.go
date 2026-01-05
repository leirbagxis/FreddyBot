package modules

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func PreCheckoutHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		query := update.PreCheckoutQuery
		if query == nil {
			return
		}

		payload := query.InvoicePayload
		parts := strings.Split(payload, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", payload)
			return
		}

		session, err := c.SessionManager.GetChannelSession(ctx, parts[1])
		if err != nil || session == nil {
			b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
				PreCheckoutQueryID: query.ID,
				ErrorMessage:       "⌛ Tempo Esgotado. Faça o processo de pagamento novamente!",
				OK:                 false,
			})
			return
		}

		b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
		})
	}
}
