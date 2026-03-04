package tutorial

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.From.ID

		topic, err := b.CreateForumTopic(ctx, &bot.CreateForumTopicParams{
			ChatID:            chatID,
			Name:              "📚 Tutoriais",
			IconCustomEmojiID: "5334882760735598374",
		})
		fmt.Println(topic, err)
		if err != nil {
			// Se o usuário não tiver tópicos habilitados pro bot, pode dar erro aqui.
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Não consegui criar a aba de tutoriais aqui no PV. Verifique se tópicos no privado estão habilitados para este bot.",
			})
			return
		}

		threadID := topic.MessageThreadID

		// 2) envia um texto dentro do tópico
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          chatID,
			MessageThreadID: threadID,
			Text:            "Bem-vindo ao tutorial! Vou te mandar os vídeos aqui 👇",
		})

	}
}
