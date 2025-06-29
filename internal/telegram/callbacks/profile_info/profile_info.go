package profileinfo

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
	"gorm.io/gorm"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID

		user, err := c.UserRepo.GetUserById(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "❌ Usuário não encontrado no banco de dados.",
				})
				return
			}
			log.Printf("Erro ao buscar usuário: %v", err)
			return
		}

		countChannel, err := c.ChannelRepo.CountUserChannels(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "❌ Usuário não encontrado no banco de dados.",
				})
				return
			}
			log.Printf("Erro ao buscar countChannel: %v", err)
			return
		}

		data := map[string]string{
			"firstName":    user.FirstName,
			"userId":       fmt.Sprintf("%d", user.UserId),
			"register":     user.CreatedAt.Format("02/01/2006"),
			"countChannel": fmt.Sprintf("%d", countChannel),
		}

		text, button := parser.GetMessage("profile-info", data)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}
