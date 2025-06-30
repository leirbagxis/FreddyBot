package mychannel

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskDeleteChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callbackData := update.CallbackQuery.Data

		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		userId := update.CallbackQuery.From.ID
		channelId, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		data := map[string]string{
			"title":     channel.Title,
			"channelId": fmt.Sprintf("%d", channelId),
		}

		text, button := parser.GetMessage("del", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func ConfirmDeleteChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callbackData := update.CallbackQuery.Data

		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		userId := update.CallbackQuery.From.ID
		channelId, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		err = c.ChannelRepo.DeleteChannelWithRelations(ctx, userId, channelId)
		if err != nil {
			if err != nil {
				log.Printf("Erro ao excluir canal: %v", err)
				return
			}
		}

		data := map[string]string{
			"title":     channel.Title,
			"channelId": parts[1],
		}

		text, button := parser.GetMessage("confirm-del", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Canal excluido com sucesso!",
		})
	}
}
