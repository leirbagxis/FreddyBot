package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func CheckAddBotMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Error in permission middleware: %v", r)
				return
			}
		}()

		if update.MyChatMember != nil {
			if !handleMyChatMember(ctx, b, update.MyChatMember) {
				return
			}
			next(ctx, b, update)
			return
		}

		if update.Message != nil && update.Message.ForwardOrigin.MessageOriginChannel != nil {
			if !handleForwardedMessage(ctx, b, update.Message) {
				return
			}
			next(ctx, b, update)
			return
		}

		log.Printf("Update not related to my_chat_member or valid forwarded channel. Ignoring.")
	}
}

func handleMyChatMember(ctx context.Context, b *bot.Bot, chatMember *models.ChatMemberUpdated) bool {
	chatId := chatMember.Chat.ID
	oldStatus := chatMember.OldChatMember.Type
	newStatus := chatMember.NewChatMember.Type

	log.Printf("Status change: %s -> %s in chat %d", oldStatus, newStatus, chatId)

	if newStatus == "left" || newStatus == "kicked" {
		log.Printf("Bot was removed from channel %d. Skipping.", chatId)
		return false
	}

	if oldStatus == "member" || oldStatus == "admnistrator" || oldStatus == "creator" {
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
	forwardedChatID := message.ForwardOrigin.MessageOriginChannel.Chat.ID
	botInfo, _ := b.GetMe(ctx)

	botMember, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: forwardedChatID,
		UserID: botInfo.ID,
	})
	if err != nil {
		log.Printf("Bot is not in forwarded chat %d, or failed to fetch info: %v",
			forwardedChatID, err)
		return false
	}

	fmt.Println(message.ForwardOrigin.MessageOriginUser.SenderUser)

	// Verificar permissões
	if !hasRequiredPermissions(botMember) {
		sendPermissionErrorMessage(ctx, b, message.From.ID)
		log.Printf("Bot is in forwarded chat %d, but lacks required permissions.", forwardedChatID)
		return false
	}

	log.Printf("Bot is present in forwarded channel %d with correct permissions.", forwardedChatID)
	return true
}

func hasRequiredPermissions(chatMember *models.ChatMember) bool {
	// Verificar se é administrator
	if chatMember.Type != "administrator" {
		return false
	}

	// Acessar diretamente o campo Administrator
	admin := chatMember.Administrator
	if admin == nil {
		return false
	}

	// Verificar as 4 permissões necessárias
	canPostMessages := admin.CanPostMessages
	canEditMessages := admin.CanEditMessages
	canDeleteMessages := admin.CanDeleteMessages
	canInviteUsers := admin.CanInviteUsers

	return canPostMessages && canEditMessages && canDeleteMessages && canInviteUsers
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
