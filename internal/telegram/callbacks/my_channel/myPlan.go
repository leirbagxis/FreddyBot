package mychannel

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func MyPlanHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		cbks := update.CallbackQuery

		userId := cbks.From.ID
		session, err := c.CacheService.GetSelectedChannel(ctx, userId)
		if err != nil {
			log.Printf("Erro ao pegar sessão: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
				ShowAlert:       true,
			})
			return
		}
		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, session)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}
		subs, err := c.SubscriptionRepo.GetChannelSubscription(ctx, channel.ID)
		if err != nil {
			log.Printf("Erro ao buscar assinatura: %v", err)
			return
		}

		data := map[string]string{
			"title":         channel.Title,
			"channelId":     fmt.Sprintf("%d", session),
			"planName":      subs.Plan.Name,
			"planValidity":  subs.EndDate.Format("2006-01-02 15:04:05"),
			"planStartDate": subs.StartDate.Format("2006-01-02 15:04:05"),
			"planId":        subs.ID,
		}
		text, button := parser.GetMessage("my_plan", data)

		// err = c.CacheService.SetSelectedChannel(ctx, session, session)
		// if err != nil {
		// 	log.Printf("Erro ao criar sessão: %v", err)
		// 	return
		// }

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

	}
}
