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

		// Verificar permissões
		messageEditAllowed := processor.permissionManager.IsMessageEditAllowed(channel, messageType)
		buttonsAllowed := processor.permissionManager.IsButtonsAllowed(channel, messageType)

		if !messageEditAllowed && !buttonsAllowed {
			return
		}

		// Processar com context separado
		go func() {
			telegramCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			var finalButtons []dbmodels.Button
			if buttonsAllowed {
				finalButtons = channel.Buttons
			}

			err := processor.ProcessMessage(telegramCtx, messageType, channel, post, finalButtons, messageEditAllowed)
			if err != nil {
				log.Printf("Erro ao processar mensagem: %v", err)
			}

			if channel.Separator != nil && (messageEditAllowed || buttonsAllowed) {
				err := processor.ProcessSeparator(telegramCtx, channel, post)
				if err != nil {
					log.Printf("Erro ao processar separator: %v", err)
				}
			}
		}()
	}
}

// GetMessageType determines the type of the incoming message.
func (mp *MessageProcessor) GetMessageType(post *models.Message) MessageType {
	if post.Text != "" {
		return MessageTypeText
	}
	if post.Audio != nil {
		return MessageTypeAudio
	}
	if post.Sticker != nil {
		return MessageTypeSticker
	}
	if post.Photo != nil && len(post.Photo) > 0 {
		return MessageTypePhoto
	}
	if post.Video != nil {
		return MessageTypeVideo
	}
	if post.Animation != nil {
		return MessageTypeAnimation
	}
	return ""
}

// ✅ CORRIGIDO: ProcessMessage usando as funções corretas do processors.go
func (mp *MessageProcessor) ProcessMessage(ctx context.Context, messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	switch messageType {
	case MessageTypeText:
		return mp.ProcessTextMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeAudio:
		// ✅ CORRIGIDO: Usar ProcessAudioMessage para áudios
		return mp.ProcessAudioMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeSticker:
		if len(buttons) > 0 {
			return mp.ProcessStickerMessage(ctx, post, buttons)
		}
		return nil
	case MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation:
		// ✅ CORRIGIDO: Usar ProcessMediaMessage apenas para fotos/vídeos/gifs
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

// // ✅ FUNÇÃO PARA REGISTRAR O HANDLER
// func LoadChannelPostHandlers(b *bot.Bot, app *container.AppContainer) {
// 	b.RegisterHandler(
// 		bot.HandlerTypeChannelPost,
// 		"",
// 		bot.MatchTypeExact,
// 		Handler(app),
// 		middleware.CheckAddBotMiddleware,
// 	)
// }
