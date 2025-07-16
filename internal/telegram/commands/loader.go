package commands

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands/admin"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands/help"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands/start"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func LoadCommandHandlers(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, start.Handler())
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, help.Handler())

	// ## ADMIM COMMANDS ### \\
	b.RegisterHandler(bot.HandlerTypeMessageText, "/users", bot.MatchTypeExact, admin.GetAllUsersHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/channels", bot.MatchTypeExact, admin.GetAllChannelsHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/backup", bot.MatchTypeExact, admin.GetBackUpHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/publi", bot.MatchTypeExact, admin.NoticeChannelsHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/info", bot.MatchTypePrefix, admin.GetInfoChannelHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/transfer", bot.MatchTypePrefix, admin.RegisterTransferHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/user", bot.MatchTypePrefix, admin.GetInfoUserHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remove", bot.MatchTypePrefix, admin.RemoveChannelHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/notice", bot.MatchTypePrefix, admin.NoticeCommandHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/send", bot.MatchTypePrefix, admin.SendMessageToIdHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, admin.AddChannelCommandHandler(c), middleware.CheckAdminMiddleware(config.OwnerID))

}

func matchAdmin(update *models.Update) bool {
	return update.Message != nil && update.Message.From.ID == config.OwnerID
}
