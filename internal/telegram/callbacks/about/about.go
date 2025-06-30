package about

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		user, _ := b.GetMe(ctx)
		text, button := parser.GetMessage("about", map[string]string{
			"ownerUser":  "@SuporteLegendas",
			"botVersion": "beta 1.0.0",
			"botId":      fmt.Sprintf("%d", user.ID),
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}
