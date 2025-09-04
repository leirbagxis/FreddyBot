package mychannel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

const layoutBR = "02/01/2006 as 15:04:05"

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

		var validity string
		if subs.EndDate.IsZero() {
			validity = "Vitalício"
		} else {
			validity = subs.EndDate.Format(layoutBR)
		}

		data := map[string]string{
			"title":         channel.Title,
			"channelId":     fmt.Sprintf("%d", session),
			"planName":      subs.Plan.Name,
			"planValidity":  validity,
			"planStartDate": subs.StartDate.Format(layoutBR),
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

func MyPlanResouces(c *container.AppContainer) bot.HandlerFunc {
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
			"planStartDate": subs.StartDate.Format(layoutBR),
			"planId":        subs.ID,
		}
		text, button := parser.GetMessage("premium_resources", data)

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

func MyPlanAskCaption(c *container.AppContainer) bot.HandlerFunc {
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
		_, err = c.SubscriptionRepo.GetChannelSubscription(ctx, channel.ID)
		if err != nil {
			log.Printf("Erro ao buscar assinatura: %v", err)
			return
		}

		data := map[string]string{}
		text, button := parser.GetMessage("ask_plan_caption", data)

		err = c.SessionManager.SetPlainCaptionSession(ctx, userId, session)
		if err != nil {
			log.Printf("Erro ao criar sessão - SetPlainCaptionSession: %v", err)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})

	}
}

func MyPlanFoundCaption(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		res, _ := json.Marshal(update)
		fmt.Println(string(res))
		// cbks := update.CallbackQuery

		// userId := cbks.From.ID
		// session, err := c.CacheService.GetSelectedChannel(ctx, userId)
		// if err != nil {
		// 	log.Printf("Erro ao pegar sessão: %v", err)
		// 	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		// 		CallbackQueryID: update.CallbackQuery.ID,
		// 		Text:            "⌛ Seção Expirada. Selecione o canal novamente!",
		// 		ShowAlert:       true,
		// 	})
		// 	return
		// }
		// channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, session)
		// if err != nil {
		// 	log.Printf("Erro ao buscar canal: %v", err)
		// 	return
		// }
		// _, err = c.SubscriptionRepo.GetChannelSubscription(ctx, channel.ID)
		// if err != nil {
		// 	log.Printf("Erro ao buscar assinatura: %v", err)
		// 	return
		// }

		// data := map[string]string{}
		// text, button := parser.GetMessage("ask_plan_caption", data)

		// // err = c.CacheService.SetSelectedChannel(ctx, session, session)
		// // if err != nil {
		// // 	log.Printf("Erro ao criar sessão: %v", err)
		// // 	return
		// // }

		// b.EditMessageText(ctx, &bot.EditMessageTextParams{
		// 	ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		// 	Text:        text,
		// 	ReplyMarkup: button,
		// 	ParseMode:   "HTML",
		// 	MessageID:   update.CallbackQuery.Message.Message.ID,
		// })

	}
}
