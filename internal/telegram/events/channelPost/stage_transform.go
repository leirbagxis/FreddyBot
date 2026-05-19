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
		formattedBase := ProcessTextWithFormatting(baseText, entities)

		// 2.1 Dynamic Links Extraction (Only from original content)
		extractedDynLinks := false
		if pCtx.Channel.DynamicLinks {
			dynButtons, cleanBase := ExtractDynamicLinks(formattedBase)
			if len(dynButtons) > 0 {
				logger.Bot("🔗 Extraídos %d botões dinâmicos do conteúdo original", len(dynButtons))
				formattedBase = cleanBase
				extractedDynLinks = true
				// Serão adicionados ao finalButtons mais tarde
				pCtx.FinalButtons = append(pCtx.FinalButtons, dynButtons...)
			}
		}

		// 3. Extract Hashtag and find Custom Caption
		hashtag := extractHashtag(formattedBase)
		var dbCaption string
		var finalButtons []dbmodels.Button = pCtx.Channel.Buttons
		var custom *dbmodels.CustomCaption

		// 3.1 Regra: Omitir botões do bot se extraiu links dinâmicos e DLBotButtons for falso
		if extractedDynLinks && !pCtx.Channel.DLBotButtons {
			logger.Bot("🚫 Omitindo botões do bot (regra DynamicLinks)")
			finalButtons = []dbmodels.Button{}
		}

		if hashtag != "" {
			logger.Bot("🔍 Hashtag detectada: #%s", hashtag)
			custom = findCustomCaption(pCtx.Channel, hashtag)
			if custom != nil {
				logger.Bot("✨ Aplicando Custom Caption para hashtag #%s", hashtag)
				// Remove hashtag and format custom caption
				cleanBase := removeHashtag(formattedBase, hashtag)
				formattedBase = cleanBase
				dbCaption = DetectParseMode(custom.Caption)
				
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
			dbCaption = DetectParseMode(pCtx.Channel.DefaultCaption.Caption)
		}

		// 4.1 Regra: Omitir legenda do bot se extraiu links dinâmicos e DLBotCaptions for falso
		if extractedDynLinks && !pCtx.Channel.DLBotCaptions {
			logger.Bot("🚫 Omitindo legenda do bot (regra DynamicLinks)")
			dbCaption = ""
		}

		// 5. Final Assembly
		// Match legacy behavior: Text appends, Media logic varies by type
		if pCtx.MessageType == MessageTypeText {
			pCtx.FormattedText = composeMessage(formattedBase, dbCaption, "\n\n", "append")
		} else if pCtx.MessageType == MessageTypeAudio {
			// For audio, we strictly use the bot's caption (replace)
			if dbCaption != "" {
				pCtx.FormattedText = dbCaption
			} else {
				pCtx.FormattedText = formattedBase
			}
		} else {
			// For other media (photos, videos, etc.), we append the bot's caption to the original
			if dbCaption != "" {
				pCtx.FormattedText = composeMessage(formattedBase, dbCaption, "\n\n", "append")
			} else {
				pCtx.FormattedText = formattedBase
			}
		}

		pCtx.FinalButtons = append(finalButtons, pCtx.FinalButtons...)
		logger.Bot("✅ Texto final preparado (tamanho: %d)", len(pCtx.FormattedText))

		// 5.2 Regra: Omitir reações se extraiu links dinâmicos e DLBotReactions for falso
		if extractedDynLinks && !pCtx.Channel.DLBotReactions {
			logger.Bot("🚫 Desativando reações (regra DynamicLinks)")
			pCtx.Permissions.CanAddReactions = false
		}

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
