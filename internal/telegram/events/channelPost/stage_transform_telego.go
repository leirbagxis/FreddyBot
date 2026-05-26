package channelpost

import (
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/core/services"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

func StageTransformTelego(c *container.AppContainer) StageTelego {
	return func(pCtx *ProcessingContextTelego) error {
		post := pCtx.Update.ChannelPost
		if post == nil {
			return nil
		}

		// 1. Determine original text and entities
		var baseText string
		var entities []telego.MessageEntity
		if pCtx.IsMediaGroup {
			for _, m := range pCtx.GroupMessages {
				if m.HasCaption {
					baseText = m.Caption
					entities = m.CaptionEntities
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

		// 2. Format base text
		formattedBase := ProcessTextWithFormattingTelego(baseText, entities)

		// 2.1 Dynamic Links
		extractedDynLinks := false
		if pCtx.Channel.DynamicLinks {
			dynButtons, cleanBase := ExtractDynamicLinks(formattedBase)
			if len(dynButtons) > 0 {
				logger.Bot("🔗 Extraídos %d botões dinâmicos do conteúdo original", len(dynButtons))
				recordChannelPostEvent(c, pCtx, "dynamic_links_extracted", services.ChannelEventStatusInfo, map[string]any{"count": len(dynButtons)}, nil)
				formattedBase = cleanBase
				extractedDynLinks = true
				pCtx.FinalButtons = append(pCtx.FinalButtons, dynButtons...)
			}
		}

		// 3. Extract Hashtag
		hashtag := extractHashtag(formattedBase)
		var dbCaption string
		var finalButtons []dbmodels.Button = pCtx.Channel.Buttons
		var custom *dbmodels.CustomCaption

		if extractedDynLinks && !pCtx.Channel.DLBotButtons {
			finalButtons = []dbmodels.Button{}
		}

		if hashtag != "" {
			custom = findCustomCaption(pCtx.Channel, hashtag)
			if custom != nil {
				cleanBase := removeHashtag(formattedBase, hashtag)
				formattedBase = cleanBase
				dbCaption = DetectParseMode(custom.Caption)

				if len(custom.Buttons) > 0 {
					finalButtons = convertCustomButtons(custom.Buttons)
				}

				if pCtx.MessageType == MessageTypeText && !custom.LinkPreview {
					pCtx.DisableLinkPreview = true
				}
			}
		}

		// 4. Fallback to Default
		if custom == nil && pCtx.Channel.DefaultCaption != nil {
			dbCaption = DetectParseMode(pCtx.Channel.DefaultCaption.Caption)
		}

		if extractedDynLinks && !pCtx.Channel.DLBotCaptions {
			dbCaption = ""
		}

		// 5. Final Assembly
		if pCtx.MessageType == MessageTypeText {
			pCtx.FormattedText = composeMessage(formattedBase, dbCaption, "\n\n", "append")
		} else if pCtx.MessageType == MessageTypeAudio {
			if dbCaption != "" {
				pCtx.FormattedText = dbCaption
			} else {
				pCtx.FormattedText = formattedBase
			}
		} else {
			if dbCaption != "" {
				pCtx.FormattedText = composeMessage(formattedBase, dbCaption, "\n\n", "append")
			} else {
				pCtx.FormattedText = formattedBase
			}
		}

		pCtx.FinalButtons = append(finalButtons, pCtx.FinalButtons...)

		if dbCaption != "" {
			recordChannelPostEvent(c, pCtx, "caption_applied", services.ChannelEventStatusInfo, map[string]any{"custom_caption": custom != nil, "message_type": pCtx.MessageType}, nil)
		}

		if extractedDynLinks && !pCtx.Channel.DLBotReactions {
			pCtx.Permissions.CanAddReactions = false
		}

		if !pCtx.Permissions.CanUseLinkPreview {
			pCtx.DisableLinkPreview = true
		}

		return nil
	}
}
