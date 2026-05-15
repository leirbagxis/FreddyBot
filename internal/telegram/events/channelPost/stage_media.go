package channelpost

import (
	"context"
	"time"

	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func StageMediaGrouping(c *container.AppContainer, executionPipeline *Pipeline) Stage {
	return func(pCtx *ProcessingContext) error {
		post := pCtx.Update.ChannelPost
		if post.MediaGroupID == "" {
			return nil
		}

		mediaGroupID := post.MediaGroupID

		// For simplicity and alignment with V2 architecture, we use the singleton manager.
		mgm := GetMediaGroupManager()

		group, loaded := mgm.GetMediaGroup(mediaGroupID)
		if !loaded {
			group = &MediaGroup{
				Messages:           make([]MediaMessage, 0, 10),
				Processed:          false,
				MessageEditAllowed: true,
				ChatID:             post.Chat.ID,
			}
			mgm.SetMediaGroup(mediaGroupID, group)
		}

		group.mu.Lock()
		defer group.mu.Unlock()

		if group.Processed || mgm.IsProcessed(mediaGroupID) {
			pCtx.StopPipeline = true
			return nil
		}

		fileID := getFileID(post)

		group.Messages = append(group.Messages, MediaMessage{
			MessageID:       post.ID,
			FileID:          fileID,
			HasCaption:      post.Caption != "",
			Caption:         post.Caption,
			CaptionEntities: post.CaptionEntities,
		})

		if group.Timer != nil {
			group.Timer.Stop()
		}

		// Backoff timeout based on number of messages
		timeout := time.Duration(800+len(group.Messages)*200) * time.Millisecond
		if timeout > 2*time.Second {
			timeout = 2 * time.Second
		}
		
		group.Timer = time.AfterFunc(timeout, func() {
			group.mu.Lock()
			if group.Processed {
				group.mu.Unlock()
				return
			}
			group.Processed = true
			msgs := group.Messages
			group.mu.Unlock()

			// Mark as processed in the manager for global tracking
			mgm.MarkProcessed(mediaGroupID)
			mgm.DeleteMediaGroup(mediaGroupID)

			logger.Bot("📸 Media group ready: %s (%d messages)", mediaGroupID, len(msgs))

			// Create a new context for the consolidated group processing
			groupCtx := &ProcessingContext{
				Ctx:           context.Background(),
				Bot:           pCtx.Bot,
				Update:        pCtx.Update,
				MessageType:   pCtx.MessageType,
				Channel:       pCtx.Channel,
				Permissions:   pCtx.Permissions,
				IsMediaGroup:  true,
				MediaGroupID:  mediaGroupID,
				GroupMessages: msgs,
				Pipeline:      executionPipeline,
			}
			
			// Enqueue the group for execution
			messageQueue.AddV2ToQueue(groupCtx, executionPipeline)
		})

		pCtx.StopPipeline = true
		return nil
	}
}

func getFileID(post *models.Message) string {
	switch {
	case post.Audio != nil:
		return post.Audio.FileID
	case post.Document != nil:
		return post.Document.FileID
	case post.Video != nil:
		return post.Video.FileID
	case post.Animation != nil:
		return post.Animation.FileID
	case post.Photo != nil && len(post.Photo) > 0:
		// Return the largest photo version
		return post.Photo[len(post.Photo)-1].FileID
	}
	return ""
}
