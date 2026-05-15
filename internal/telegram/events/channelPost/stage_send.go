package channelpost

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func StageSend(c *container.AppContainer) Stage {
	return func(pCtx *ProcessingContext) error {
		// Perform sending with retry/backoff logic
		logger.Bot("📤 Iniciando envio final (Tipo: %s, Album: %v)", pCtx.MessageType, pCtx.IsMediaGroup)
		
		err := processWithRetry(pCtx.Ctx, func() error {
			if pCtx.IsMediaGroup {
				return ProcessMediaGroupDispatch(pCtx)
			}

			switch pCtx.MessageType {
			case MessageTypeText:
				return ProcessTextDispatch(pCtx)
			case MessageTypeAudio, MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation, MessageTypeDocument:
				return ProcessMediaDispatch(pCtx)
			case MessageTypeSticker:
				return ProcessStickerDispatch(pCtx)
			}
			return nil
		})

		if err != nil {
			logger.Error("BOT", "❌ Falha final no envio: %v", err)
			return err
		}

		logger.Bot("✅ Postagem concluída com sucesso no canal %d", pCtx.Channel.ID)
		return nil
	}
}

func processWithRetry(ctx context.Context, fn func() error) error {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		
		if strings.Contains(err.Error(), "Too Many Requests") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter == 0 {
				retryAfter = (attempt + 1) * 2
			}
			logger.Bot("⏳ Rate limit atingido pelo Telegram, aguardando %d segundos (tentativa %d/%d)...", retryAfter, attempt+1, maxRetries)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("failed after %d attempts", maxRetries)
}
