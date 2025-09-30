package events

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	addchannel "github.com/leirbagxis/FreddyBot/internal/telegram/events/addChannel"
	channelpost "github.com/leirbagxis/FreddyBot/internal/telegram/events/channelPost"
)

func LoadEvents(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandlerMatchFunc(matchMyChatMember, addchannel.AskAddChannelHandler(c), middleware.CheckAddBotMiddleware(c))
	b.RegisterHandlerMatchFunc(matchForwardedChannel, addchannel.AskForwadedChannelHandler(c), middleware.CheckAddBotMiddleware(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-yes:", bot.MatchTypePrefix, addchannel.AddYesHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-not:", bot.MatchTypePrefix, addchannel.AddNotHandler(c))

	// ## CHANNEL POST ## \\
	b.RegisterHandlerMatchFunc(matchChannelPost, channelpost.Handler(c))

}

func matchMyChatMember(update *models.Update) bool {
	fmt.Println("Checking: ", update.MyChatMember != nil)
	return update.MyChatMember != nil
}

func matchForwardedChannel(update *models.Update) bool {
	return update.Message != nil && update.Message.ForwardOrigin != nil && update.Message.ForwardOrigin.MessageOriginChannel != nil
}

func matchChannelPost(update *models.Update) bool {
	return update.ChannelPost != nil
}
