package suporte

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		textInput := update.Message.Text
		command := "/ouvidoria"
		
		message := strings.TrimSpace(textInput[len(command):])

		if len(message) <= 0 {
			text, button := parser.GetMessage("support-usage", map[string]string{})
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        text,
				ReplyMarkup: button,
				ParseMode:   models.ParseModeHTML,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		firstName := html.EscapeString(utils.RemoveHTMLTags(update.Message.From.FirstName))
		safeMessage := html.EscapeString(message)

		user, _ := b.GetMe(ctx)

		data := map[string]string{
			"firstName":   firstName,
			"userId":      fmt.Sprintf("%d", update.Message.From.ID),
			"botId":       fmt.Sprintf("%d", user.ID),
			"userMessage": safeMessage,
		}

		adminText, _ := parser.GetMessage("support-msg-admin", data)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    config.OwnerID,
			Text:      adminText,
			ParseMode: models.ParseModeHTML,
		})
		if err != nil {
			logger.Error("BOT", "Erro ao enviar mensagem de ouvidoria pro admin: %v", err)
		}

		text, button := parser.GetMessage("support-sent", data)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
		if err != nil {
			logger.Error("BOT", "Erro ao enviar confirmação de ouvidoria pro usuário %d: %v", update.Message.From.ID, err)
		}
	}
}
