package channelpost

import (
	"context"
	"strings"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func StagePreflight(c *container.AppContainer) Stage {
	return func(pCtx *ProcessingContext) error {
		post := pCtx.Update.ChannelPost
		if post == nil {
			pCtx.StopPipeline = true
			return nil
		}

		// 1. Basic Filters
		botInfo, _ := pCtx.Bot.GetMe(pCtx.Ctx)
		if post.ViaBot != nil && post.ViaBot.ID == botInfo.ID {
			logger.Bot("⏭️ Ignoring inline message.")
			pCtx.StopPipeline = true
			return nil
		}

		maintenance, _ := c.ServerService.GetMaintenance(pCtx.Ctx)
		if maintenance {
			pCtx.StopPipeline = true
			return nil
		}

		// 2. Load Channel
		channel, err := c.ChannelService.GetChannelWithRelations(pCtx.Ctx, post.Chat.ID)
		if err != nil {
			logger.Error("PIPELINE", "❌ Canal %d não encontrado no banco: %v", post.Chat.ID, err)
			pCtx.StopPipeline = true
			return nil
		}
		pCtx.Channel = channel
		logger.Bot("📁 Canal carregado: %s (%d)", channel.Title, channel.ID)

		// 3. Blacklist Check
		if channel.Owner != nil && channel.Owner.IsBlacklisted {
			logger.Bot("🚫 Canal %d ignorado: Proprietário está na Blacklist", channel.ID)
			pCtx.StopPipeline = true
			return nil
		}

		// 4. Metadata Sync (Proactive/Debounced)
		syncMetadata(pCtx, c)

		// 5. Detect Message Type
		pCtx.MessageType = GetMessageType(post)
		if pCtx.MessageType == "" {
			logger.Bot("⏭️ Tipo de mensagem não suportado ignorado")
			pCtx.StopPipeline = true
			return nil
		}
		logger.Bot("📝 Tipo detectado: %s", pCtx.MessageType)

		// 6. Check Permissions
		pm := GetPermissionManager()
		pCtx.Permissions = pm.CheckPermissions(channel, pCtx.MessageType)
		if !pCtx.Permissions.CanEdit && !pCtx.Permissions.CanAddButtons {
			logger.Error("PIPELINE", "❌ Sem permissões de Edição ou Botões para o canal %d (Tipo: %s)", channel.ID, pCtx.MessageType)
			pCtx.StopPipeline = true
			return nil
		}
		logger.Bot("⚖️ Permissões: Edit=%v, Buttons=%v", pCtx.Permissions.CanEdit, pCtx.Permissions.CanAddButtons)

		return nil
	}
}

func syncMetadata(pCtx *ProcessingContext, c *container.AppContainer) {
	post := pCtx.Update.ChannelPost
	channel := pCtx.Channel

	titleChanged := post.Chat.Title != "" && utils.RemoveHTMLTags(post.Chat.Title) != channel.Title
	usernameChanged := post.Chat.Username != "" && ("@" + post.Chat.Username) != channel.InviteURL && !strings.HasPrefix(channel.InviteURL, "https://t.me/+")

	shouldUpdateNow := titleChanged || usernameChanged

	if !shouldUpdateNow && !c.CacheService.ShouldUpdateChannel(pCtx.Ctx, post.Chat.ID) {
		return
	}

	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var chatObj interface{}
		if shouldUpdateNow {
			chatObj = &post.Chat
		}

		updatedChannel, hasChanges := UpdateChannelBasicInfo(updateCtx, pCtx.Bot, post.Chat.ID, channel, chatObj)
		_ = c.CacheService.SetLastChannelUpdate(updateCtx, post.Chat.ID)

		if hasChanges {
			if err := c.ChannelService.UpdateChannelBasicInfoAndFirstButton(updateCtx, updatedChannel); err != nil {
				logger.Error("PIPELINE", "❌ Error saving metadata for %d: %v", post.Chat.ID, err)
			}
			// Cache is already invalidated by UpdateChannelBasicInfoAndFirstButton
		}
	}()
}
