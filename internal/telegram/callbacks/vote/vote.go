package vote

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}

		data := update.CallbackQuery.Data
		if !strings.HasPrefix(data, "vote:") {
			return
		}

		// Extrair o emoji do callback data
		votedEmoji := strings.TrimPrefix(data, "vote:")

		msg := update.CallbackQuery.Message
		if msg.Message == nil {
			logger.Warn("VOTE", "Recebido callback de voto sem mensagem ou mensagem inacessível: %v", update.CallbackQuery.ID)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Erro: Mensagem não encontrada.",
			})
			return
		}

		// 1. Registrar/Alternar voto no banco de dados
		added, _, err := c.VoteRepo.ToggleVote(ctx, msg.Message.Chat.ID, msg.Message.ID, update.CallbackQuery.From.ID, votedEmoji)
		if err != nil {
			logger.Error("VOTE", "Erro ao processar voto no banco: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Erro ao computar voto.",
			})
			return
		}

		// Feedback visual para o usuário
		feedback := fmt.Sprintf("Voto removido de %s", votedEmoji)
		if added {
			feedback = fmt.Sprintf("Você votou em %s!", votedEmoji)
		}
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            feedback,
		})

		// 2. Buscar contagens atualizadas para todos os emojis desta mensagem
		counts, err := c.VoteRepo.GetVoteCounts(ctx, msg.Message.Chat.ID, msg.Message.ID)
		if err != nil {
			logger.Error("VOTE", "Erro ao buscar contagens: %v", err)
			return
		}

		ikb := msg.Message.ReplyMarkup
		if ikb == nil {
			logger.Warn("VOTE", "Mensagem sem teclado inline: %d", msg.Message.ID)
			return
		}

		// 3. Atualizar todos os botões de voto com as contagens reais do banco
		updated := false
		for i, row := range ikb.InlineKeyboard {
			for j, btn := range row {
				if strings.HasPrefix(btn.CallbackData, "vote:") {
					emoji := strings.TrimPrefix(btn.CallbackData, "vote:")
					count := counts[emoji]

					var newText string
					if count > 0 {
						newText = fmt.Sprintf("%s %d", emoji, count)
					} else {
						newText = emoji
					}

					if ikb.InlineKeyboard[i][j].Text != newText {
						ikb.InlineKeyboard[i][j].Text = newText
						updated = true
					}
				}
			}
		}

		if updated {
			// Editar a mensagem original com o teclado atualizado
			_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
				ChatID:      msg.Message.Chat.ID,
				MessageID:   msg.Message.ID,
				ReplyMarkup: ikb,
			})
			if err != nil {
				// Ignorar erro "message is not modified" que pode ocorrer em condições de corrida
				if !strings.Contains(err.Error(), "message is not modified") {
					logger.Error("VOTE", "Erro ao editar teclado: %v", err)
				}
			} else {
				logger.Bot("Voto atualizado para %s na mensagem %d (Chat: %d)", votedEmoji, msg.Message.ID, msg.Message.Chat.ID)
			}
		}
	}
}
