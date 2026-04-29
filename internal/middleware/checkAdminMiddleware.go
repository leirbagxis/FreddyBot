package middleware

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func CheckAdminMiddleware(app *container.AppContainer) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			var userID int64
			ownerID := config.OwnerID

			if update.Message != nil && update.Message.From != nil {
				userID = update.Message.From.ID
			} else if update.CallbackQuery != nil {
				userID = update.CallbackQuery.From.ID
			} else if update.InlineQuery != nil && update.InlineQuery.From != nil {
				userID = update.InlineQuery.From.ID
			}

			user, err := app.UserRepo.GetUserById(ctx, userID)
			if err != nil || user == nil {
				if userID != ownerID {
					fmt.Println("User not found or not admin: ", userID)
					return
				}
				// If owner but not in DB, we still want to allow
			} else if !user.IsAdmin && user.UserId != ownerID {
				fmt.Println("caindo aqui: ", !user.IsAdmin || user.UserId != ownerID)
				return
			}

			next(ctx, b, update)
		}
	}
}
