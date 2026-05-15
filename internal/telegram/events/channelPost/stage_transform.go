package channelpost

import (
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func StageTransform(c *container.AppContainer) Stage {
	return func(pCtx *ProcessingContext) error {
		post := pCtx.Update.ChannelPost
		if post == nil {
			return nil
		}

		// 1. Determine original text and entities
		var baseText string
		var entities []models.MessageEntity
		if pCtx.IsMediaGroup {
			for _, m := range pCtx.GroupMessages {
				if m.HasCaption {
					baseText = m.Caption
					// Note: CaptionEntities are already converted to interface{} in stage_media
					// For transformation we might need to convert back or use a different strategy.
					// For now let's assume we use the first caption found.
					if post.ID == m.MessageID {
						entities = post.CaptionEntities
					}
					break
				}
			}
		} else {
			if pCtx.MessageType == MessageTypeText {
				baseText = post.Text
				entities = post.Entities
			} else {
				baseText = post.Caption
				entities = post.CaptionEntities
			}
		}
		pCtx.OriginalCaption = baseText

		// 2. Format base text with existing Telegram formatting
		formattedBase := processTextWithFormatting(baseText, entities)

		// 3. Extract Hashtag and find Custom Caption
		hashtag := extractHashtag(formattedBase)
		var dbCaption string
		var finalButtons []dbmodels.Button = pCtx.Channel.Buttons
		var custom *dbmodels.CustomCaption

		if hashtag != "" {
			logger.Bot("🔍 Hashtag detectada: #%s", hashtag)
			custom = findCustomCaption(pCtx.Channel, hashtag)
			if custom != nil {
				logger.Bot("✨ Aplicando Custom Caption para hashtag #%s", hashtag)
				// Remove hashtag and format custom caption
				cleanBase := removeHashtag(formattedBase, hashtag)
				formattedBase = cleanBase
				dbCaption = detectParseMode(custom.Caption)
				
				if len(custom.Buttons) > 0 {
					logger.Bot("🔘 Usando botões da Custom Caption")
					finalButtons = convertCustomButtons(custom.Buttons)
				}
				
				// Apply custom caption specific permissions (like link preview)
				if pCtx.MessageType == MessageTypeText && !custom.LinkPreview {
					pCtx.DisableLinkPreview = true
				}
			} else {
				logger.Bot("⚠️ Hashtag #%s não encontrada no banco", hashtag)
			}
		}

		// 4. Fallback to Default Caption if no custom caption was found
		if custom == nil && pCtx.Channel.DefaultCaption != nil {
			logger.Bot("📜 Aplicando Legenda Padrão")
			dbCaption = detectParseMode(pCtx.Channel.DefaultCaption.Caption)
		}

		// 5. Final Assembly
		// Match legacy behavior: Text appends, Media replaces
		if pCtx.MessageType == MessageTypeText {
			pCtx.FormattedText = composeMessage(formattedBase, dbCaption, "\n\n", "append")
		} else {
			if dbCaption != "" {
				pCtx.FormattedText = dbCaption
			} else {
				pCtx.FormattedText = formattedBase
			}
		}
		pCtx.FinalButtons = finalButtons
		logger.Bot("✅ Texto final preparado (tamanho: %d)", len(pCtx.FormattedText))

		// 6. Global Link Preview override from permissions
		if !pCtx.Permissions.CanUseLinkPreview {
			pCtx.DisableLinkPreview = true
		}

		return nil
	}
}

func convertCustomButtons(cbs []dbmodels.CustomCaptionButton) []dbmodels.Button {
	btns := make([]dbmodels.Button, len(cbs))
	for i, cb := range cbs {
		btns[i] = dbmodels.Button{
			ButtonID:   cb.ButtonID,
			NameButton: cb.NameButton,
			ButtonURL:  cb.ButtonURL,
			PositionX:  cb.PositionX,
			PositionY:  cb.PositionY,
		}
	}
	return btns
}
