package mychannel

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	separatorModels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskStickerSeparatorHandler(c *container.AppContainer) bot.HandlerFunc {
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
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", session),
		}

		text, button := parser.GetMessage("ask-separator-message", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func RequireStickerSeparatorHandler(c *container.AppContainer) bot.HandlerFunc {
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
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", session),
		}
		c.CacheService.SetAwaitingStickerSeparator(ctx, userId, session)

		text, button := parser.GetMessage("require-separator-message", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func SetStickerSeparatorHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		userId := update.Message.From.ID
		channelId, err := c.CacheService.GetAwaitingStickerSeparator(ctx, userId)
		if err != nil {
			log.Printf("Erro ao buscar cache sticker: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				//CallbackQueryID: update.CallbackQuery.ID,
				Text:      "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert: true,
			})
			return
		}

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		stickerId := update.Message.Sticker.FileID

		file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: stickerId})
		if err != nil {
			fmt.Println("erro ao obter sticker")
		}
		stickerLink := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", config.TelegramBotToken, file.FilePath)
		fmt.Println(stickerLink)

		separator := &separatorModels.Separator{
			OwnerChannelID: channelId,
			SeparatorID:    stickerId,
			SeparatorURL:   stickerLink,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err = c.SeparatorRepo.SaveSeparator(ctx, separator); err != nil {
			text, button := parser.GetMessage("failed-save-separator", map[string]string{
				"channelId": fmt.Sprintf("%d", channelId),
			})

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        text,
				ReplyMarkup: button,
				ParseMode:   "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
		}

		c.CacheService.DeleteAwaitingStickerSeparator(ctx, userId)
		text, button := parser.GetMessage("success-save-separator", map[string]string{
			"channelId":   fmt.Sprintf("%d", channelId),
			"channelName": channel.Title,
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func DeleteSeparatorHandler(c *container.AppContainer) bot.HandlerFunc {
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

		separator, err := c.SeparatorRepo.GetSeparatorByOwnerChannelID(ctx, channel.ID)

		if separator == nil || err != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você ainda não possui nenhum separador vinculado.",
				ShowAlert:       true,
			})
			return
		}
		fmt.Println(separator, err)

		err = c.SeparatorRepo.DeleteSeparatorByOwnerChannelId(ctx, session)
		if err != nil {
			log.Printf("Erro ao excluir separator: %v", err)
			return
		}

		data := map[string]string{
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", session),
		}

		text, button := parser.GetMessage("success-delete-separator", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Separador excluido com sucesso!",
		})
	}
}
