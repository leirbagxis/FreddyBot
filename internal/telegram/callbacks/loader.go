package callbacks

import (
	"github.com/go-telegram/bot"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/about"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/help"
	mychannel "github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/my_channel"
	profileinfo "github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/profile_info"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/start"
)

func LoadCallbacksHandlers(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "help", bot.MatchTypeExact, help.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start", bot.MatchTypeExact, start.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "about", bot.MatchTypeExact, about.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "profile-info", bot.MatchTypeExact, profileinfo.Handler(c))

	// ## MY CHANNEL HANDLERS ## \\
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "profile-user-channels", bot.MatchTypeExact, mychannel.Handler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "config:", bot.MatchTypePrefix, mychannel.ConfigHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "del:", bot.MatchTypePrefix, mychannel.AskDeleteChannelHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm-del:", bot.MatchTypePrefix, mychannel.ConfirmDeleteChannelHandler(c))
}
