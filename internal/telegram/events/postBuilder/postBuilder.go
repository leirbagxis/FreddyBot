package postbuilder

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
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
		user, err := c.UserService.GetUserByID(ctx, update.Message.From.ID)
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
		} else if update.Message.Sticker != nil {
			mediaID = update.Message.Sticker.FileID
			mediaType = "sticker"
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
	// Processar texto com entidades (negrito, itálico, links) e markdown
	formattedText := channelpost.ProcessTextWithFormatting(text, update.Message.Entities)

	// Reset prompt ID if it exists (but don't delete)
	if state.PromptMessageID != 0 {
		state.PromptMessageID = 0
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
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		state.Reactions = text
		state.Step = ""
	case "awaiting_button":
		lines := strings.Split(text, "\n")
		if len(lines) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "❌ Formato inválido. Envie o <b>Nome</b> em uma linha e o <b>Link</b> na linha de baixo.",
				ParseMode: models.ParseModeHTML,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}
		name := strings.TrimSpace(lines[0])
		url := strings.TrimSpace(lines[1])

		if !strings.HasPrefix(url, "http") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ URL inválida. Deve começar com http:// ou https://. Tente novamente:",
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		state.Buttons = append(state.Buttons, cache.PostBuilderButton{Text: name, URL: url})
		state.Step = ""
	default:
		return
	}

	// Always clear MenuMessageID when re-sending menu after input
	state.MenuMessageID = 0
	c.CacheService.SetPostBuilderState(ctx, update.Message.From.ID, *state)

	if state.Step == "awaiting_button" {
		showButtonManager(ctx, b, update.Message.Chat.ID, update.Message.From.ID, c, state)
	} else {
		// Enviar novo menu respondendo à mensagem do usuário
		showMenu(ctx, b, update.Message.Chat.ID, update.Message.From.ID, c, state, update.Message.ID)
	}
}

func showMenu(ctx context.Context, b *bot.Bot, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState, replyToMessageID ...int) {
	var sb strings.Builder
	sb.WriteString("🛠️ <b>Post Builder - Menu</b>\n\n")
	sb.WriteString(fmt.Sprintf("📝 <b>Título:</b> %s\n", state.Title))
	sb.WriteString(fmt.Sprintf("📄 <b>Corpo:</b> %s\n", state.Body))
	sb.WriteString(fmt.Sprintf("👣 <b>Rodapé:</b> %s\n", state.Footer))
	sb.WriteString(fmt.Sprintf("🎭 <b>Reações:</b> %s\n", state.Reactions))
	sb.WriteString(fmt.Sprintf("🔘 <b>Botões:</b> %d\n\n", len(state.Buttons)))
	sb.WriteString("Escolha o que deseja editar:")

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
				{Text: "🔘 Botões", CallbackData: "pb-manage-buttons"},
				{Text: "📥 Importar Canal", CallbackData: "pb-import-channel"},
			},
			{
				{Text: "👁️ Preview", CallbackData: "pb-preview"},
			},
			{
				{Text: "✅ Salvar", CallbackData: "pb-save"},
				{Text: "❌ Cancelar", CallbackData: "pb-cancel"},
			},
		},
	}

	// Se houver um ID para responder, ignoramos o Edit e enviamos uma nova mensagem
	if len(replyToMessageID) > 0 && replyToMessageID[0] != 0 {
		params := &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        sb.String(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
			ReplyParameters: &models.ReplyParameters{
				MessageID: replyToMessageID[0],
			},
		}
		msg, _ := b.SendMessage(ctx, params)
		if msg != nil {
			state.MenuMessageID = msg.ID
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
		}
		return
	}

	if state.MenuMessageID != 0 {
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   state.MenuMessageID,
			Text:        sb.String(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err == nil {
			return
		}
	}

	msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        sb.String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: kb,
	})

	if msg != nil {
		state.MenuMessageID = msg.ID
		c.CacheService.SetPostBuilderState(ctx, userID, *state)
	}
}

func showButtonManager(ctx context.Context, b *bot.Bot, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState) {
	var sb strings.Builder
	sb.WriteString("🔘 <b>Gerenciamento de Botões</b>\n\n")
	if len(state.Buttons) == 0 {
		sb.WriteString("<i>Nenhum botão adicionado ainda.</i>")
	} else {
		sb.WriteString("Clique em um botão para <b>excluí-lo</b>:")
	}

	var rows [][]models.InlineKeyboardButton

	// Listar botões atuais para exclusão
	for i, btn := range state.Buttons {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "❌ " + btn.Text, CallbackData: fmt.Sprintf("pb-del-button:%d", i)},
		})
	}

	// Botões de ação
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: "➕ Adicionar Novo Botão", CallbackData: "pb-add-button"},
	})
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: "🔙 Voltar ao Menu", CallbackData: "pb-start"},
	})

	kb := &models.InlineKeyboardMarkup{InlineKeyboard: rows}

	if state.MenuMessageID != 0 {
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   state.MenuMessageID,
			Text:        sb.String(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err == nil {
			return
		}
	}

	msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        sb.String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: kb,
	})

	if msg != nil {
		state.MenuMessageID = msg.ID
		c.CacheService.SetPostBuilderState(ctx, userID, *state)
	}
}

func CallbackHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.Message.Chat.ID
		data := update.CallbackQuery.Data

		// Check Blacklist
		user, err := c.UserService.GetUserByID(ctx, userID)
		if err == nil && user != nil && user.IsBlacklisted {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você está na blacklist.",
				ShowAlert:       true,
			})
			return
		}

		state, _ := c.CacheService.GetPostBuilderState(ctx, userID)
		if state == nil && data != "pb-cancel" && !strings.HasPrefix(data, "pb-send-") {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Sessão expirada ou não encontrada.",
			})
			return
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		// --- Handlers Dinâmicos (Prefixos) ---
		if strings.HasPrefix(data, "pb-import-apply:") {
			channelIDStr := strings.TrimPrefix(data, "pb-import-apply:")
			channelID, _ := strconv.ParseInt(channelIDStr, 10, 64)

			channel, err := c.ChannelService.GetChannelWithRelations(ctx, channelID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "❌ Erro ao obter dados do canal.",
				})
				return
			}

			// Mapear dados do canal para o state
			if channel.DefaultCaption != nil {
				state.Body = channel.DefaultCaption.Caption
			}
			state.Reactions = channel.Reactions

			// Mapear botões
			state.Buttons = make([]cache.PostBuilderButton, 0)
			for _, btn := range channel.Buttons {
				state.Buttons = append(state.Buttons, cache.PostBuilderButton{
					Text: btn.NameButton,
					URL:  btn.ButtonURL,
				})
			}

			c.CacheService.SetPostBuilderState(ctx, userID, *state)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      fmt.Sprintf("✅ Dados importados do canal <b>%s</b>!", channel.Title),
				ParseMode: models.ParseModeHTML,
			})
			showMenu(ctx, b, chatID, userID, c, state)
			return
		}

		if strings.HasPrefix(data, "pb-del-button:") {
			indexStr := strings.TrimPrefix(data, "pb-del-button:")
			index, _ := strconv.Atoi(indexStr)

			if index >= 0 && index < len(state.Buttons) {
				// Remover item do slice
				state.Buttons = append(state.Buttons[:index], state.Buttons[index+1:]...)
				c.CacheService.SetPostBuilderState(ctx, userID, *state)
			}
			showButtonManager(ctx, b, chatID, userID, c, state)
			return
		}

		if strings.HasPrefix(data, "pb-send-to-channels:") {
			sessionID := strings.TrimPrefix(data, "pb-send-to-channels:")
			handleSendToChannels(ctx, b, chatID, userID, sessionID, c)
			return
		}

		if strings.HasPrefix(data, "pb-send-apply:") {
			parts := strings.Split(strings.TrimPrefix(data, "pb-send-apply:"), ":")
			if len(parts) == 2 {
				channelID, _ := strconv.ParseInt(parts[0], 10, 64)
				sessionID := parts[1]
				handleSendApply(ctx, b, chatID, userID, channelID, sessionID, c)
			}
			return
		}

		switch data {
		case "pb-start":
			showMenu(ctx, b, chatID, userID, c, state)
		case "pb-manage-buttons":
			showButtonManager(ctx, b, chatID, userID, c, state)
		case "pb-import-channel":
			channels, err := c.ChannelService.GetUserChannels(ctx, userID)
			if err != nil || len(channels) == 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:    chatID,
					Text:      "❌ Você não possui canais cadastrados ou ocorreu um erro.",
					ParseMode: models.ParseModeHTML,
				})
				return
			}

			var rows [][]models.InlineKeyboardButton
			for _, ch := range channels {
				rows = append(rows, []models.InlineKeyboardButton{
					{Text: "📣 " + ch.Title, CallbackData: fmt.Sprintf("pb-import-apply:%d", ch.ID)},
				})
			}
			rows = append(rows, []models.InlineKeyboardButton{
				{Text: "🔙 Voltar", CallbackData: "pb-start"},
			})

			kb := &models.InlineKeyboardMarkup{InlineKeyboard: rows}
			text := "📥 <b>Importar de Canal</b>\n\nEscolha o canal de onde deseja copiar a legenda padrão, reações e botões:"

			if state.MenuMessageID != 0 {
				_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:      chatID,
					MessageID:   state.MenuMessageID,
					Text:        text,
					ParseMode:   models.ParseModeHTML,
					ReplyMarkup: kb,
				})
				if err == nil {
					return
				}
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      chatID,
				Text:        text,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		case "pb-edit-title":
			state.Step = "awaiting_title"
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "📝 Envie o <b>Título</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.ID
			}
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
		case "pb-edit-body":
			state.Step = "awaiting_body"
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "📄 Envie o <b>Corpo</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.ID
			}
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
		case "pb-edit-footer":
			state.Step = "awaiting_footer"
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "👣 Envie o <b>Rodapé</b> da postagem (suporta formatação):",
				ParseMode: models.ParseModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.ID
			}
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
		case "pb-edit-reactions":
			state.Step = "awaiting_reactions"
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "🎭 Envie as <b>Reações</b> separadas por vírgula (ex: 👍,👎,❤️):",
				ParseMode: models.ParseModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.ID
			}
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
		case "pb-add-button":
			state.Step = "awaiting_button"
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "🔘 Envie os dados do botão no formato:\n\n<code>Nome do Botão\nhttps://link.com</code>",
				ParseMode: models.ParseModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.ID
			}
			c.CacheService.SetPostBuilderState(ctx, userID, *state)
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
					{
						{Text: "📢 Enviar para Canais", CallbackData: "pb-send-to-channels:" + id},
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

func handleSendToChannels(ctx context.Context, b *bot.Bot, chatID, userID int64, sessionID string, c *container.AppContainer) {
	channels, err := c.ChannelService.GetUserChannels(ctx, userID)
	if err != nil || len(channels) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      "❌ Você não possui canais cadastrados para envio direto.",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	var rows [][]models.InlineKeyboardButton
	for _, ch := range channels {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "📣 " + ch.Title, CallbackData: fmt.Sprintf("pb-send-apply:%d:%s", ch.ID, sessionID)},
		})
	}

	kb := &models.InlineKeyboardMarkup{InlineKeyboard: rows}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "📢 <b>Enviar para Canal</b>\n\nSelecione o canal para o qual deseja enviar esta postagem:",
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: kb,
	})
}

func handleSendApply(ctx context.Context, b *bot.Bot, chatID, userID int64, channelID int64, sessionID string, c *container.AppContainer) {
	state, err := c.CacheService.GetPostBuilderSession(ctx, sessionID)
	if err != nil || state == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Sessão de postagem não encontrada ou expirada.",
		})
		return
	}

	sendFinalPost(ctx, b, channelID, userID, c, state, false)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "✅ Postagem enviada com sucesso para o canal!",
	})
}

func sendFinalPost(ctx context.Context, b *bot.Bot, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState, deleteState bool) {
	var sb strings.Builder
	if state.Title != "" {
		sb.WriteString(channelpost.DetectParseMode(state.Title) + "\n\n")
	}
	if state.Body != "" {
		sb.WriteString(channelpost.DetectParseMode(state.Body) + "\n\n")
	}
	if state.Footer != "" {
		sb.WriteString(channelpost.DetectParseMode(state.Footer))
	}
	caption := sb.String()

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
	case "sticker":
		_, err := b.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:      chatID,
			Sticker:     &models.InputFileString{Data: state.MediaFileID},
			ReplyMarkup: kb,
		})
		if err != nil {
			logger.Error("BOT", "PostBuilder: Error sending sticker: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Erro ao enviar sticker: %v", err),
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

func ChosenInlineResultHandler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.ChosenInlineResult == nil {
			return
		}

		sessionID := update.ChosenInlineResult.ResultID
		inlineMessageID := update.ChosenInlineResult.InlineMessageID

		logger.Bot("📥 ChosenInlineResult recebido: SessionID=%s, InlineMessageID=%s", sessionID, inlineMessageID)

		if sessionID != "" && inlineMessageID != "" {
			// Mapeia o inline_message_id para o sessionID da postagem
			// Expira em 24h (mesmo tempo da sessão pb_session)
			key := fmt.Sprintf("pb_inline_map:%s", inlineMessageID)
			err := c.CacheService.Set(ctx, key, sessionID, 24*time.Hour)
			if err != nil {
				logger.Error("BOT", "❌ Erro ao salvar mapeamento inline no Redis: %v", err)
			} else {
				logger.Bot("🔗 Mapeamento salvo no Redis: %s -> %s", key, sessionID)
			}
		} else {
			logger.Warn("BOT", "⚠️ ChosenInlineResult com dados incompletos: SessionID=%s, InlineMessageID=%s", sessionID, inlineMessageID)
		}
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

		var sb strings.Builder
		if state.Title != "" {
			sb.WriteString(channelpost.DetectParseMode(state.Title) + "\n\n")
		}
		if state.Body != "" {
			sb.WriteString(channelpost.DetectParseMode(state.Body) + "\n\n")
		}
		if state.Footer != "" {
			sb.WriteString(channelpost.DetectParseMode(state.Footer))
		}
		caption := sb.String()

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
						})					}
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
		case "sticker":
			result = &models.InlineQueryResultCachedSticker{
				ID:            id,
				StickerFileID: state.MediaFileID,
				ReplyMarkup:   kb,
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
