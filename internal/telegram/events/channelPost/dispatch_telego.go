package channelpost

import (
	"context"
	"fmt"
	"strings"
	"time"

	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

func ProcessTextDispatchTelego(pCtx *ProcessingContextTelego) error {
	post := pCtx.Update.ChannelPost

	if !pCtx.Permissions.CanEdit {
		if pCtx.FinalKeyboard != nil {
			_, err := pCtx.Bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
				ChatID:      telego.ChatID{ID: post.Chat.ID},
				MessageID:   post.MessageID,
				ReplyMarkup: pCtx.FinalKeyboard,
			})
			return err
		}
		return nil
	}

	params := &telego.EditMessageTextParams{
		ChatID:    telego.ChatID{ID: post.Chat.ID},
		MessageID: post.MessageID,
		Text:      pCtx.FormattedText,
		ParseMode: telego.ModeHTML,
	}
	if pCtx.DisableLinkPreview {
		params.LinkPreviewOptions = &telego.LinkPreviewOptions{IsDisabled: true}
	}
	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageText(context.Background(), params)
	if err == nil {
		HandleSeparatorAfterDispatchTelego(pCtx)
	}
	return err
}

func ProcessMediaDispatchTelego(pCtx *ProcessingContextTelego) error {
	post := pCtx.Update.ChannelPost

	if !pCtx.Permissions.CanEdit {
		if pCtx.FinalKeyboard != nil {
			_, err := pCtx.Bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
				ChatID:      telego.ChatID{ID: post.Chat.ID},
				MessageID:   post.MessageID,
				ReplyMarkup: pCtx.FinalKeyboard,
			})
			return err
		}
		return nil
	}

	params := &telego.EditMessageCaptionParams{
		ChatID:    telego.ChatID{ID: post.Chat.ID},
		MessageID: post.MessageID,
		Caption:   pCtx.FormattedText,
		ParseMode: telego.ModeHTML,
	}

	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageCaption(context.Background(), params)
	if err == nil {
		HandleSeparatorAfterDispatchTelego(pCtx)
	}
	return err
}

func ProcessStickerDispatchTelego(pCtx *ProcessingContextTelego) error {
	post := pCtx.Update.ChannelPost

	if pCtx.FinalKeyboard == nil {
		return nil
	}

	_, err := pCtx.Bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
		ChatID:      telego.ChatID{ID: post.Chat.ID},
		MessageID:   post.MessageID,
		ReplyMarkup: pCtx.FinalKeyboard,
	})
	if err == nil {
		HandleSeparatorAfterDispatchTelego(pCtx)
	}
	return err
}

func ProcessMediaGroupDispatchTelego(pCtx *ProcessingContextTelego) error {
	if len(pCtx.GroupMessages) == 0 {
		return nil
	}

	if pCtx.MessageType == MessageTypeAudio || pCtx.MessageType == MessageTypeDocument {
		return dispatchReSendMediaGroupTelego(pCtx)
	}

	targetMessage := pCtx.GroupMessages[0]
	for _, message := range pCtx.GroupMessages {
		if message.HasCaption {
			targetMessage = message
			break
		}
	}

	params := &telego.EditMessageCaptionParams{
		ChatID:    telego.ChatID{ID: pCtx.Channel.ID},
		MessageID: targetMessage.MessageID,
		Caption:   pCtx.FormattedText,
		ParseMode: telego.ModeHTML,
	}
	if pCtx.FinalKeyboard != nil {
		params.ReplyMarkup = pCtx.FinalKeyboard
	}

	_, err := pCtx.Bot.EditMessageCaption(context.Background(), params)
	if err == nil {
		logger.Bot("✅ Media Group %s (Photos/Videos) processed", pCtx.MediaGroupID)
		HandleSeparatorAfterDispatchTelego(pCtx)
	}

	return err
}

func dispatchReSendMediaGroupTelego(pCtx *ProcessingContextTelego) error {
	for i, m := range pCtx.GroupMessages {
		if i > 0 {
			time.Sleep(time.Duration(200+i*150) * time.Millisecond)
		}

		var err error
		if pCtx.MessageType == MessageTypeAudio {
			params := &telego.SendAudioParams{
				ChatID:    telego.ChatID{ID: pCtx.Channel.ID},
				Audio:     telego.InputFile{FileID: m.FileID},
				Caption:   pCtx.FormattedText,
				ParseMode: telego.ModeHTML,
			}
			if pCtx.FinalKeyboard != nil {
				params.ReplyMarkup = pCtx.FinalKeyboard
			}
			_, err = pCtx.Bot.SendAudio(context.Background(), params)
		} else {
			params := &telego.SendDocumentParams{
				ChatID:    telego.ChatID{ID: pCtx.Channel.ID},
				Document:  telego.InputFile{FileID: m.FileID},
				Caption:   pCtx.FormattedText,
				ParseMode: telego.ModeHTML,
			}
			if pCtx.FinalKeyboard != nil {
				params.ReplyMarkup = pCtx.FinalKeyboard
			}
			_, err = pCtx.Bot.SendDocument(context.Background(), params)
		}

		if err != nil {
			logger.Error("BOT", "❌ Failed to re-send media %d in group %s: %v", m.MessageID, pCtx.MediaGroupID, err)
			continue
		}

		time.Sleep(200 * time.Millisecond)
		_ = pCtx.Bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    telego.ChatID{ID: pCtx.Channel.ID},
			MessageID: m.MessageID,
		})
	}

	logger.Bot("✅ Media Group %s (Re-sent) processed", pCtx.MediaGroupID)
	HandleSeparatorAfterDispatchTelego(pCtx)
	return nil
}

func HandleSeparatorAfterDispatchTelego(pCtx *ProcessingContextTelego) {
	if pCtx.Channel.Separator == nil || pCtx.Channel.Separator.SeparatorID == "" {
		return
	}

	time.AfterFunc(1*time.Second, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = ProcessSeparatorTelego(ctx, pCtx.Bot, pCtx.Channel, nil)
	})
}

func ProcessSeparatorTelego(ctx context.Context, b *telego.Bot, channel *dbmodels.Channel, post *telego.Message) error {
	if channel == nil || channel.Separator == nil || channel.Separator.SeparatorID == "" {
		return nil
	}

	if post != nil && post.MediaGroupID != "" && post.Audio != nil {
		return nil
	}

	var chatID int64
	if post != nil {
		chatID = post.Chat.ID
	} else {
		chatID = channel.ID
	}

	sendCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	maxRetries := 2
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := b.SendSticker(sendCtx, &telego.SendStickerParams{
			ChatID:  telego.ChatID{ID: chatID},
			Sticker: telego.InputFile{FileID: channel.Separator.SeparatorID},
		})
		if err == nil {
			return nil
		}

		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "too many requests") || strings.Contains(lower, "429") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter <= 0 {
				retryAfter = (attempt + 1) * 2
			}
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("failed after %d attempts", maxRetries)
}
