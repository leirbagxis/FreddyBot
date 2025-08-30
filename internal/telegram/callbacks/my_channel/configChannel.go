package mychannel

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
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

		userIDStr := fmt.Sprintf("%d", userID)

		data := map[string]string{
			"title":     channel.Title,
			"channelId": channelIdString,
			"webAppUrl": auth.GenerateMiniAppUrl(userIDStr, channelIdString),
		}
		text, button := parser.GetMessage("config-channel", data)

		err = c.CacheService.SetSelectedChannel(ctx, userID, channelId)
		if err != nil {
			log.Printf("Erro ao criar sessão: %v", err)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

	}
}
