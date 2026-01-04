package modules

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/telegram/logs"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func SendChannelActivationPayment(
	ctx context.Context,
	b *bot.Bot,
	userID int64,
	messageID int,
	sessionKey,
	description string,
) {
	b.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:      userID,
		Title:       "🔓 Liberar cadastro",
		Description: description,
		Payload:     "active-pay-channel:" + sessionKey,
		Currency:    "XTR",
		Prices: []models.LabeledPrice{
			{
				Label:  "Ativação do canal",
				Amount: 1,
			},
		},
		ReplyParameters: &models.ReplyParameters{
			MessageID: messageID,
		},
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "⭐ Pagar 1 Stars",
						Pay:  true,
					},
				},
				{
					{
						Text:         "Cancelar",
						CallbackData: "add-not:" + sessionKey,
					},
				},
			},
		},
	})
}

func SuccessfulPaymentChannel(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.SuccessfulPayment == nil {
			return
		}

		botInfo, _ := b.GetMe(ctx)

		msg := update.Message
		payment := msg.SuccessfulPayment
		payload := payment.InvoicePayload
		from := msg.From

		parts := strings.Split(payload, ":")
		if len(parts) != 2 {
			log.Println("Payload inválido:", payload)
			return
		}

		sessionKey := parts[1]

		getSession, err := c.SessionManager.GetChannelSession(ctx, sessionKey)
		if err != nil || getSession == nil {
			log.Printf(
				"Pagamento recebido com sessão inválida. Session=%s User=%d",
				sessionKey,
				from.ID,
			)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text: "⚠️ O pagamento foi confirmado, mas a solicitação expirou.\n\n" +
					"Entre em contato com o suporte informando o ID abaixo:\n\n" +
					"🆔 " + sessionKey,
			})
			return
		}

		// 🧹 Remove a sessão (evita replay)
		_ = c.SessionManager.DeleteChannelSession(ctx, sessionKey)

		channelInfo, err := b.GetChat(ctx, &bot.GetChatParams{
			ChatID: getSession.ChannelID,
		})
		if err != nil {
			log.Printf("Erro ao obter info do canal: %v", err)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   "❌ Erro ao obter informações do canal. Fale com o suporte.",
			})
			return
		}

		inviteURL := channelInfo.InviteLink
		if channelInfo.Username != "" {
			inviteURL = fmt.Sprintf("https://t.me/%s", channelInfo.Username)
		}

		newPackCaption := fmt.Sprintf(`╔═━──━═༻✧༺═━──━═╗

𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
(｡•́︿•̀｡)っ✧.*ೃ༄
˗ˏˋ [%s](%s) ⋆｡˚ ☁︎
    彡♡ ₊˚

⋆｡˚ ❀ @%s ☽⁺₊

╚═━──━═༻✧༺═━──━═╝`,
			getSession.Title,
			inviteURL,
			botInfo.Username,
		)

		defaultCaption := fmt.Sprintf("➽ 𝐛𝐲 @%s", botInfo.Username)

		channel, err := c.ChannelRepo.CreateChannelWithDefaults(
			ctx,
			getSession.ChannelID,
			getSession.Title,
			inviteURL,
			newPackCaption,
			defaultCaption,
			getSession.OwnerID,
		)
		if err != nil {
			log.Printf("Erro ao criar canal: %v", err)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   "❌ Erro ao finalizar a criação do canal. Fale com o suporte.",
			})
			return
		}

		userID := fmt.Sprintf("%d", from.ID)
		channelID := fmt.Sprintf("%d", channel.ID)

		data := map[string]string{
			"firstName":   utils.RemoveHTMLTags(from.FirstName),
			"botId":       fmt.Sprintf("%d", botInfo.ID),
			"channelName": utils.RemoveHTMLTags(channel.Title),
			"channelId":   channelID,
			"miniAppUrl":  auth.GenerateMiniAppUrl(userID, channelID),
		}

		text, button := parser.GetMessage("toadd-success-message", data)

		// ✅ Mensagem final de sucesso
		sentMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      msg.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
		})
		if err != nil {
			log.Printf("Erro ao enviar mensagem final: %v", err)
			return
		}

		// 🎉 Reação (não crítica)
		_, err = b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
			ChatID:    sentMsg.Chat.ID,
			MessageID: sentMsg.ID,
			Reaction: []models.ReactionType{
				{
					Type: "emoji",
					ReactionTypeEmoji: &models.ReactionTypeEmoji{
						Type:  "emoji",
						Emoji: "🎉",
					},
				},
			},
			IsBig: bot.True(),
		})
		if err != nil {
			log.Printf("Aviso: não foi possível adicionar reação: %v", err)
		}

		logs.LogAdmin(ctx, b, channel)
	}
}

// func SuccessfulPaymentChannel(c *container.AppContainer) bot.HandlerFunc {
// 	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
// 		botInfo, _ := b.GetMe(ctx)
// 		msg := update.Message
// 		payment := msg.SuccessfulPayment
// 		payload := payment.InvoicePayload
// 		from := msg.From

// 		parts := strings.Split(payload, ":")
// 		if len(parts) != 2 {
// 			log.Println("Callback invalido:", payload)
// 			return
// 		}

// 		getSession, err := c.SessionManager.GetChannelSession(ctx, parts[1])
// 		if err != nil || getSession == nil {
// 			log.Printf("Pagamento recebido com sessão inválida. Payload=%s User=%d",
// 				parts[1],
// 				msg.From.ID,
// 			)

// 			b.SendMessage(ctx, &bot.SendMessageParams{
// 				ChatID: msg.Chat.ID,
// 				Text: "⚠️ O pagamento foi confirmado, mas a solicitação expirou.\n" +
// 					"Entre em contato com o suporte para regularizar.\n\n" +
// 					"ID: " + parts[1],
// 			})

// 			return
// 		}

// 		// 🧹 Limpa a sessão (evita replay)
// 		_ = c.SessionManager.DeleteChannelSession(ctx, parts[1])

// 		channelInfo, err := b.GetChat(ctx, &bot.GetChatParams{
// 			ChatID: getSession.ChannelID,
// 		})
// 		if err != nil {
// 			log.Printf("Erro ao obter info do canal: %v", err)
// 			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
// 				CallbackQueryID: update.CallbackQuery.ID,
// 				Text:            "❌ Erro ao obter informações do canal!",
// 				ShowAlert:       true,
// 			})
// 			return
// 		}

// 		inviteURL := channelInfo.InviteLink
// 		if channelInfo.Username != "" {
// 			inviteURL = fmt.Sprintf("t.me/%s", channelInfo.Username)
// 		}

// 		// Gerar newPackCaption
// 		newPackCaption := fmt.Sprintf(`╔═━──━═༻✧༺═━──━═╗

//         𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
//         (｡•́︿•̀｡)っ✧.*ೃ༄
//         ˗ˏˋ [$title]($link) ⋆｡˚ ☁︎
//             彡♡ ₊˚

// ⋆｡˚ ❀ @%s ☽⁺₊

// ╚═━──━═༻✧༺═━──━═╝`, botInfo.Username)

// 		defaultCaption := fmt.Sprintf("➽ 𝐛𝐲 @%s", botInfo.Username)

// 		channel, err := c.ChannelRepo.CreateChannelWithDefaults(
// 			ctx,
// 			getSession.ChannelID,
// 			getSession.Title,
// 			inviteURL,
// 			newPackCaption,
// 			defaultCaption,
// 			getSession.OwnerID,
// 		)
// 		if err != nil {
// 			log.Printf("Erro ao criar canal: %v", err)
// 			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
// 				CallbackQueryID: callback.ID,
// 				Text:            "❌ Erro ao criar canal!",
// 				ShowAlert:       true,
// 			})
// 			return
// 		}

// 		userID := fmt.Sprintf("%d", from.ID)
// 		channelID := fmt.Sprintf("%d", channel.ID)

// 		data := map[string]string{
// 			"firstName":   utils.RemoveHTMLTags(from.FirstName),
// 			"botId":       fmt.Sprintf("%d", botInfo.ID),
// 			"channelName": utils.RemoveHTMLTags(channel.Title),
// 			"channelId":   fmt.Sprintf("%d", channel.ID),
// 			"miniAppUrl":  auth.GenerateMiniAppUrl(userID, channelID),
// 		}

// 		text, button := parser.GetMessage("toadd-success-message", data)

// 		editMsg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
// 			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
// 			MessageID:   update.CallbackQuery.Message.Message.ID,
// 			Text:        text,
// 			ReplyMarkup: button,
// 			ParseMode:   "HTML",
// 		})
// 		if err != nil {
// 			log.Printf("Erro ao editar mensagem: %v", err)
// 			return
// 		}

// 		// CORREÇÃO: Adicionar reação de forma correta
// 		reactionParams := &bot.SetMessageReactionParams{
// 			ChatID:    editMsg.Chat.ID,
// 			MessageID: editMsg.ID,
// 			Reaction: []models.ReactionType{
// 				{
// 					Type: "emoji",
// 					ReactionTypeEmoji: &models.ReactionTypeEmoji{
// 						Type:  "emoji",
// 						Emoji: "🎉",
// 					},
// 				},
// 			},
// 			IsBig: bot.True(),
// 		}

// 		// Tentar adicionar reação (não crítico se falhar)
// 		_, err = b.SetMessageReaction(ctx, reactionParams)
// 		if err != nil {
// 			log.Printf("Aviso: Não foi possível adicionar reação: %v", err)
// 			// Não retornar erro, apenas logar
// 		}

// 		// Responder ao callback
// 		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
// 			CallbackQueryID: callback.ID,
// 			Text:            "✅ Canal adicionado com sucesso!",
// 		})

// 		logs.LogAdmin(ctx, b, channel)
// 	}
// }
