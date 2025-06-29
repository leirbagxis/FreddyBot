package mychannel

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID

		channelsButtons, err := c.ButtonRepo.GetUserChannelsAsButtons(ctx, userID)
		if err != nil || len(channelsButtons) == 0 {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você ainda não possui nenhum canal vinculado.",
				ShowAlert:       true,
			})
			return
		}

		text, homeButtons := parser.GetMessage("profile-user-channels", map[string]string{})

		var finalButtons [][]parser.Button

		// Adiciona os botões dinâmicos
		for _, row := range channelsButtons {
			finalButtons = append(finalButtons, row)
		}

		if homeButtons != nil {
			for _, row := range homeButtons.InlineKeyboard {
				var newRow []parser.Button
				for _, btn := range row {
					newRow = append(newRow, parser.Button{
						Text:              btn.Text,
						CallbackData:      btn.CallbackData,
						SwitchInlineQuery: btn.SwitchInlineQuery,
						URL:               btn.URL,
					})
				}
				finalButtons = append(finalButtons, newRow)
			}
		}

		replyMarkup := parser.BuildInlineKeyboard(finalButtons)

		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        text,
			ReplyMarkup: replyMarkup,
			ParseMode:   "HTML",
		})

		if err != nil {
			log.Printf("Erro ao editar mensagem: %v", err)
		}
	}
}
