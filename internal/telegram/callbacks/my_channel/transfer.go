package mychannel

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func AskTransferAccessHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callbackData := update.CallbackQuery.Data

		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		userId := update.CallbackQuery.From.ID
		channelId, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		data := map[string]string{
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", channelId),
		}

		text, button := parser.GetMessage("ask-paccess-message", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func TransferAcessHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callbackData := update.CallbackQuery.Data

		parts := strings.Split(callbackData, ":")
		if len(parts) != 2 {
			log.Println("Callback invalido:", callbackData)
			return
		}

		userId := update.CallbackQuery.From.ID
		channelId, _ := strconv.ParseInt(parts[1], 10, 64)

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			return
		}

		user, err := c.UserRepo.GetUserById(ctx, userId)
		if err != nil {
			log.Printf("Erro ao buscar usuario: %v", err)
			return
		}

		data := map[string]string{
			"channelName": channel.Title,
			"channelId":   fmt.Sprintf("%d", channelId),
			"ownerId":     fmt.Sprintf("%d", user.UserId),
			"ownerName":   user.FirstName,
		}
		c.CacheService.SetAwaitingStickerSeparator(ctx, userId, channelId)

		text, button := parser.GetMessage("require-paccess-message", data)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			MessageID:   update.CallbackQuery.Message.Message.ID,
		})
	}
}

func SetTransferAccessHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		botInfo, _ := b.GetMe(ctx)
		userId := update.Message.From.ID

		channelId, err := c.CacheService.GetAwaitingStickerSeparator(ctx, userId)
		if err != nil {
			log.Printf("Erro ao buscar cache sticker: %v", err)
			return
		}

		channel, err := c.ChannelRepo.GetChannelByTwoID(ctx, userId, channelId)
		if err != nil {
			log.Printf("Erro ao buscar canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Canal não encontrado ou você não tem permissão para alterá-lo.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}
		if channel == nil {
			log.Printf("Canal retornado é nil para channelId=%d e userId=%d", channelId, userId)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Erro interno: canal não encontrado.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		newOwnerID, err := strconv.ParseInt(update.Message.Text, 10, 64)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ ID de usuário inválido! Tente novamente com um ID válido.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		if newOwnerID == userId {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ O novo dono precisa ser diferente de voce.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		newOwner, err := b.GetChat(ctx, &bot.GetChatParams{ChatID: newOwnerID})
		if err != nil {
			log.Println("Erro ao obter chat do novo dono:", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ O novo dono precisa iniciar o bot pelo menos uma vez. Peça para ele mandar uma mensagem no bot antes de transferir o canal.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		// Verifica se o novo dono é um bot
		if newOwnerID == botInfo.ID {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ O novo dono não pode ser eu.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		admins, err := b.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
			ChatID: channelId,
		})
		if err != nil {
			log.Println("Erro ao buscar administradores do canal:", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Erro ao consultar administradores do canal.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		isAdmin := false
		for _, admin := range admins {
			if admin.Administrator != nil && admin.Administrator.User.ID == newOwnerID {
				isAdmin = true
				break
			}
			if admin.Owner != nil && admin.Owner.User != nil && admin.Owner.User.ID == newOwnerID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ O novo dono precisa ser administrador do canal.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		// Deletar dados vinculados ao antigo dono
		_ = c.SeparatorRepo.DeleteSeparatorByOwnerChannelId(ctx, userId)

		err = c.ChannelRepo.UpdateOwnerChannel(ctx, channelId, userId, newOwnerID)
		if err != nil {
			log.Printf("Erro ao transferir posse do canal: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Erro ao passar a posse para o novo usuário.",
				ParseMode: "HTML",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		data := map[string]string{
			"channelId":    fmt.Sprintf("%d", channelId),
			"channelName":  channel.Title,
			"newOwnerName": newOwner.LastName,
			"newOwnerId":   fmt.Sprintf("%d", newOwnerID),
			"miniAppUrl":   "https://caption.chelodev.shop",
		}

		textOld, buttonOld := parser.GetMessage("success-old-paccess-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        textOld,
			ReplyMarkup: buttonOld,
			ParseMode:   "HTML",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})

		textNew, buttonNew := parser.GetMessage("success-new-paccess-message", data)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      newOwnerID,
			Text:        textNew,
			ReplyMarkup: buttonNew,
			ParseMode:   "HTML",
		})
	}
}
