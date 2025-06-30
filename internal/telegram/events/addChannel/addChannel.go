package addchannel

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chat := update.MyChatMember.Chat
		from := update.MyChatMember.From

		if update.MyChatMember.OldChatMember.Type != "left" {
			return
		}

		getChannel, _ := c.ChannelRepo.GetChannelByID(ctx, chat.ID)

		data := map[string]string{
			"channelName": chat.Title,
		}

		if getChannel != nil {
			text, button := parser.GetMessage("toadd-exist-channel", data)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      from.ID,
				Text:        text,
				ReplyMarkup: button,
				ParseMode:   "HTML",
			})
			return
		}

		session, err := c.SessionManager.CreateChannelSession(ctx, chat.ID, from.ID, chat.Title)
		if err != nil {
			log.Printf("Erro ao criar sessão: %v", err)
			return
		}
		data["sessionKey"] = session.Key

		text, button := parser.GetMessage("toadd-require-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      from.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
		})
		fmt.Println(from)
	}
}

func AddYesHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		getSession, err := c.SessionManager.GetChannelSession(ctx, parts[1])
		if err != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Tempo Esgotado. Faça o processo novamente!",
				ShowAlert:       true,
			})
			return
		}

	}
}
