package middleware

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func CheckBlacklistMiddleware(c *container.AppContainer) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, upt *models.Update) {
			// Ignorar postagens de canal aqui, pois elas são tratadas no handler de channelPost
			if upt.ChannelPost != nil || upt.EditedChannelPost != nil {
				next(ctx, b, upt)
				return
			}

			userID := getUpdateUserID(upt)
			if userID == 0 {
				next(ctx, b, upt)
				return
			}

			// Owner nunca é bloqueado
			if userID == config.OwnerID {
				next(ctx, b, upt)
				return
			}

			user, err := c.UserService.GetUserByID(ctx, userID)
			if err != nil || user == nil {
				next(ctx, b, upt)
				return
			}

			if !user.IsBlacklisted {
				next(ctx, b, upt)
				return
			}

			// Se está na blacklist, apenas /start e /ouvidoria são permitidos
			if upt.Message != nil && (strings.HasPrefix(upt.Message.Text, "/start") || strings.HasPrefix(upt.Message.Text, "/ouvidoria")) {
				next(ctx, b, upt)
				return
			}

			// Bloquear outras interações
			if upt.CallbackQuery != nil {
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: upt.CallbackQuery.ID,
					Text:            "❌ Você está na blacklist e seus comandos estão bloqueados.",
					ShowAlert:       true,
				})
				return
			}

			// Para outras mensagens, podemos enviar um aviso silencioso ou apenas ignorar
		}
	}
}

