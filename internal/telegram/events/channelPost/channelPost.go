package channelpost

import (
	"context"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		post := update.ChannelPost
		if post == nil {
			return
		}

		processor := NewMessageProcessor(b)
		chat := post.Chat

		channel, err := c.ChannelRepo.GetChannelWithRelations(dbCtx, chat.ID)
		if err != nil {
			log.Printf("Canal %d não encontrado: %v", chat.ID, err)
			return
		}

		messageType := processor.GetMessageType(post)
		if messageType == "" {
			return
		}

		// ✅ VERIFICAR PERMISSÕES USANDO O SISTEMA CORRIGIDO
		permissions := processor.CheckPermissions(channel, messageType)
		if !permissions.CanEdit && !permissions.CanAddButtons {
			log.Printf("❌ Sem permissões para processar mensagem no canal %d", channel.ID)
			return
		}

		// Processar com context separado
		go func() {
			telegramCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			var finalButtons []dbmodels.Button
			if permissions.CanAddButtons {
				finalButtons = channel.Buttons
			}

			err := processor.ProcessMessage(telegramCtx, messageType, channel, post, finalButtons, permissions.CanEdit)
			if err != nil {
				log.Printf("Erro ao processar mensagem: %v", err)
			}

			if channel.Separator != nil && (permissions.CanEdit || permissions.CanAddButtons) {
				err := processor.ProcessSeparator(telegramCtx, channel, post)
				if err != nil {
					log.Printf("Erro ao processar separator: %v", err)
				}
			}
		}()
	}
}

// ✅ CORRIGIDO: ProcessMessage usando as funções corretas do processors.go
func (mp *MessageProcessor) ProcessMessage(ctx context.Context, messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	switch messageType {
	case MessageTypeText:
		return mp.ProcessTextMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeAudio:
		return mp.ProcessAudioMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeSticker:
		if len(buttons) > 0 {
			return mp.ProcessStickerMessage(ctx, channel, post, buttons)
		}
		return nil
	case MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation:
		return mp.ProcessMediaMessage(ctx, channel, post, buttons, messageEditAllowed)
	default:
		return nil
	}
}

func (mp *MessageProcessor) ProcessSeparator(ctx context.Context, channel *dbmodels.Channel, post *models.Message) error {
	if channel.Separator == nil || channel.Separator.SeparatorID == "" {
		return nil
	}

	_, err := mp.bot.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  post.Chat.ID,
		Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
	})
	return err
}
