package events

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	addchannel "github.com/leirbagxis/FreddyBot/internal/telegram/events/addChannel"
)

func LoadEvents(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandlerMatchFunc(matchMyChatMember, addchannel.Handler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add-yes:", bot.MatchTypePrefix, addchannel.AddYesHandler(c))
}

func matchMyChatMember(update *models.Update) bool {
	fmt.Println("Checking: ", update.MyChatMember != nil)
	return update.MyChatMember != nil
}
