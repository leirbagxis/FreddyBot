package postbuilder

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
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

func HandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		if update.Message == nil || update.Message.From == nil {
			return nil
		}

		bot := ctx.Bot()

		// Check Blacklist
		user, err := c.UserService.GetUserByID(context.Background(), update.Message.From.ID)
		if err == nil && user != nil && user.IsBlacklisted {
			return nil
		}

		// Detect media
		var mediaID string
		var mediaType string

		// Se o usuário estiver configurando um sticker separador, o PostBuilder não deve interceptar
		awaitingStickerChannel, _ := c.CacheService.GetAwaitingStickerSeparator(context.Background(), update.Message.From.ID)
		if awaitingStickerChannel != 0 && update.Message.Sticker != nil {
			return nil
		}

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
			state, _ := c.CacheService.GetPostBuilderState(context.Background(), update.Message.From.ID)
			if state != nil && state.Step != "" {
				return handleTextInputTelego(ctx, update, c, state)
			}
			return nil
		}

		// Media detected, offer Post Builder
		kb := &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
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
		c.CacheService.SetPostBuilderState(context.Background(), update.Message.From.ID, state)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      "✨ Media detectada! Deseja usar o <b>Post Builder</b> para criar uma postagem personalizada?",
			ParseMode: telego.ModeHTML,
			ReplyMarkup: kb,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})

		return nil
	}
}

func handleTextInputTelego(ctx *telegohandler.Context, update telego.Update, c *container.AppContainer, state *cache.PostBuilderState) error {
	text := update.Message.Text
	bot := ctx.Bot()

	// Processar texto com entidades
	formattedText := channelpost.ProcessTextWithFormattingTelego(text, update.Message.Entities)

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
		parts := strings.Split(text, ",")
		var finalReactions []string
		valid := true

		// Mapear entidades para facilitar a busca por posição
		entityMap := make(map[int]string)
		for _, e := range update.Message.Entities {
			if e.Type == "custom_emoji" {
				entityMap[e.Offset] = e.CustomEmojiID
			}
		}

		currentOffset := 0
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed == "" {
				currentOffset += len(p) + 1
				continue
			}

			// Verificar se nesta posição do texto original existe uma entidade de emoji customizado
			pos := strings.Index(text[currentOffset:], trimmed) + currentOffset
			if eid, ok := entityMap[pos]; ok {
				finalReactions = append(finalReactions, "eid:"+eid)
			} else if isEmoji(trimmed) {
				finalReactions = append(finalReactions, trimmed)
			} else {
				valid = false
				break
			}
			currentOffset += len(p) + 1
		}

		if !valid {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      "❌ Apenas emojis (padrão ou customizados) são permitidos como reações. Tente novamente:",
				ParseMode: telego.ModeHTML,
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
			return nil
		}

		state.Reactions = strings.Join(finalReactions, ",")
		state.Step = ""
	case "awaiting_button":
		lines := strings.Split(text, "\n")
		if len(lines) < 2 {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      "❌ Formato inválido. Envie o <b>Nome</b> em uma linha e o <b>Link</b> na linha de baixo.",
				ParseMode: telego.ModeHTML,
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
			return nil
		}
		name := strings.TrimSpace(lines[0])
		url := strings.TrimSpace(lines[1])

		// Extrair CustomEmojiID do nome (primeira linha)
		var customEmojiID string
		firstLineLen := len(lines[0])
		for _, entity := range update.Message.Entities {
			if entity.Type == "custom_emoji" && entity.Offset < firstLineLen {
				customEmojiID = entity.CustomEmojiID
				break
			}
		}

		if !strings.HasPrefix(url, "http") {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text:   "❌ URL inválida. Deve começar com http:// ou https://. Tente novamente:",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
			return nil
		}

		state.Buttons = append(state.Buttons, cache.PostBuilderButton{Text: name, URL: url, CustomEmojiID: customEmojiID})
		state.Step = ""
	default:
		return nil
	}

	state.MenuMessageID = 0
	c.CacheService.SetPostBuilderState(context.Background(), update.Message.From.ID, *state)

	if state.Step == "awaiting_button" {
		showButtonManagerTelego(ctx, update.Message.Chat.ID, update.Message.From.ID, c, state)
	} else {
		showMenuTelego(ctx, update.Message.Chat.ID, update.Message.From.ID, c, state, update.Message.MessageID)
	}

	return nil
}

func showMenuTelego(ctx *telegohandler.Context, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState, replyToMessageID ...int) {
	var sb strings.Builder
	bot := ctx.Bot()

	sb.WriteString("🛠️ <b>Post Builder - Menu</b>\n\n")
	sb.WriteString(fmt.Sprintf("📝 <b>Título:</b> %s\n", state.Title))
	sb.WriteString(fmt.Sprintf("📄 <b>Corpo:</b> %s\n", state.Body))
	sb.WriteString(fmt.Sprintf("👣 <b>Rodapé:</b> %s\n", state.Footer))
	// Formatar reações para exibição no menu
	displayReactions := state.Reactions
	if strings.Contains(displayReactions, "eid:") {
		parts := strings.Split(displayReactions, ",")
		for i, p := range parts {
			if strings.HasPrefix(p, "eid:") {
				parts[i] = "🖼️" // Placeholder para emoji customizado
			}
		}
		displayReactions = strings.Join(parts, ", ")
	}

	sb.WriteString(fmt.Sprintf("🎭 <b>Reações:</b> %s\n", displayReactions))
	sb.WriteString(fmt.Sprintf("🔘 <b>Botões:</b> %d\n\n", len(state.Buttons)))
	sb.WriteString("Escolha o que deseja editar:")

	kb := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
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

	if len(replyToMessageID) > 0 && replyToMessageID[0] != 0 {
		msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: chatID},
			Text:      sb.String(),
			ParseMode: telego.ModeHTML,
			ReplyMarkup: kb,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: replyToMessageID[0],
			},
		})
		if msg != nil {
			state.MenuMessageID = msg.MessageID
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		}
		return
	}

	if state.MenuMessageID != 0 {
		_, err := bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: chatID},
			MessageID: state.MenuMessageID,
			Text:      sb.String(),
			ParseMode: telego.ModeHTML,
			ReplyMarkup: kb,
		})
		if err == nil {
			return
		}
	}

	msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		Text:      sb.String(),
		ParseMode: telego.ModeHTML,
		ReplyMarkup: kb,
	})

	if msg != nil {
		state.MenuMessageID = msg.MessageID
		c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
	}
}

func showButtonManagerTelego(ctx *telegohandler.Context, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState) {
	var sb strings.Builder
	bot := ctx.Bot()

	sb.WriteString("🔘 <b>Gerenciamento de Botões</b>\n\n")
	if len(state.Buttons) == 0 {
		sb.WriteString("<i>Nenhum botão adicionado ainda.</i>")
	} else {
		sb.WriteString("Clique em um botão para <b>excluí-lo</b>:")
	}

	var rows [][]telego.InlineKeyboardButton

	for i, btn := range state.Buttons {
		rows = append(rows, []telego.InlineKeyboardButton{
			{Text: "❌ " + btn.Text, CallbackData: fmt.Sprintf("pb-del-button:%d", i)},
		})
	}

	rows = append(rows, []telego.InlineKeyboardButton{
		{Text: "➕ Adicionar Novo Botão", CallbackData: "pb-add-button"},
	})
	rows = append(rows, []telego.InlineKeyboardButton{
		{Text: "🔙 Voltar ao Menu", CallbackData: "pb-start"},
	})

	kb := &telego.InlineKeyboardMarkup{InlineKeyboard: rows}

	if state.MenuMessageID != 0 {
		_, err := bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: chatID},
			MessageID: state.MenuMessageID,
			Text:      sb.String(),
			ParseMode: telego.ModeHTML,
			ReplyMarkup: kb,
		})
		if err == nil {
			return
		}
	}

	msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		Text:      sb.String(),
		ParseMode: telego.ModeHTML,
		ReplyMarkup: kb,
	})

	if msg != nil {
		state.MenuMessageID = msg.MessageID
		c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
	}
}

func CallbackHandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ID
		data := update.CallbackQuery.Data
		bot := ctx.Bot()

		user, err := c.UserService.GetUserByID(context.Background(), userID)
		if err == nil && user != nil && user.IsBlacklisted {
			_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Você está na blacklist.",
				ShowAlert:       true,
			})
			return nil
		}

		state, _ := c.CacheService.GetPostBuilderState(context.Background(), userID)
		if state == nil && data != "pb-cancel" && !strings.HasPrefix(data, "pb-send-") {
			_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Sessão expirada ou não encontrada.",
			})
			return nil
		}

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		if strings.HasPrefix(data, "pb-import-apply:") {
			channelIDStr := strings.TrimPrefix(data, "pb-import-apply:")
			channelID, _ := strconv.ParseInt(channelIDStr, 10, 64)

			channel, err := c.ChannelService.GetChannelWithRelations(context.Background(), channelID)
			if err != nil {
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatID},
					Text:   "❌ Erro ao obter dados do canal.",
				})
				return nil
			}

			if channel.DefaultCaption != nil {
				state.Body = channelpost.DetectParseMode(channel.DefaultCaption.Caption)
			}
			state.Reactions = channel.Reactions

			state.Buttons = make([]cache.PostBuilderButton, 0)
			for _, btn := range channel.Buttons {
				state.Buttons = append(state.Buttons, cache.PostBuilderButton{
					Text: btn.NameButton,
					URL:  btn.ButtonURL,
				})
			}

			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      fmt.Sprintf("✅ Dados importados do canal <b>%s</b>!", channel.Title),
				ParseMode: telego.ModeHTML,
			})
			showMenuTelego(ctx, chatID, userID, c, state)
			return nil
		}

		if strings.HasPrefix(data, "pb-del-button:") {
			indexStr := strings.TrimPrefix(data, "pb-del-button:")
			index, _ := strconv.Atoi(indexStr)

			if index >= 0 && index < len(state.Buttons) {
				state.Buttons = append(state.Buttons[:index], state.Buttons[index+1:]...)
				c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
			}
			showButtonManagerTelego(ctx, chatID, userID, c, state)
			return nil
		}

		if strings.HasPrefix(data, "pb-send-to-channels:") {
			sessionID := strings.TrimPrefix(data, "pb-send-to-channels:")
			handleSendToChannelsTelego(ctx, chatID, userID, sessionID, c)
			return nil
		}

		if strings.HasPrefix(data, "pb-send-apply:") {
			parts := strings.Split(strings.TrimPrefix(data, "pb-send-apply:"), ":")
			if len(parts) == 2 {
				channelID, _ := strconv.ParseInt(parts[0], 10, 64)
				sessionID := parts[1]
				handleSendApplyTelego(ctx, chatID, userID, channelID, sessionID, c)
			}
			return nil
		}

		switch data {
		case "pb-start":
			showMenuTelego(ctx, chatID, userID, c, state)
		case "pb-manage-buttons":
			showButtonManagerTelego(ctx, chatID, userID, c, state)
		case "pb-import-channel":
			channels, err := c.ChannelService.GetUserChannels(context.Background(), userID)
			if err != nil || len(channels) == 0 {
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID:    telego.ChatID{ID: chatID},
					Text:      "❌ Você não possui canais cadastrados ou ocorreu um erro.",
					ParseMode: telego.ModeHTML,
				})
				return nil
			}

			var rows [][]telego.InlineKeyboardButton
			for _, ch := range channels {
				rows = append(rows, []telego.InlineKeyboardButton{
					{Text: "📣 " + ch.Title, CallbackData: fmt.Sprintf("pb-import-apply:%d", ch.ID)},
				})
			}
			rows = append(rows, []telego.InlineKeyboardButton{
				{Text: "🔙 Voltar", CallbackData: "pb-start"},
			})

			kb := &telego.InlineKeyboardMarkup{InlineKeyboard: rows}
			text := "📥 <b>Importar de Canal</b>\n\nEscolha o canal de onde deseja copiar a legenda padrão, reações e botões:"

			if state.MenuMessageID != 0 {
				_, err := bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
					ChatID:      telego.ChatID{ID: chatID},
					MessageID:   state.MenuMessageID,
					Text:        text,
					ParseMode:   telego.ModeHTML,
					ReplyMarkup: kb,
				})
				if err == nil {
					return nil
				}
			}

			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      text,
				ParseMode: telego.ModeHTML,
				ReplyMarkup: kb,
			})
		case "pb-edit-title":
			state.Step = "awaiting_title"
			msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      "📝 Envie o <b>Título</b> da postagem (suporta formatação):",
				ParseMode: telego.ModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.MessageID
			}
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		case "pb-edit-body":
			state.Step = "awaiting_body"
			msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      "📄 Envie o <b>Corpo</b> da postagem (suporta formatação):",
				ParseMode: telego.ModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.MessageID
			}
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		case "pb-edit-footer":
			state.Step = "awaiting_footer"
			msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      "👣 Envie o <b>Rodapé</b> da postagem (suporta formatação):",
				ParseMode: telego.ModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.MessageID
			}
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		case "pb-edit-reactions":
			state.Step = "awaiting_reactions"
			msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      "🎭 Envie as <b>Reações</b> separadas por vírgula (ex: 👍,👎,❤️):",
				ParseMode: telego.ModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.MessageID
			}
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		case "pb-add-button":
			state.Step = "awaiting_button"
			msg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      "🔘 Envie os dados do botão no formato:\n\n<code>Nome do Botão\nhttps://link.com</code>",
				ParseMode: telego.ModeHTML,
			})
			if msg != nil {
				state.PromptMessageID = msg.MessageID
			}
			c.CacheService.SetPostBuilderState(context.Background(), userID, *state)
		case "pb-preview":
			sendFinalPostTelego(ctx, chatID, userID, c, state, false)
			showMenuTelego(ctx, chatID, userID, c, state)
		case "pb-save":
			id, err := c.CacheService.SavePostBuilderSession(context.Background(), *state)
			if err != nil {
				logger.Error("BOT", "PostBuilder: Error saving session: %v", err)
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatID},
					Text:   "❌ Erro ao salvar postagem.",
				})
				return nil
			}
			botInfo, _ := bot.GetMe(context.Background())

			query := "pb " + id
			kb := &telego.InlineKeyboardMarkup{
				InlineKeyboard: [][]telego.InlineKeyboardButton{
					{
						{Text: "🚀 Compartilhar", SwitchInlineQuery: &query},
					},
					{
						{Text: "📢 Enviar para Canais", CallbackData: "pb-send-to-channels:" + id},
					},
				},
			}

			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chatID},
				Text:      fmt.Sprintf("✅ <b>Postagem salva com sucesso!</b>\n\nUtilize o modo inline para enviar:\n<code>@%s pb %s</code>", botInfo.Username, id),
				ParseMode: telego.ModeHTML,
				ReplyMarkup: kb,
			})
			c.CacheService.DeletePostBuilderState(context.Background(), userID)
		case "pb-cancel":
			c.CacheService.DeletePostBuilderState(context.Background(), userID)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: chatID},
				Text:   "❌ Post Builder cancelado.",
			})
		}

		return nil
	}
}

func handleSendToChannelsTelego(ctx *telegohandler.Context, chatID, userID int64, sessionID string, c *container.AppContainer) {
	bot := ctx.Bot()
	channels, err := c.ChannelService.GetUserChannels(context.Background(), userID)
	if err != nil || len(channels) == 0 {
		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: chatID},
			Text:      "❌ Você não possui canais cadastrados para envio direto.",
			ParseMode: telego.ModeHTML,
		})
		return
	}

	var rows [][]telego.InlineKeyboardButton
	for _, ch := range channels {
		rows = append(rows, []telego.InlineKeyboardButton{
			{Text: "📣 " + ch.Title, CallbackData: fmt.Sprintf("pb-send-apply:%d:%s", ch.ID, sessionID)},
		})
	}

	kb := &telego.InlineKeyboardMarkup{InlineKeyboard: rows}
	_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		Text:      "📢 <b>Enviar para Canal</b>\n\nSelecione o canal para o qual deseja enviar esta postagem:",
		ParseMode: telego.ModeHTML,
		ReplyMarkup: kb,
	})
}

func handleSendApplyTelego(ctx *telegohandler.Context, chatID, userID int64, channelID int64, sessionID string, c *container.AppContainer) {
	bot := ctx.Bot()
	state, err := c.CacheService.GetPostBuilderSession(context.Background(), sessionID)
	if err != nil || state == nil {
		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "❌ Sessão de postagem não encontrada ou expirada.",
		})
		return
	}

	sendFinalPostTelego(ctx, channelID, userID, c, state, false)

	_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: chatID},
		Text:   "✅ Postagem enviada com sucesso para o canal!",
	})
}

func sendFinalPostTelego(ctx *telegohandler.Context, chatID, userID int64, c *container.AppContainer, state *cache.PostBuilderState, deleteState bool) {
	var sb strings.Builder
	bot := ctx.Bot()

	if state.Title != "" {
		sb.WriteString(state.Title + "\n\n")
	}
	if state.Body != "" {
		sb.WriteString(state.Body + "\n\n")
	}
	if state.Footer != "" {
		sb.WriteString(state.Footer)
	}
	caption := sb.String()

	// Safeguard: se houver Markdown não convertido e não contiver tags HTML básicas
	if channelpost.IsMarkdown(caption) && !strings.Contains(caption, "<a href=") && !strings.Contains(caption, "<b>") && !strings.Contains(caption, "<tg-emoji") {
		caption = channelpost.DetectParseMode(caption)
	}

	var kb telego.ReplyMarkup
	if len(state.Buttons) > 0 || state.Reactions != "" {
		ikb := &telego.InlineKeyboardMarkup{}
		for _, btn := range state.Buttons {
			ikb.InlineKeyboard = append(ikb.InlineKeyboard, []telego.InlineKeyboardButton{
				{Text: btn.Text, URL: btn.URL, IconCustomEmojiID: btn.CustomEmojiID},
			})
		}

		if state.Reactions != "" {
			reactions := strings.Split(state.Reactions, ",")
			var reactionRow []telego.InlineKeyboardButton
			for _, r := range reactions {
				val := strings.TrimSpace(r)
				if val != "" {
					btn := telego.InlineKeyboardButton{
						CallbackData: "vote:" + val,
					}
					if strings.HasPrefix(val, "eid:") {
						btn.IconCustomEmojiID = strings.TrimPrefix(val, "eid:")
						btn.Text = " " // Texto mínimo para botões com ícone
					} else {
						btn.Text = val
					}
					reactionRow = append(reactionRow, btn)
				}
			}
			if len(reactionRow) > 0 {
				ikb.InlineKeyboard = append(ikb.InlineKeyboard, reactionRow)
			}
		}
		kb = ikb
	}

	paramsPhoto := &telego.SendPhotoParams{
		ChatID:    telego.ChatID{ID: chatID},
		Photo:     telego.InputFile{FileID: state.MediaFileID},
		Caption:   caption,
		ParseMode: telego.ModeHTML,
	}
	if kb != nil {
		paramsPhoto.ReplyMarkup = kb
	}

	var err error
	switch state.MediaType {
	case "photo":
		_, err = bot.SendPhoto(context.Background(), paramsPhoto)
	case "video":
		params := &telego.SendVideoParams{
			ChatID:    telego.ChatID{ID: chatID},
			Video:     telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendVideo(context.Background(), params)
	case "animation":
		params := &telego.SendAnimationParams{
			ChatID:    telego.ChatID{ID: chatID},
			Animation: telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendAnimation(context.Background(), params)
	case "audio":
		params := &telego.SendAudioParams{
			ChatID:    telego.ChatID{ID: chatID},
			Audio:     telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendAudio(context.Background(), params)
	case "document":
		params := &telego.SendDocumentParams{
			ChatID:    telego.ChatID{ID: chatID},
			Document:  telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendDocument(context.Background(), params)
	case "sticker":
		params := &telego.SendStickerParams{
			ChatID:  telego.ChatID{ID: chatID},
			Sticker: telego.InputFile{FileID: state.MediaFileID},
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendSticker(context.Background(), params)
	default:
		params := &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: chatID},
			Text:      caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err = bot.SendMessage(context.Background(), params)
	}

	if err != nil {
		logger.Error("BOT", "PostBuilder: Error sending final post: %v", err)
	}

	if deleteState {
		c.CacheService.DeletePostBuilderState(context.Background(), userID)
	}
}

func ChosenInlineResultHandlerTelego(c *container.AppContainer) telegohandler.ChosenInlineResultHandler {
	return func(ctx *telegohandler.Context, result telego.ChosenInlineResult) error {
		sessionID := result.ResultID
		inlineMessageID := result.InlineMessageID

		if sessionID != "" && inlineMessageID != "" {
			key := fmt.Sprintf("pb_inline_map:%s", inlineMessageID)
			err := c.CacheService.Set(context.Background(), key, sessionID, 24*time.Hour)
			if err != nil {
				logger.Error("BOT", "❌ Erro ao salvar mapeamento inline no Redis: %v", err)
			}
		}

		return nil
	}
}

func InlineHandlerTelego(c *container.AppContainer) telegohandler.InlineQueryHandler {
	return func(ctx *telegohandler.Context, inlineQuery telego.InlineQuery) error {
		bot := ctx.Bot()
		query := inlineQuery.Query
		if !strings.HasPrefix(query, "pb ") {
			return nil
		}

		id := strings.TrimSpace(strings.TrimPrefix(query, "pb "))
		if id == "" {
			return nil
		}

		state, err := c.CacheService.GetPostBuilderSession(context.Background(), id)
		if err != nil || state == nil {
			logger.Warn("BOT", "InlineHandler: Sessão %s não encontrada ou expirada", id)
			_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
				InlineQueryID: inlineQuery.ID,
				Results: []telego.InlineQueryResult{
					&telego.InlineQueryResultArticle{
						Type:  "article",
						ID:    "not_found",
						Title: "❌ Postagem não encontrada",
						InputMessageContent: &telego.InputTextMessageContent{
							MessageText: "Esta postagem não existe ou já expirou.",
						},
					},
				},
				CacheTime: 0,
			})
			return nil
		}

		var sb strings.Builder
		if state.Title != "" {
			sb.WriteString(state.Title + "\n\n")
		}
		if state.Body != "" {
			sb.WriteString(state.Body + "\n\n")
		}
		if state.Footer != "" {
			sb.WriteString(state.Footer)
		}
		caption := sb.String()

		// Safeguard: se houver Markdown não convertido e não contiver tags HTML básicas (links/bold)
		if channelpost.IsMarkdown(caption) && !strings.Contains(caption, "<a href=") && !strings.Contains(caption, "<b>") && !strings.Contains(caption, "<tg-emoji") {
			caption = channelpost.DetectParseMode(caption)
		}

		displayCaption := caption
		if displayCaption == "" {
			displayCaption = "Postagem sem texto."
		}

		var kb *telego.InlineKeyboardMarkup
		if len(state.Buttons) > 0 || state.Reactions != "" {
			kb = &telego.InlineKeyboardMarkup{}
			for _, btn := range state.Buttons {
				kb.InlineKeyboard = append(kb.InlineKeyboard, []telego.InlineKeyboardButton{
					{Text: btn.Text, URL: btn.URL, IconCustomEmojiID: btn.CustomEmojiID},
				})
			}

			if state.Reactions != "" {
				reactions := strings.Split(state.Reactions, ",")
				var reactionRow []telego.InlineKeyboardButton
				for _, r := range reactions {
					val := strings.TrimSpace(r)
					if val != "" {
						btn := telego.InlineKeyboardButton{
							CallbackData: "vote:" + val,
						}
						if strings.HasPrefix(val, "eid:") {
							btn.IconCustomEmojiID = strings.TrimPrefix(val, "eid:")
							btn.Text = " "
						} else {
							btn.Text = val
						}
						reactionRow = append(reactionRow, btn)
					}
				}
				if len(reactionRow) > 0 {
					kb.InlineKeyboard = append(kb.InlineKeyboard, reactionRow)
				}
			}
		}

		var result telego.InlineQueryResult
		switch state.MediaType {
		case "photo":
			res := &telego.InlineQueryResultCachedPhoto{
				Type:        "photo",
				ID:          id,
				PhotoFileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   telego.ModeHTML,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		case "video":
			res := &telego.InlineQueryResultCachedVideo{
				Type:        "video",
				ID:          id,
				VideoFileID: state.MediaFileID,
				Title:       "Video Post",
				Caption:     caption,
				ParseMode:   telego.ModeHTML,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		case "animation":
			res := &telego.InlineQueryResultCachedMpeg4Gif{
				Type:        "mpeg4_gif",
				ID:          id,
				Mpeg4FileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   telego.ModeHTML,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		case "audio":
			res := &telego.InlineQueryResultCachedAudio{
				Type:        "audio",
				ID:          id,
				AudioFileID: state.MediaFileID,
				Caption:     caption,
				ParseMode:   telego.ModeHTML,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		case "document":
			res := &telego.InlineQueryResultCachedDocument{
				Type:           "document",
				ID:             id,
				DocumentFileID: state.MediaFileID,
				Title:          "Document Post",
				Caption:        caption,
				ParseMode:      telego.ModeHTML,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		case "sticker":
			res := &telego.InlineQueryResultCachedSticker{
				Type:          "sticker",
				ID:            id,
				StickerFileID: state.MediaFileID,
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		default:
			res := &telego.InlineQueryResultArticle{
				Type:  "article",
				ID:    id,
				Title: "Post Builder",
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: displayCaption,
					ParseMode:   telego.ModeHTML,
				},
			}
			if kb != nil {
				res.ReplyMarkup = kb
			}
			result = res
		}

		if err := bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
			InlineQueryID: inlineQuery.ID,
			Results:       []telego.InlineQueryResult{result},
			CacheTime:     0,
		}); err != nil {
			logger.Error("BOT", "Erro ao responder Inline Query: %v", err)
		}

		return nil
	}
}
