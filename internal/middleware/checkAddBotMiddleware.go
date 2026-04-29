package middleware

import (
	"context"
	"reflect"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func CheckAddBotMiddleware(c *container.AppContainer) func(next bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("BOT", "Error in permission middleware: %v", r)
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

			logger.Bot("Update not related to my_chat_member or valid forwarded channel. Ignoring.")
		}
	}
}

func handleMyChatMember(ctx context.Context, b *bot.Bot, chatMember *models.ChatMemberUpdated, c *container.AppContainer) bool {
	if chatMember == nil {
		logger.Error("BOT", "ChatMemberUpdated is nil")
		return false
	}

	if isEmpty(chatMember.Chat) {
		logger.Error("BOT", "Chat is empty")
		return false
	}

	if isEmpty(chatMember.OldChatMember) {
		logger.Error("BOT", "OldChatMember is nil")
		return false
	}

	if isEmpty(chatMember.NewChatMember) {
		logger.Error("BOT", "NewChatMember is nil")
		return false
	}

	if isEmpty(chatMember.From) {
		logger.Error("BOT", "From is empty")
		return false
	}

	chatId := chatMember.Chat.ID
	oldStatus := chatMember.OldChatMember.Type
	newStatus := chatMember.NewChatMember.Type

	logger.Bot("Status change: %s -> %s in chat %d", oldStatus, newStatus, chatId)

	if newStatus == "left" || newStatus == "kicked" {
		logger.Bot("Bot was removed from channel %d. Skipping.", chatId)

		channel, err := c.ChannelRepo.GetChannelByID(ctx, chatId)
		if err != nil {
			logger.Error("BOT", "Canal nao Encontrado: %v", err)
		}

		err = c.ChannelRepo.DeleteChannelWithRelations(ctx, channel.OwnerID, chatId)
		if err != nil {
			logger.Error("BOT", "Erro ao remover canal: %v", err)
		}

		return false
	}

	if oldStatus == "member" || oldStatus == "administrator" || oldStatus == "creator" {
		logger.Bot("Bot was already in channel %d. Skipping.", chatId)
		return false
	}

	if newStatus != "administrator" {
		logger.Bot("Bot was not promoted to admin in channel %d. Skipping.", chatId)
		return false
	}

	if !hasRequiredPermissions(&chatMember.NewChatMember) {
		sendPermissionErrorMessage(ctx, b, chatMember.From.ID)
		logger.Error("BOT", "Bot added to channel %d without required permissions.", chatId)
		return false
	}

	logger.Bot("Bot successfully added with all necessary permissions to channel %d.", chatId)
	return true
}

func handleForwardedMessage(ctx context.Context, b *bot.Bot, message *models.Message) bool {
	if message == nil {
		logger.Error("BOT", "Message is nil")
		return false
	}

	if isEmpty(message.From) {
		logger.Error("BOT", "Message.From is empty")
		return false
	}

	if message.ForwardOrigin == nil {
		logger.Error("BOT", "ForwardOrigin is nil")
		return false
	}

	if message.ForwardOrigin.MessageOriginChannel == nil {
		logger.Error("BOT", "MessageOriginChannel is nil")
		return false
	}

	if isEmpty(message.ForwardOrigin.MessageOriginChannel.Chat) {
		logger.Error("BOT", "MessageOriginChannel.Chat is empty")
		return false
	}

	forwardedChatID := message.ForwardOrigin.MessageOriginChannel.Chat.ID

	botInfo, err := b.GetMe(ctx)
	if err != nil {
		logger.Error("BOT", "Failed to get bot info: %v", err)
		return false
	}

	botMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: forwardedChatID,
		UserID: botInfo.ID,
	})
	if err != nil {
		logger.Error("BOT", "Bot is not in forwarded chat %d, or failed to fetch info: %v",
			forwardedChatID, err)
		return false
	}

	if !hasRequiredPermissions(botMember) {
		sendPermissionErrorMessage(ctx, b, message.From.ID)
		logger.Error("BOT", "Bot is in forwarded chat %d, but lacks required permissions.", forwardedChatID)
		return false
	}

	logger.Bot("Bot is present in forwarded channel %d with correct permissions.", forwardedChatID)
	return true
}

func hasRequiredPermissions(chatMember *models.ChatMember) bool {
	if chatMember == nil {
		logger.Error("BOT", "ChatMember is nil")
		return false
	}

	switch chatMember.Type {
	case "administrator":
		admin := chatMember.Administrator
		if admin == nil {
			logger.Error("BOT", "Administrator field is nil")
			return false
		}

		canPost := admin.CanPostMessages
		canEdit := admin.CanEditMessages
		canDelete := admin.CanDeleteMessages
		canInvite := admin.CanInviteUsers

		logger.Bot("Admin permissions - Post: %v, Edit: %v, Delete: %v, Invite: %v",
			canPost, canEdit, canDelete, canInvite)

		return canPost && canEdit && canDelete && canInvite

	case "creator":
		logger.Bot("User is creator, has all permissions")
		return true

	default:
		logger.Bot("ChatMember type is not administrator or creator: %s", chatMember.Type)
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
