package commands

import (
	"github.com/go-telegram/bot"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands/help"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands/start"
)

func LoadCommandHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, start.Handler())
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, help.Handler())
}
