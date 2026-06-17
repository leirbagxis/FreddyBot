package channelpost

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/core/services"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

func StageSendTelego(c *container.AppContainer) StageTelego {
	return func(pCtx *ProcessingContextTelego) error {
		logger.Bot("📤 Iniciando envio final (Tipo: %s, Album: %v)", pCtx.MessageType, pCtx.IsMediaGroup)

		ownerID := pCtx.Channel.OwnerID

		// Find the target message ID (handles both single posts and media groups)
		targetMsgID := int32(0)
		if pCtx.IsMediaGroup {
			for _, m := range pCtx.GroupMessages {
				if m.HasCaption || targetMsgID == 0 {
					targetMsgID = int32(m.MessageID)
				}
				if m.HasCaption {
					break
				}
			}
		} else if pCtx.Update.ChannelPost != nil {
			targetMsgID = int32(pCtx.Update.ChannelPost.MessageID)
		}

		// MTProto tenta editar TUDO (texto + entidades + botões).
		// Se conseguir texto mas não botões (sempre em canais), o BOT aplica só os botões.
		mtprotoOK := false
		if ownerID != 0 && pCtx.FormattedText != "" && targetMsgID > 0 &&
			strings.Contains(pCtx.FormattedText, "<tg-emoji") {
			logger.Bot("🔄 MTProto: editando msg %d (texto + entidades + tentativa botões)...", targetMsgID)
			mtprotoCtx, mtprotoCancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := c.TelegramPosterService.EditMessage(
				mtprotoCtx,
				ownerID,
				pCtx.Channel.ID,
				targetMsgID,
				pCtx.FormattedText,
				nil, // keyboard via MTProto não funciona em canais; bot faz fallback
			)
			mtprotoCancel()
			if err == nil {
				logger.Info("MTPOST", "Texto editado via MTProto como user %d no canal %d (msgID=%d)", ownerID, pCtx.Channel.ID, targetMsgID)
				mtprotoOK = true
			} else {
				logger.Info("MTPOST", "MTProto: fallback total para bot: %v", err)
			}
		}

		if mtprotoOK {
			// MTProto tentou botões (silenciosamente ignorado em canais).
			// Fallback: BOT aplica apenas os botões na mensagem já editada pelo MTProto.
			if pCtx.FinalKeyboard != nil && pCtx.Permissions.CanEdit {
				logger.Bot("⌨️ Bot: aplicando apenas botões na msg %d (fallback MTProto)...", targetMsgID)
				_, kbErr := pCtx.Bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
					ChatID:      telego.ChatID{ID: pCtx.Channel.ID},
					MessageID:   int(targetMsgID),
					ReplyMarkup: pCtx.FinalKeyboard,
				})
				if kbErr != nil {
					logger.Warn("MTPOST", "Bot keyboard update failed: %v", kbErr)
				}
			} else if pCtx.FinalKeyboard == nil && pCtx.Permissions.CanEdit {
				_, _ = pCtx.Bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
					ChatID:    telego.ChatID{ID: pCtx.Channel.ID},
					MessageID: int(targetMsgID),
				})
			}
			HandleSeparatorAfterDispatchTelego(pCtx)
			recordChannelPostEvent(c, pCtx, "post_processed", services.ChannelEventStatusSuccess, map[string]any{"album": pCtx.IsMediaGroup, "method": "mtproto", "buttons": len(pCtx.FinalButtons), "has_caption": pCtx.FormattedText != ""}, nil)
			logger.Bot("✅ Postagem MTProto concluída no canal %d", pCtx.Channel.ID)
			return nil
		}

		err := processWithRetryTelego(pCtx.Ctx, func() error {
			if pCtx.IsMediaGroup {
				return ProcessMediaGroupDispatchTelego(pCtx)
			}
			switch pCtx.MessageType {
			case MessageTypeText:
				return ProcessTextDispatchTelego(pCtx)
			case MessageTypeAudio, MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation, MessageTypeDocument:
				return ProcessMediaDispatchTelego(pCtx)
			case MessageTypeSticker:
				return ProcessStickerDispatchTelego(pCtx)
			}
			return nil
		})

		if err != nil {
			logger.Error("BOT", "❌ Falha final no envio via bot: %v", err)
			recordChannelPostEvent(c, pCtx, "post_failed", services.ChannelEventStatusError, map[string]any{"album": pCtx.IsMediaGroup, "method": "bot"}, err)
			return err
		}

		recordChannelPostEvent(c, pCtx, "post_processed", services.ChannelEventStatusSuccess, map[string]any{"album": pCtx.IsMediaGroup, "method": "bot", "buttons": len(pCtx.FinalButtons), "has_caption": pCtx.FormattedText != ""}, nil)
		logger.Bot("✅ Postagem via bot concluída no canal %d", pCtx.Channel.ID)
		return nil
	}
}

func processWithRetryTelego(ctx context.Context, fn func() error) error {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "Too Many Requests") || strings.Contains(err.Error(), "429") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter == 0 {
				retryAfter = (attempt + 1) * 2
			}
			logger.Bot("⏳ Rate limit Telego, aguardando %d segundos...", retryAfter)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("failed after %d attempts", maxRetries)
}
