package mychannel

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func ConfigHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		channelIdString := parts[1]
		channelId, err := strconv.ParseInt(channelIdString, 10, 64)
		if err != nil {
			log.Println("Error parsing channelId:", err)
			return
		}
		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userID, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		data := map[string]string{
			"title":     channel.Title,
			"channelId": channelIdString,
			"webAppUrl": "https://caption.chelodev.shop/703450014/-1002824722434?signature=f830adccc8cadcf0b84ad8f6236dc4c22f39f67ef7bc1e0c223fa498a9a8cf89",
		}
		text, button := parser.GetMessage("config-channel", data)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

	}
}
