package addchannel

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
	"github.com/leirbagxis/FreddyBot/internal/telegram/modules"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskAddChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Erro ao adicionar-chatmember: %v", r)
				//sendErrorMessage(ctx, b, update)
			}
		}()

		botInfo, _ := b.GetMe(ctx)
		chat := update.MyChatMember.Chat
		from := update.MyChatMember.From

		if update.MyChatMember.OldChatMember.Type != "left" {
			return
		}

		getChannel, _ := c.ChannelRepo.GetChannelByID(ctx, chat.ID)

		data := map[string]string{
			"channelName": utils.RemoveHTMLTags(chat.Title),
			"channelId":   fmt.Sprintf("%d", chat.ID),
			"botId":       fmt.Sprintf("%d", botInfo.ID),
			"firstName":   utils.RemoveHTMLTags(from.FirstName),
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
	}
}

func AskForwadedChannelHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Erro ao adicionar-forward: %v", r)
				//sendErrorMessage(ctx, b, update)
			}
		}()

		botInfo, _ := b.GetMe(ctx)
		from := update.Message.From
		forwardedChannelID := update.Message.ForwardOrigin.MessageOriginChannel.Chat.ID
		forwardedChannelTitle := update.Message.ForwardOrigin.MessageOriginChannel.Chat.Title

		getChannel, _ := c.ChannelRepo.GetChannelByID(ctx, forwardedChannelID)

		data := map[string]string{
			"channelName": utils.RemoveHTMLTags(forwardedChannelTitle),
			"channelId":   fmt.Sprintf("%d", forwardedChannelID),
			"botId":       fmt.Sprintf("%d", botInfo.ID),
			"firstName":   utils.RemoveHTMLTags(from.FirstName),
		}

		if getChannel != nil {
			text, button := parser.GetMessage("toadd-exist-channel", data)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      from.ID,
				Text:        text,
				ReplyMarkup: button,
				ParseMode:   "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		session, err := c.SessionManager.CreateChannelSession(ctx, forwardedChannelID, from.ID, forwardedChannelTitle)
		if err != nil {
			log.Printf("Erro ao criar sessão: %v", err)
			return
		}
		data["sessionKey"] = session.Key

		countChannel, _ := c.ChannelRepo.CountUserChannels(ctx, from.ID)
		if countChannel >= 1 {
			modules.SendChannelActivationPayment(ctx, b, update, c, session.Key)
			return
		}

		text, button := parser.GetMessage("toadd-require-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      from.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})

		fmt.Println(from, forwardedChannelID)

	}
}

func AddYesHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		botInfo, _ := b.GetMe(ctx)
		callback := update.CallbackQuery
		from := update.CallbackQuery.From

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
        ˗ˏˋ [$title]($link) ⋆｡˚ ☁︎
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

		userID := fmt.Sprintf("%d", from.ID)
		channelID := fmt.Sprintf("%d", channel.ID)

		c.SessionManager.DeleteChannelSession(ctx, parts[1])

		data := map[string]string{
			"firstName":   utils.RemoveHTMLTags(from.FirstName),
			"botId":       fmt.Sprintf("%d", botInfo.ID),
			"channelName": utils.RemoveHTMLTags(channel.Title),
			"channelId":   fmt.Sprintf("%d", channel.ID),
			"miniAppUrl":  auth.GenerateMiniAppUrl(userID, channelID),
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

		// CORREÇÃO: Adicionar reação de forma correta
		reactionParams := &bot.SetMessageReactionParams{
			ChatID:    editMsg.Chat.ID,
			MessageID: editMsg.ID,
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
		}

		// Tentar adicionar reação (não crítico se falhar)
		_, err = b.SetMessageReaction(ctx, reactionParams)
		if err != nil {
			log.Printf("Aviso: Não foi possível adicionar reação: %v", err)
			// Não retornar erro, apenas logar
		}

		// Responder ao callback
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "✅ Canal adicionado com sucesso!",
		})

		logs.LogAdmin(ctx, b, channel)

	}
}

func AddNotHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callback := update.CallbackQuery

		callbackData := update.CallbackQuery.Data
		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		_, err := c.SessionManager.GetChannelSession(ctx, parts[1])
		if err != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "⌛ Tempo Esgotado. Faça o processo novamente!",
				ShowAlert:       true,
			})
			return
		}

		c.SessionManager.DeleteChannelSession(ctx, parts[1])

		data := map[string]string{}

		text, button := parser.GetMessage("toadd-cancel-message", data)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
		})

		// Responder ao callback
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "❌ Operação Cancelada!",
		})

	}
}
