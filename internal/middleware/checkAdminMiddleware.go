package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
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
					logger.Error("MID", "Acesso negado: Usuário não encontrado ou não admin: %d", userID)
					return
				}
				// If owner but not in DB, we still want to allow
			} else if !user.IsAdmin && user.UserId != ownerID {
				logger.Error("MID", "Acesso negado: Usuário %d tentou comando de admin", userID)
				return
			}

			next(ctx, b, update)
		}
	}
}
