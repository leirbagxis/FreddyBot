package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	addchannel "github.com/leirbagxis/FreddyBot/internal/telegram/events/addChannel"
	channelpost "github.com/leirbagxis/FreddyBot/internal/telegram/events/channelPost"
	postbuilder "github.com/leirbagxis/FreddyBot/internal/telegram/events/postBuilder"
)

func LoadEvents(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandlerMatchFunc(matchMyChatMember, addchannel.AskAddChannelHandler(c), middleware.CheckAddBotMiddleware(c))
	b.RegisterHandlerMatchFunc(matchForwardedChannel, addchannel.AskForwadedChannelHandler(c), middleware.CheckAddBotMiddleware(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-yes:", bot.MatchTypePrefix, addchannel.AddYesHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-not:", bot.MatchTypePrefix, addchannel.AddNotHandler(c))

	// ## POST BUILDER ## \\
	b.RegisterHandlerMatchFunc(matchPostBuilder(c), postbuilder.Handler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "pb-", bot.MatchTypePrefix, postbuilder.CallbackHandler(c))
	b.RegisterHandlerMatchFunc(matchPostBuilderInline, postbuilder.InlineHandler(c))

	// ## CHANNEL POST ## \\
	b.RegisterHandlerMatchFunc(matchChannelPost, channelpost.Handler(c))

}

func matchPostBuilderInline(update *models.Update) bool {
	return update.InlineQuery != nil && strings.HasPrefix(update.InlineQuery.Query, "pb ")
}

func matchPostBuilder(c *container.AppContainer) bot.MatchFunc {
	return func(update *models.Update) bool {
		if update.Message == nil || update.Message.Chat.Type != models.ChatTypePrivate {
			return false
		}
		// Se for comando, deixa os comandos tratarem
		if strings.HasPrefix(update.Message.Text, "/") {
			return false
		}
		// Match if it has media
		if update.Message.Photo != nil || update.Message.Video != nil || update.Message.Animation != nil || update.Message.Audio != nil || update.Message.Document != nil {
			return true
		}
		// Match if in active session for text input
		state, _ := c.CacheService.GetPostBuilderState(context.Background(), update.Message.From.ID)
		return state != nil && state.Step != ""
	}
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
