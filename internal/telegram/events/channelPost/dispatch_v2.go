package channelpost

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func ProcessTextDispatch(pCtx *ProcessingContext) error {
	post := pCtx.Update.ChannelPost
	
	if !pCtx.Permissions.CanEdit {
		if pCtx.FinalKeyboard != nil {
			_, err := pCtx.Bot.EditMessageReplyMarkup(pCtx.Ctx, &bot.EditMessageReplyMarkupParams{
				ChatID:      post.Chat.ID,
				MessageID:   post.ID,
				ReplyMarkup: pCtx.FinalKeyboard,
			})
			return err
		}
		return nil
	}

	params := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: post.ID,
		Text:      pCtx.FormattedText,
		ParseMode: "HTML",
	}
	if pCtx.DisableLinkPreview {
		val := true
		params.LinkPreviewOptions = &models.LinkPreviewOptions{IsDisabled: &val}
	}
	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageText(pCtx.Ctx, params)
	if err == nil {
		HandleSeparatorAfterDispatch(pCtx)
	}
	return err
}

func ProcessMediaDispatch(pCtx *ProcessingContext) error {
	post := pCtx.Update.ChannelPost

	if !pCtx.Permissions.CanEdit {
		if pCtx.FinalKeyboard != nil {
			_, err := pCtx.Bot.EditMessageReplyMarkup(pCtx.Ctx, &bot.EditMessageReplyMarkupParams{
				ChatID:      post.Chat.ID,
				MessageID:   post.ID,
				ReplyMarkup: pCtx.FinalKeyboard,
			})
			return err
		}
		return nil
	}

	params := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: post.ID,
		Caption:   pCtx.FormattedText,
		ParseMode: "HTML",
	}
	// Note: EditMessageCaptionParams does NOT support LinkPreviewOptions in the current Bot API version.
	
	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageCaption(pCtx.Ctx, params)
	if err == nil {
		HandleSeparatorAfterDispatch(pCtx)
	}
	return err
}

func ProcessStickerDispatch(pCtx *ProcessingContext) error {
	post := pCtx.Update.ChannelPost

	if pCtx.FinalKeyboard == nil {
		return nil
	}

	_, err := pCtx.Bot.EditMessageReplyMarkup(pCtx.Ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      post.Chat.ID,
		MessageID:   post.ID,
		ReplyMarkup: pCtx.FinalKeyboard,
	})
	if err == nil {
		HandleSeparatorAfterDispatch(pCtx)
	}
	return err
}

func ProcessMediaGroupDispatch(pCtx *ProcessingContext) error {
	if len(pCtx.GroupMessages) == 0 {
		return nil
	}

	// Audio and Documents are handled by re-sending every item and deleting the original
	// to ensure buttons and new captions are applied correctly.
	if pCtx.MessageType == MessageTypeAudio || pCtx.MessageType == MessageTypeDocument {
		return dispatchReSendMediaGroup(pCtx)
	}

	// Photos and Videos are handled by editing the caption of the target message
	targetMessage := pCtx.GroupMessages[0]
	params := &bot.EditMessageCaptionParams{
		ChatID:    pCtx.Channel.ID,
		MessageID: targetMessage.MessageID,
		Caption:   pCtx.FormattedText,
		ParseMode: "HTML",
	}
	if pCtx.DisableLinkPreview {
		// Note: EditMessageCaption doesn't support LinkPreviewOptions in current Bot API
	}
	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageCaption(pCtx.Ctx, params)
	if err == nil {
		logger.Bot("✅ Media Group %s (Photos/Videos) processed", pCtx.MediaGroupID)
		HandleSeparatorAfterDispatch(pCtx)
	}
	
	return err
}

func dispatchReSendMediaGroup(pCtx *ProcessingContext) error {
	for i, m := range pCtx.GroupMessages {
		// Apply delay between re-sends to avoid floods
		if i > 0 {
			time.Sleep(time.Duration(200+i*150) * time.Millisecond)
		}

		var err error
		if pCtx.MessageType == MessageTypeAudio {
			params := &bot.SendAudioParams{
				ChatID:    pCtx.Channel.ID,
				Audio:     &models.InputFileString{Data: m.FileID},
				Caption:   pCtx.FormattedText,
				ParseMode: "HTML",
			}
			if pCtx.FinalKeyboard != nil {
				params.ReplyMarkup = pCtx.FinalKeyboard
			}
			_, err = pCtx.Bot.SendAudio(pCtx.Ctx, params)
		} else {
			params := &bot.SendDocumentParams{
				ChatID:    pCtx.Channel.ID,
				Document:  &models.InputFileString{Data: m.FileID},
				Caption:   pCtx.FormattedText,
				ParseMode: "HTML",
			}
			if pCtx.FinalKeyboard != nil {
				params.ReplyMarkup = pCtx.FinalKeyboard
			}
			_, err = pCtx.Bot.SendDocument(pCtx.Ctx, params)
		}

		if err != nil {
			logger.Error("BOT", "❌ Failed to re-send media %d in group %s: %v", m.MessageID, pCtx.MediaGroupID, err)
			continue
		}

		// Delete original message
		time.Sleep(200 * time.Millisecond)
		_, _ = pCtx.Bot.DeleteMessage(pCtx.Ctx, &bot.DeleteMessageParams{
			ChatID:    pCtx.Channel.ID,
			MessageID: m.MessageID,
		})
	}

	logger.Bot("✅ Media Group %s (Re-sent) processed", pCtx.MediaGroupID)
	HandleSeparatorAfterDispatch(pCtx)
	return nil
}

func HandleSeparatorAfterDispatch(pCtx *ProcessingContext) {
	if pCtx.Channel.Separator == nil || pCtx.Channel.Separator.SeparatorID == "" {
		return
	}

	// Small delay to ensure order in Telegram
	time.AfterFunc(1*time.Second, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = ProcessSeparator(ctx, pCtx.Bot, pCtx.Channel, nil)
	})
}

// Envia o separador independentemente de CanEdit/CanAddButtons; suprime apenas no início de álbum de áudio.
// Use esta função no Handler após enfileirar a mensagem; o finalizador de grupo enviará no fim do álbum.
func ProcessSeparator(ctx context.Context, b *bot.Bot, channel *dbmodels.Channel, post *models.Message) error {
	if channel == nil || channel.Separator == nil || channel.Separator.SeparatorID == "" {
		logger.Bot("⚠️ Separator não configurado para o canal")
		return nil
	}

	// Suprime no início do álbum de áudio: deixa para o finalizador do grupo
	if post != nil && post.MediaGroupID != "" && post.Audio != nil {
		logger.Bot("ℹ️ Separator suprimido no início do álbum de áudio (groupID=%s)", post.MediaGroupID)
		return nil
	}

	// Determinar chat alvo
	var chatID int64
	if post != nil {
		chatID = post.Chat.ID
	} else {
		chatID = channel.ID
	}

	// Contexto próprio com timeout
	sendCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Bot("🔄 Enviando separator para chat %d", chatID)

	maxRetries := 2
	baseDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := b.SendSticker(sendCtx, &bot.SendStickerParams{
			ChatID:  chatID,
			Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
		})
		if err == nil {
			logger.Bot("✅ Separator enviado com sucesso para chat %d", chatID)
			return nil
		}

		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "too many requests") || strings.Contains(lower, "429") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter <= 0 {
				retryAfter = int(baseDelay.Seconds()) * (attempt + 1)
			}
			logger.Bot("⏳ Rate limit no separator, aguardando %d segundos (tentativa %d/%d)", retryAfter, attempt+1, maxRetries)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		logger.Error("BOT", "❌ Erro ao enviar separator: %v", err)
		return err
	}

	return fmt.Errorf("falha após %d tentativas no envio do separator", maxRetries)
}
