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

		var chatID int64
		var messageID int
		var inlineMessageID string
		var replyMarkup *models.InlineKeyboardMarkup

		if update.CallbackQuery.Message.Message != nil {
			msg := update.CallbackQuery.Message.Message
			chatID = msg.Chat.ID
			messageID = msg.ID
			replyMarkup = msg.ReplyMarkup
		} else {
			inlineMessageID = update.CallbackQuery.InlineMessageID
		}

		if chatID == 0 && messageID == 0 && inlineMessageID == "" {
			logger.Warn("VOTE", "Recebido callback de voto sem identificador de mensagem: %v", update.CallbackQuery.ID)
			return
		}

		// 1. Registrar/Alternar voto no banco de dados
		added, _, err := c.VoteService.ToggleVote(ctx, chatID, messageID, inlineMessageID, update.CallbackQuery.From.ID, votedEmoji)
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
		counts, err := c.VoteService.GetVoteCounts(ctx, chatID, messageID, inlineMessageID)
		if err != nil {
			logger.Error("VOTE", "Erro ao buscar contagens: %v", err)
			return
		}

		// Se for mensagem normal, usamos o ReplyMarkup que veio nela.
		// Se for inline, tentamos recuperar o sessionID mapeado para reconstruir o teclado.
		ikb := replyMarkup
		updated := false

		if ikb == nil && inlineMessageID != "" {
			var sessionID string
			key := fmt.Sprintf("pb_inline_map:%s", inlineMessageID)
			err := c.CacheService.Get(ctx, key, &sessionID)
			if err != nil || sessionID == "" {
				logger.Warn("VOTE", "❌ Mapeamento inline não encontrado no Redis. Chave: %s, Erro: %v", key, err)
			} else {
				state, _ := c.CacheService.GetPostBuilderSession(ctx, sessionID)
				if state != nil {
					ikb = &models.InlineKeyboardMarkup{}
					// Botões de URL
					for _, btn := range state.Buttons {
						ikb.InlineKeyboard = append(ikb.InlineKeyboard, []models.InlineKeyboardButton{
							{Text: btn.Text, URL: btn.URL},
						})
					}
					// Reações (já reconstrói com a contagem atual)
					if state.Reactions != "" {
						reactions := strings.Split(state.Reactions, ",")
						var reactionRow []models.InlineKeyboardButton
						for _, r := range reactions {
							emoji := strings.TrimSpace(r)
							if emoji == "" {
								continue
							}

							count := counts[emoji]
							text := emoji
							if count > 0 {
								text = fmt.Sprintf("%s %d", emoji, count)
							}

							reactionRow = append(reactionRow, models.InlineKeyboardButton{
								Text:         text,
								CallbackData: "vote:" + emoji,
							})
						}
						if len(reactionRow) > 0 {
							ikb.InlineKeyboard = append(ikb.InlineKeyboard, reactionRow)
						}
					}
					logger.Bot("✅ Teclado inline reconstruído para sessão %s", sessionID)
					updated = true // Força atualização pois o teclado original do Telegram não tem contagem
				}
			}
		}

		if ikb == nil {
			if inlineMessageID == "" {
				logger.Warn("VOTE", "Mensagem sem teclado inline: %d", messageID)
			}
			return
		}

		// Atualizar o teclado (seja o original ou o reconstruído)
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
			editParams := &bot.EditMessageReplyMarkupParams{
				ReplyMarkup: ikb,
			}
			if inlineMessageID != "" {
				editParams.InlineMessageID = inlineMessageID
			} else {
				editParams.ChatID = chatID
				editParams.MessageID = messageID
			}

			_, err := b.EditMessageReplyMarkup(ctx, editParams)
			if err != nil {
				errStr := err.Error()
				isNotModified := strings.Contains(errStr, "message is not modified")
				isUnmarshalBool := strings.Contains(errStr, "cannot unmarshal bool")

				if !isNotModified && !isUnmarshalBool {
					logger.Error("VOTE", "Erro ao editar teclado: %v", err)
				}
			}
		}
	}
}
