package middleware

import (
	"context"
	"log"
	"reflect"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func CheckAddBotMiddleware(c *container.AppContainer) func(next bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Error in permission middleware: %v", r)
					return
				}
			}()

			if update.MyChatMember != nil {
				if !handleMyChatMember(ctx, b, update.MyChatMember, c) {
					return
				}
				next(ctx, b, update)
				return
			}

			// Verificar mensagem encaminhada com verificações de nil adequadas
			if update.Message != nil &&
				update.Message.ForwardOrigin != nil &&
				update.Message.ForwardOrigin.MessageOriginChannel != nil {
				if !handleForwardedMessage(ctx, b, update.Message) {
					return
				}
				next(ctx, b, update)
				return
			}

			log.Printf("Update not related to my_chat_member or valid forwarded channel. Ignoring.")
		}
	}
}

func handleMyChatMember(ctx context.Context, b *bot.Bot, chatMember *models.ChatMemberUpdated, c *container.AppContainer) bool {
	if chatMember == nil {
		log.Printf("ChatMemberUpdated is nil")
		return false
	}

	if isEmpty(chatMember.Chat) {
		log.Printf("Chat is empty")
		return false
	}

	if isEmpty(chatMember.OldChatMember) {
		log.Printf("OldChatMember is nil")
		return false
	}

	if isEmpty(chatMember.NewChatMember) {
		log.Printf("NewChatMember is nil")
		return false
	}

	if isEmpty(chatMember.From) {
		log.Printf("From is empty")
		return false
	}

	chatId := chatMember.Chat.ID
	oldStatus := chatMember.OldChatMember.Type
	newStatus := chatMember.NewChatMember.Type

	log.Printf("Status change: %s -> %s in chat %d", oldStatus, newStatus, chatId)

	if newStatus == "left" || newStatus == "kicked" {
		log.Printf("Bot was removed from channel %d. Skipping.", chatId)

		channel, err := c.ChannelRepo.GetChannelByID(ctx, chatId)
		if err != nil {
			log.Printf("Canal nao Encontrado: %v", err)
		}

		err = c.ChannelRepo.DeleteChannelWithRelations(ctx, channel.OwnerID, chatId)
		if err != nil {
			log.Printf("Erro ao remover canal: %v", err)
		}

		return false
	}

	if oldStatus == "member" || oldStatus == "administrator" || oldStatus == "creator" {
		log.Printf("Bot was already in channel %d. Skipping.", chatId)
		return false
	}

	if newStatus != "administrator" {
		log.Printf("Bot was not promoted to admin in channel %d. Skipping.", chatId)
		return false
	}

	if !hasRequiredPermissions(&chatMember.NewChatMember) {
		sendPermissionErrorMessage(ctx, b, chatMember.From.ID)
		log.Printf("Bot added to channel %d without required permissions.", chatId)
		return false
	}

	log.Printf("Bot successfully added with all necessary permissions to channel %d.", chatId)
	return true
}

func handleForwardedMessage(ctx context.Context, b *bot.Bot, message *models.Message) bool {
	if message == nil {
		log.Printf("Message is nil")
		return false
	}

	if isEmpty(message.From) {
		log.Printf("Message.From is empty")
		return false
	}

	if message.ForwardOrigin == nil {
		log.Printf("ForwardOrigin is nil")
		return false
	}

	if message.ForwardOrigin.MessageOriginChannel == nil {
		log.Printf("MessageOriginChannel is nil")
		return false
	}

	if isEmpty(message.ForwardOrigin.MessageOriginChannel.Chat) {
		log.Printf("MessageOriginChannel.Chat is empty")
		return false
	}

	forwardedChatID := message.ForwardOrigin.MessageOriginChannel.Chat.ID

	botInfo, err := b.GetMe(ctx)
	if err != nil {
		log.Printf("Failed to get bot info: %v", err)
		return false
	}

	botMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: forwardedChatID,
		UserID: botInfo.ID,
	})
	if err != nil {
		log.Printf("Bot is not in forwarded chat %d, or failed to fetch info: %v",
			forwardedChatID, err)
		return false
	}

	if !hasRequiredPermissions(botMember) {
		sendPermissionErrorMessage(ctx, b, message.From.ID)
		log.Printf("Bot is in forwarded chat %d, but lacks required permissions.", forwardedChatID)
		return false
	}

	log.Printf("Bot is present in forwarded channel %d with correct permissions.", forwardedChatID)
	return true
}

func hasRequiredPermissions(chatMember *models.ChatMember) bool {
	if chatMember == nil {
		log.Printf("ChatMember is nil")
		return false
	}

	switch chatMember.Type {
	case "administrator":
		admin := chatMember.Administrator
		if admin == nil {
			log.Printf("Administrator field is nil")
			return false
		}

		canPost := admin.CanPostMessages
		canEdit := admin.CanEditMessages
		canDelete := admin.CanDeleteMessages
		canInvite := admin.CanInviteUsers

		log.Printf("Admin permissions - Post: %v, Edit: %v, Delete: %v, Invite: %v",
			canPost, canEdit, canDelete, canInvite)

		return canPost && canEdit && canDelete && canInvite

	case "creator":
		log.Printf("User is creator, has all permissions")
		return true

	default:
		log.Printf("ChatMember type is not administrator or creator: %s", chatMember.Type)
		return false
	}
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return val.IsZero()
}

func sendPermissionErrorMessage(ctx context.Context, b *bot.Bot, userID int64) {
	text, button := parser.GetMessage("toadd-notfound-permissions-message", nil)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        text,
		ReplyMarkup: button,
		ParseMode:   "HTML",
	})
}
