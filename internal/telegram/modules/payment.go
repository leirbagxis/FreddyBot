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

	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

func SendChannelActivationPayment(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	c *container.AppContainer,
	sessionKey string,
) {

	getSession, err := c.SessionManager.GetChannelSession(ctx, sessionKey)
	if err != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "⌛ Tempo Esgotado. Faça o processo novamente!",
			ShowAlert:       true,
		})
		return
	}

	getPrice, _ := c.PaymentService.GetPricePlan(ctx, "add_to_channel_fee")
	price := getPrice.Amount

	payment := dbmodels.Payment{
		UserID:  getSession.OwnerID,
		Amount:  int(price),
		Payload: sessionKey,
	}
	_, err = c.PaymentService.CreateNewPayment(ctx, payment)
	if err != nil {
		log.Fatal("Erro ao criar pagamento: %w", err.Error())
		return
	}

	user, err := b.GetChat(ctx, &bot.GetChatParams{
		ChatID: getSession.OwnerID,
	})
	if err != nil {
		log.Printf("Erro ao buscar usuário: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: getSession.OwnerID, Text: "❌ Erro ao buscar informações do usuário."})
		return
	}

	data := map[string]string{
		"firstName": user.FirstName,
		"value":     fmt.Sprintf("%d", price),
	}

	text, button := parser.GetMessage("toadd-payment-require-message", data)
	b.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:      getSession.OwnerID,
		Title:       "🔓 Liberar cadastro",
		Description: text,
		Payload:     "active-pay-channel:" + sessionKey,
		Currency:    "XTR",
		Prices: []models.LabeledPrice{
			{
				Label:  "Ativação do canal",
				Amount: int(price),
			},
		},
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
		ReplyMarkup: button,
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
