package postbuilder

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/container"
	channelpost "github.com/leirbagxis/FreddyBot/internal/telegram/events/channelPost"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func isEmoji(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		// Check Blacklist
		user, err := c.UserRepo.GetUserById(ctx, update.Message.From.ID)
		if err == nil && user != nil && user.IsBlacklisted {
			return
		}

		// Detect media
		var mediaID string
		var mediaType string

		if update.Message.Photo != nil {
			mediaID = update.Message.Photo[len(update.Message.Photo)-1].FileID
			mediaType = "photo"
		} else if update.Message.Video != nil {
			mediaID = update.Message.Video.FileID
			mediaType = "video"
		} else if update.Message.Animation != nil {
			mediaID = update.Message.Animation.FileID
			mediaType = "animation"
		} else if update.Message.Audio != nil {
			mediaID = update.Message.Audio.FileID
			mediaType = "audio"
		} else if update.Message.Document != nil {
			mediaID = update.Message.Document.FileID
			mediaType = "document"
		}

		if mediaID == "" {
			// Check if we are in a state of awaiting text input
			state, _ := c.CacheService.GetPostBuilderState(ctx, update.Message.From.ID)
			if state != nil && state.Step != "" {
				handleTextInput(ctx, b, update, c, state)
				return
			}
			return
		}

		// Media detected, offer Post Builder
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "🛠️ Post Builder", CallbackData: "pb-start"},
				},
			},
		}

		// Store initial state
		state := cache.PostBuilderState{
			MediaType:   mediaType,
			MediaFileID: mediaID,
			Step:        "",
		}
		logger.Bot("PostBuilder: Saving initial state for user %d", update.Message.From.ID)
		c.CacheService.SetPostBuilderState(ctx, update.Message.From.ID, state)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "✨ Media detectada! Deseja usar o <b>Post Builder</b> para criar uma postagem personalizada?",
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}

func handleTextInput(ctx context.Context, b *bot.Bot, update *models.Update, c *container.AppContainer, state *cache.PostBuilderState) {
	text := update.Message.Text
	// Get formatted text if available
	formattedText := text
	if len(update.Message.Entities) > 0 {
		formattedText = channelpost.ProcessEntitiesOnly(text, update.Message.Entities)
	}

	switch state.Step {
	case "awaiting_title":
		state.Title = formattedText
		state.Step = ""
	case "awaiting_body":
		state.Body = formattedText
		state.Step = ""
	case "awaiting_footer":
		state.Footer = formattedText
		state.Step = ""
	case "awaiting_reactions":
		// Validação de emojis
		parts := strings.Split(text, ",")
		valid := true
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if !isEmoji(p) {
				valid = false
				break
			}
		}

		if !valid {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Apenas emojis são permitidos como reações. Tente novamente:",
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		state.Reactions = text
		state.Step = ""
	case "awaiting_button":
		lines := strings.Split(text, "\n")
		if len(lines) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Formato inválido. Envie o <b>Nome</b> em uma linha e o <b>Link</b> na linha de baixo.",
				ParseMode: models.ParseModeHTML,
			})
			return
		}
		name := strings.TrimSpace(lines[0])
		url := strings.TrimSpace(lines[1])

		if !strings.HasPrefix(url, "http") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ URL inválida. Deve começar com http:// ou https://. Tente novamente:",
			})
			return
		}

		state.Buttons = append(state.Buttons, cache.PostBuilderButton{Text: name, URL: url})
		state.Step = ""
	default:
		return
	}

	c.CacheService.SetPostBuilderState(ctx, update.Message.From.ID, *state)
	showMenu(ctx, b, update.Message.Chat.ID, update.Message.From.ID, c, state)
}

func showMenu(ctx context.Context, b *bot.Bot, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState) {
	text := "🛠️ <b>Post Builder - Menu</b>\n\n"
	text += fmt.Sprintf("📝 <b>Título:</b> %s\n", state.Title)
	text += fmt.Sprintf("📄 <b>Corpo:</b> %s\n", state.Body)
	text += fmt.Sprintf("👣 <b>Rodapé:</b> %s\n", state.Footer)
	text += fmt.Sprintf("🎭 <b>Reações:</b> %s\n", state.Reactions)
	text += fmt.Sprintf("🔘 <b>Botões:</b> %d\n\n", len(state.Buttons))
	text += "Escolha o que deseja editar:"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "📝 Título", CallbackData: "pb-edit-title"},
				{Text: "📄 Corpo", CallbackData: "pb-edit-body"},
			},
			{
				{Text: "👣 Rodapé", CallbackData: "pb-edit-footer"},
				{Text: "🎭 Reações", CallbackData: "pb-edit-reactions"},
			},
			{
				{Text: "🔘 Botão", CallbackData: "pb-add-button"},
				{Text: "👁️ Preview", CallbackData: "pb-preview"},
			},
			{
				{Text: "✅ Salvar", CallbackData: "pb-save"},
				{Text: "❌ Cancelar", CallbackData: "pb-cancel"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: kb,
	})
}

func CallbackHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.Message.Chat.ID
		data := update.CallbackQuery.Data

		// Check Blacklist
		user, err := c.UserRepo.GetUserById(ctx, userID)
		if err == nil && user != nil && user.IsBlacklisted {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você está na blacklist.",
				ShowAlert:       true,
			})
			return
		}

		state, _ := c.CacheService.GetPostBuilderState(ctx, userID)
		if state == nil && data != "pb-cancel" {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Sessão expirada ou não encontrada.",
			})
			return
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		switch data {
		case "pb-start":
			showMenu(ctx, b, chatID, userID, c, state)
		case "pb-edit-title":
			state.Step = "awaiting_title"
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "📝 Envie o <b>Título</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
		case "pb-edit-body":
			state.Step = "awaiting_body"
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "📄 Envie o <b>Corpo</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
		case "pb-edit-footer":
			state.Step = "awaiting_footer"
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "👣 Envie o <b>Rodapé</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
		case "pb-edit-reactions":
			state.Step = "awaiting_reactions"
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "🎭 Envie as <b>Reações</b> separadas por vírgula (ex: 👍,👎,❤️):",
				ParseMode: models.ParseModeHTML,
			})
		case "pb-add-button":
			state.Step = "awaiting_button"
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "🔘 Envie os dados do botão no formato:\n\n<code>Nome do Botão\nhttps://link.com</code>",
				ParseMode: models.ParseModeHTML,
			})
		case "pb-preview":
			sendFinalPost(ctx, b, chatID, userID, c, state, false)
			showMenu(ctx, b, chatID, userID, c, state)
		case "pb-save":
			id, err := c.CacheService.SavePostBuilderSession(ctx, *state)
			if err != nil {
				logger.Error("BOT", "PostBuilder: Error saving session: %v", err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "❌ Erro ao salvar postagem.",
				})
				return
			}
			botInfo, _ := b.GetMe(ctx)

			kb := &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "🚀 Compartilhar", SwitchInlineQuery: "pb " + id},
					},
				},
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      chatID,
				Text:        fmt.Sprintf("✅ <b>Postagem salva com sucesso!</b>\n\nUtilize o modo inline para enviar:\n<code>@%s pb %s</code>", botInfo.Username, id),
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
			c.CacheService.DeletePostBuilderState(ctx, userID)
		case "pb-cancel":
			c.CacheService.DeletePostBuilderState(ctx, userID)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Post Builder cancelado.",
			})
		}
	}
}

func sendFinalPost(ctx context.Context, b *bot.Bot, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState, deleteState bool) {
	caption := ""
	if state.Title != "" {
		caption += state.Title + "\n\n"
	}
	if state.Body != "" {
		caption += state.Body + "\n\n"
	}
	if state.Footer != "" {
		caption += state.Footer
	}

	var kb models.ReplyMarkup
	if len(state.Buttons) > 0 || state.Reactions != "" {
		ikb := &models.InlineKeyboardMarkup{}
		for _, btn := range state.Buttons {
			ikb.InlineKeyboard = append(ikb.InlineKeyboard, []models.InlineKeyboardButton{
				{Text: btn.Text, URL: btn.URL},
			})
		}
		
		if state.Reactions != "" {
			reactions := strings.Split(state.Reactions, ",")
			var reactionRow []models.InlineKeyboardButton
			for _, r := range reactions {
				emoji := strings.TrimSpace(r)
				if emoji != "" {
					reactionRow = append(reactionRow, models.InlineKeyboardButton{
						Text:         emoji,
						CallbackData: "vote:" + emoji,
					})
				}
			}
			if len(reactionRow) > 0 {
				ikb.InlineKeyboard = append(ikb.InlineKeyboard, reactionRow)
			}
		}
		
		kb = ikb
	}

	switch state.MediaType {
	case "photo":
		_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      chatID,
			Photo:       &models.InputFileString{Data: state.MediaFileID},
			Caption:     caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending photo: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar foto: %v", err),
			})
		}
	case "video":
		_, err := b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:      chatID,
			Video:       &models.InputFileString{Data: state.MediaFileID},
			Caption:     caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending video: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar vídeo: %v", err),
			})
		}
	case "animation":
		_, err := b.SendAnimation(ctx, &bot.SendAnimationParams{
			ChatID:      chatID,
			Animation:   &models.InputFileString{Data: state.MediaFileID},
			Caption:     caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending animation: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar animação: %v", err),
			})
		}
	case "audio":
		_, err := b.SendAudio(ctx, &bot.SendAudioParams{
			ChatID:      chatID,
			Audio:       &models.InputFileString{Data: state.MediaFileID},
			Caption:     caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending audio: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar áudio: %v", err),
			})
		}
	case "document":
		_, err := b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:      chatID,
			Document:    &models.InputFileString{Data: state.MediaFileID},
			Caption:     caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending document: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar documento: %v", err),
			})
		}
	default:
		logger.Warn("BOT", "PostBuilder: No media type matched, sending text only. Type: %s", state.MediaType)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        caption,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
	}

	if deleteState {
		c.CacheService.DeletePostBuilderState(ctx, userID)
	}
}

func InlineHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.InlineQuery == nil {
			return
		}

		query := update.InlineQuery.Query
		if !strings.HasPrefix(query, "pb ") {
			return
		}

		id := strings.TrimSpace(strings.TrimPrefix(query, "pb "))
		if id == "" {
			return
		}

		state, err := c.CacheService.GetPostBuilderSession(ctx, id)
		if err != nil || state == nil {
			logger.Warn("BOT", "PostBuilder: InlineQuery session %s not found", id)
			_, _ = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
				InlineQueryID: update.InlineQuery.ID,
				Results: []models.InlineQueryResult{
					&models.InlineQueryResultArticle{
						ID:    "not_found",
						Title: "❌ Postagem não encontrada",
						InputMessageContent: &models.InputTextMessageContent{
							MessageText: "Esta postagem não existe ou já expirou.",
						},
					},
				},
				CacheTime: 0,
			})
			return
		}

		caption := ""
		if state.Title != "" {
			caption += state.Title + "\n\n"
		}
		if state.Body != "" {
			caption += state.Body + "\n\n"
		}
		if state.Footer != "" {
			caption += state.Footer
		}

		// Garante que o caption não seja vazio para o caso de Article
		displayCaption := caption
		if displayCaption == "" {
			displayCaption = "Postagem sem texto."
		}

		var kb *models.InlineKeyboardMarkup
		if len(state.Buttons) > 0 || state.Reactions != "" {
			kb = &models.InlineKeyboardMarkup{}
			for _, btn := range state.Buttons {
				kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
					{Text: btn.Text, URL: btn.URL},
				})
			}

			if state.Reactions != "" {
				reactions := strings.Split(state.Reactions, ",")
				var reactionRow []models.InlineKeyboardButton
				for _, r := range reactions {
					emoji := strings.TrimSpace(r)
					if emoji != "" {
						reactionRow = append(reactionRow, models.InlineKeyboardButton{
							Text:         emoji,
							CallbackData: "vote:" + emoji,
						})
					}
				}
				if len(reactionRow) > 0 {
					kb.InlineKeyboard = append(kb.InlineKeyboard, reactionRow)
				}
			}
		}

		var result models.InlineQueryResult
		switch state.MediaType {
		case "photo":
			result = &models.InlineQueryResultCachedPhoto{
				ID:          id,
				PhotoFileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			}
		case "video":
			result = &models.InlineQueryResultCachedVideo{
				ID:          id,
				VideoFileID: state.MediaFileID,
				Title:       "Video Post",
				Caption:     caption,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			}
		case "animation":
			result = &models.InlineQueryResultCachedMpeg4Gif{
				ID:          id,
				Mpeg4FileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			}
		case "audio":
			result = &models.InlineQueryResultCachedAudio{
				ID:          id,
				AudioFileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			}
		case "document":
			result = &models.InlineQueryResultCachedDocument{
				ID:             id,
				DocumentFileID: state.MediaFileID,
				Title:          "Document Post",
				Caption:        caption,
				ParseMode:      models.ParseModeHTML,
				ReplyMarkup:    kb,
			}
		default:
			result = &models.InlineQueryResultArticle{
				ID:    id,
				Title: "Post Builder",
				InputMessageContent: &models.InputTextMessageContent{
					MessageText: displayCaption,
					ParseMode:   models.ParseModeHTML,
				},
				ReplyMarkup: kb,
			}
		}

		_, err = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results:       []models.InlineQueryResult{result},
			CacheTime:     0,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error answering inline query: %v", err)
		}
	}
}
