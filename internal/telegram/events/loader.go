package events

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	addchannel "github.com/leirbagxis/FreddyBot/internal/telegram/events/addChannel"
)

func LoadEvents(b *bot.Bot) {
	b.RegisterHandlerMatchFunc(matchMyChatMember, addchannel.Handler())
}

func matchMyChatMember(update *models.Update) bool {
	fmt.Println("Checking: ", update.MyChatMember != nil)
	return update.MyChatMember != nil
}
