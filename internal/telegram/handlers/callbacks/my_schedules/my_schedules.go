package myschedules

import (
	"context"
	"fmt"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func HandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		userID := update.CallbackQuery.From.ID
		bot := ctx.Bot()

		schedules, err := c.SchedulerService.GetUserSchedules(context.Background(), userID)
		if err != nil {
			logger.Error("BOT", "Erro ao buscar agendamentos: %v", err)
			_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao buscar agendamentos.",
			})
			return nil
		}

		if len(schedules) == 0 {
			_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "📋 Nenhum agendamento encontrado.",
			})
			return nil
		}

		var text string
		text = "📅 <b>Meus Agendamentos</b>\n\n"

		var buttons [][]telego.InlineKeyboardButton
		for i, s := range schedules {
			statusEmoji := "🟢"
			if s.Status == "paused" {
				statusEmoji = "🟡"
			} else if s.Status == "cancelled" {
				statusEmoji = "🔴"
			} else if s.Status == "completed" {
				statusEmoji = "✅"
			}

			typeLabel := s.ScheduleType
			timeLabel := s.ScheduleTime
			if s.ScheduleType == "interval" {
				typeLabel = fmt.Sprintf("intervalo (%dmin)", s.IntervalMin)
				if s.WindowStart != "" {
					timeLabel = fmt.Sprintf("%s-%s", s.WindowStart, s.WindowEnd)
				} else {
					timeLabel = "24h"
				}
			}

			text += fmt.Sprintf("%s <b>%s</b> → %s\n⏰ %s | %s\n\n",
				statusEmoji,
				s.ChannelTitle,
				typeLabel,
				timeLabel,
				s.Status,
			)

			buttons = append(buttons, []telego.InlineKeyboardButton{
				{Text: fmt.Sprintf("📋 %s", s.ChannelTitle), CallbackData: fmt.Sprintf("schedule-detail:%s", s.ID)},
			})

			_ = i
		}

		buttons = append(buttons, []telego.InlineKeyboardButton{
			{Text: "🔙 Voltar", CallbackData: "start"},
		})

		kb := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
			MessageID:   update.CallbackQuery.Message.GetMessageID(),
			Text:        text,
			ParseMode:   telego.ModeHTML,
			ReplyMarkup: kb,
		})

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		return nil
	}
}
