package addchannel

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chat := update.MyChatMember.Chat
		from := update.MyChatMember.From

		if update.MyChatMember.OldChatMember.Type != "left" {
			return
		}

		getChannel, _ := c.ChannelRepo.GetChannelByID(ctx, chat.ID)

		data := map[string]string{
			"channelName": chat.Title,
		}

		if getChannel != nil {
			text, button := parser.GetMessage("toadd-exist-channel", data)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      from.ID,
				Text:        text,
				ReplyMarkup: button,
				ParseMode:   "HTML",
			})
			return
		}

		session, err := c.SessionManager.CreateChannelSession(ctx, chat.ID, from.ID, chat.Title)
		if err != nil {
			log.Printf("Erro ao criar sessão: %v", err)
			return
		}
		data["sessionKey"] = session.Key

		text, button := parser.GetMessage("toadd-require-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      from.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
		})
		fmt.Println(from)
	}
}

func AddYesHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		botInfo, _ := b.GetMe(ctx)
		callback := update.CallbackQuery
		from := update.CallbackQuery.From
		userID := from.ID

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		getSession, err := c.SessionManager.GetChannelSession(ctx, parts[1])
		if err != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Tempo Esgotado. Faça o processo novamente!",
				ShowAlert:       true,
			})
			return
		}

		channelInfo, err := b.GetChat(ctx, &bot.GetChatParams{
			ChatID: getSession.ChannelID,
		})
		if err != nil {
			log.Printf("Erro ao obter info do canal: %v", err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao obter informações do canal!",
				ShowAlert:       true,
			})
			return
		}

		inviteURL := channelInfo.InviteLink
		if channelInfo.Username != "" {
			inviteURL = fmt.Sprintf("t.me/%s", channelInfo.Username)
		}

		// Gerar newPackCaption
		newPackCaption := fmt.Sprintf(`╔═━──━═༻✧༺═━──━═╗

        𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
        (｡•́︿•̀｡)っ✧.*ೃ༄
        ˗ˏˋ [$name]($link) ⋆｡˚ ☁︎
            彡♡ ₊˚

⋆｡˚ ❀ @%s ☽⁺₊

╚═━──━═༻✧༺═━──━═╝`, botInfo.Username)

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
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: callback.ID,
				Text:            "❌ Erro ao criar canal!",
				ShowAlert:       true,
			})
			return
		}

		c.SessionManager.DeleteChannelSession(ctx, parts[1])

		data := map[string]string{
			"firstName":   from.FirstName,
			"botId":       fmt.Sprintf("%s", botInfo.ID),
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", channel.ID),
			"miniAppUrl":  "https://caption.chelodev.shop/6762185696/-1001765135605?signature=53b8be8058f96458794c406e0f31fe91bb43e1a9cac2ed9e6f4e8b87efeccb86",
		}

		text, button := parser.GetMessage("toadd-success-message", data)

		editMsg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
		})
		if err != nil {
			log.Printf("Erro ao editar mensagem: %v", err)
			return
		}

		b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
			ChatID:    editMsg.Chat.ID,
			MessageID: editMsg.ID,
			Reaction: []models.ReactionType{
				{
					Type: "emoji",
					ReactionTypeCustomEmoji: &models.ReactionTypeCustomEmoji{
						CustomEmojiID: "🎉",
					},
				},
			},
		})

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "✅ Canal adicionado com sucesso!",
		})

		fmt.Println(userID, getSession)

	}
}
