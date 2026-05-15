package mychannel

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskDeleteChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		cbks := update.CallbackQuery

		userId := cbks.From.ID
		session, err := c.CacheService.GetSelectedChannel(ctx, userId)
		if err != nil {
			logger.Error("BOT", "Erro ao pegar sessão: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert:       true,
			})
			return
		}

		channel, err := c.ChannelService.GetChannelByTwoID(ctx, userId, session)
		if err != nil {
			logger.Error("BOT", "Erro ao buscar canal: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Canal não encontrado ou não pertence a você!",
				ShowAlert:       true,
			})
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
			logger.Error("BOT", "Erro ao criar sessão: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert:       true,
			})
			return
		}

		channel, err := c.ChannelService.GetChannelByTwoID(ctx, userId, session)
		if err != nil {
			logger.Error("BOT", "Erro ao buscar canal: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Canal não encontrado ou não pertence a você!",
				ShowAlert:       true,
			})
			return
		}

		err = c.ChannelService.DisconnectChannel(ctx, userId, session)
		if err != nil {
			logger.Error("BOT", "Erro ao excluir canal: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao excluir canal. Tente novamente!",
				ShowAlert:       true,
			})
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
			logger.Error("BOT", "Erro ao excluir all sessions: %v", err)
			return
		}
	}
}
