package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	tgbotModels "github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"gorm.io/gorm"
)

func SaveUserMiddleware(db *gorm.DB) bot.Middleware {
	userRepo := repositories.NewUserRepository(db)

	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *tgbotModels.Update) {
			var userId int64
			var firstName string
			var username string

			if update.Message != nil && update.Message.From != nil {
				userId = update.Message.From.ID
				firstName = update.Message.From.FirstName
				username = fmt.Sprintf("@%s", update.Message.From.Username)
			} else if update.CallbackQuery != nil {
				userId = update.CallbackQuery.From.ID
				firstName = update.CallbackQuery.From.FirstName
				username = fmt.Sprintf("@%s", update.CallbackQuery.From.Username)
			} else if update.InlineQuery != nil && update.InlineQuery.From != nil {
				userId = update.InlineQuery.From.ID
				firstName = update.InlineQuery.From.FirstName
				username = fmt.Sprintf("@%s", update.InlineQuery.From.Username)
			}

			if userId != 0 {
				user := &models.User{
					UserId:    userId,
					FirstName: utils.RemoveHTMLTags(firstName),
					Username:  username,
				}

				err := userRepo.UpsertUser(ctx, user)
				if err != nil {
					log.Printf("❌ Erro ao upsert do usuário: %v", err)
				}
			}

			next(ctx, b, update)
		}
	}
}
