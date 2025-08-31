package help

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

type Deps struct {
	BotUsername string
}

func Handler(deps Deps) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		text, button := parser.GetMessage("help", map[string]string{
			"botUsername": deps.BotUsername,
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
