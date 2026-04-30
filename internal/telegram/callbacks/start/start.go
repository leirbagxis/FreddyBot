package start

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		text, button := parser.GetMessage("start", map[string]string{
			"firstName": utils.RemoveHTMLTags(update.CallbackQuery.From.FirstName),
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func CheckSubscriptionHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const channelID = -1003767126116

		member, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
			ChatID: channelID,
			UserID: update.CallbackQuery.From.ID,
		})

		if err != nil {
			// Se der erro, pode ser que o bot não esteja no canal como admin
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao verificar sua inscrição. O bot precisa ser administrador no canal oficial!",
				ShowAlert:       true,
			})
			return
		}

		isMember := false
		if member != nil {
			switch member.Type {
			case "creator", "administrator", "member", "restricted":
				isMember = true
			}
		}

		if !isMember {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você ainda não entrou no canal!",
				ShowAlert:       true,
			})
			return
		}

		// Se entrou, mostra o start normal
		text, button := parser.GetMessage("start", map[string]string{
			"firstName": utils.RemoveHTMLTags(update.CallbackQuery.From.FirstName),
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Obrigado por entrar!",
		})
	}
}
