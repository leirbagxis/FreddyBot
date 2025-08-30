package mychannel

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskDeleteChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		cbks := update.CallbackQuery

		userId := cbks.From.ID
		session, err := c.CacheService.GetSelectedChannel(ctx, userId)
		if err != nil {
			log.Printf("Erro ao pegar sessão: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert:       true,
			})
			return
		}

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, session)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		data := map[string]string{
			"title":     channel.Title,
			"channelId": fmt.Sprintf("%d", session),
		}
		c.CacheService.SetDeleteChannel(ctx, userId, session)

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
		cbks := update.CallbackQuery

		userId := cbks.From.ID
		session, err := c.CacheService.GetDeleteChannel(ctx, userId)
		if err != nil {
			log.Printf("Erro ao criar sessão: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert:       true,
			})
			return
		}

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, session)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		err = c.ChannelRepo.DeleteChannelWithRelations(ctx, userId, session)
		if err != nil {
			log.Printf("Erro ao excluir canal: %v", err)
			return
		}

		data := map[string]string{
			"title":     channel.Title,
			"channelId": fmt.Sprintf("%d", session),
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

		_, err = c.CacheService.DeleteAllUserSessionsBySuffix(ctx, userId)
		if err != nil {
			log.Printf("Erro ao excluir all sessions: %v", err)
			return
		}
	}
}
