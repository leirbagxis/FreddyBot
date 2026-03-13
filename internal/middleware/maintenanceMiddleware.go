package middleware

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"gorm.io/gorm"
)

func CheckMaintenceMiddleware(db *gorm.DB) bot.Middleware {
	serverRepo := repositories.NewServerConfigRepository(db)
	userRepo := repositories.NewUserRepository(db)

	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, upt *models.Update) {
			if upt.ChannelPost != nil {
				next(ctx, b, upt)
				return
			}

			userID := getUpdateUserID(upt)
			if userID == 0 {
				return
			}

			user, err := userRepo.GetUserById(ctx, userID)
			if err != nil {
				fmt.Println(err)
				return
			}

			if user.IsAdmin {
				next(ctx, b, upt)
				return
			}

			maintence, err := serverRepo.GetMaintence(ctx)
			if err != nil {
				fmt.Printf("erro ao pegar o maintence: %v\n", err)
				next(ctx, b, upt)
				return
			}

			if !maintence {
				next(ctx, b, upt)
				return
			}

			sendMaintenceResponse(ctx, b, upt, userID)
		}
	}
}

func getUpdateUserID(upt *models.Update) int64 {
	switch {
	case upt.Message != nil:
		return upt.Message.From.ID
	case upt.CallbackQuery != nil:
		return upt.CallbackQuery.From.ID
	default:
		return 0
	}
}

func sendMaintenceResponse(ctx context.Context, b *bot.Bot, upt *models.Update, userID int64) {
	switch {
	case upt.Message != nil:
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      "⚙️ <b>Manutenção em andamento</b>\n\nO bot está temporariamente indisponível para melhorias no sistema.\n\n⏳ <i>Voltaremos em breve.</i>\n\nObrigado pela sua paciência 💙",
			ParseMode: models.ParseModeHTML,
			ReplyParameters: &models.ReplyParameters{
				MessageID: upt.Message.ID,
			},
		})

	case upt.CallbackQuery != nil:
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: upt.CallbackQuery.ID,
			Text:            "🚧 O bot está em manutenção no momento.",
			ShowAlert:       true,
			CacheTime:       0,
		})

	case upt.ChannelPost != nil:
		// Se quiser bloquear post de canal silenciosamente, não faz nada.
		return
	}
}
