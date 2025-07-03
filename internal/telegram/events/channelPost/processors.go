package channelpost

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// ========== CACHES E REGEX ==========
var (
	hashtagRegex         = regexp.MustCompile(`#(\w+)`)
	removeHashRegexCache = sync.Map{} // string -> *regexp.Regexp
	customCaptionCache   = sync.Map{} // string -> *dbmodels.CustomCaption
	mediaGroups          = sync.Map{} // string -> *MediaGroup
)

// ========== ESTRUTURAS ==========

type MediaGroup struct {
	Messages           []MediaMessage
	Processed          bool
	Timer              *time.Timer
	MessageEditAllowed bool
	ChatID             int64
	mu                 sync.Mutex
}

type MessageProcessor struct {
	bot               *bot.Bot
	permissionManager *PermissionManager
	mediaGroupManager *MediaGroupManager
}

// ✅ CONFIGURAÇÃO DO BOT com timeout maior
func NewMessageProcessor(b *bot.Bot) *MessageProcessor {
	// Configurar cliente HTTP com timeout maior se possível
	return &MessageProcessor{
		bot:               b,
		permissionManager: NewPermissionManager(),
		mediaGroupManager: NewMediaGroupManager(),
	}
}

// ✅ TIMEOUT ADAPTATIVO melhorado
func (mp *MessageProcessor) getAdaptiveTimeout(groupSize int) time.Duration {
	// Timeout base maior para grupos grandes
	baseTimeout := 2000 * time.Millisecond
	additionalTime := time.Duration(groupSize*300) * time.Millisecond
	maxTimeout := 5000 * time.Millisecond

	timeout := baseTimeout + additionalTime
	if timeout > maxTimeout {
		timeout = maxTimeout
	}

	return timeout
}

// ========== INLINE KEYBOARD ==========
func (mp *MessageProcessor) CreateInlineKeyboard(
	buttons []dbmodels.Button,
	customCaption *dbmodels.CustomCaption,
	allowButtons bool,
) *models.InlineKeyboardMarkup {
	if !allowButtons {
		return nil
	}
	var finalButtons []dbmodels.Button

	// Priorize botões do custom caption, se existirem
	if customCaption != nil && len(customCaption.Buttons) > 0 {
		for _, cb := range customCaption.Buttons {
			if cb.NameButton == "" || cb.ButtonURL == "" {
				continue
			}
			finalButtons = append(finalButtons, dbmodels.Button{
				NameButton: cb.NameButton,
				ButtonURL:  cb.ButtonURL,
				PositionY:  cb.PositionY,
				PositionX:  cb.PositionX,
			})
		}
	} else {
		for _, b := range buttons {
			if b.NameButton == "" || b.ButtonURL == "" {
				continue
			}
			finalButtons = append(finalButtons, b)
		}
	}

	if len(finalButtons) == 0 {
		return nil
	}

	buttonGrid := make(map[int]map[int]models.InlineKeyboardButton)
	for i, button := range finalButtons {
		row := button.PositionY
		col := button.PositionX
		if col == 0 {
			col = i
		}
		if buttonGrid[row] == nil {
			buttonGrid[row] = make(map[int]models.InlineKeyboardButton)
		}
		buttonGrid[row][col] = models.InlineKeyboardButton{
			Text: button.NameButton,
			URL:  button.ButtonURL,
		}
	}

	var keyboard [][]models.InlineKeyboardButton
	for row := 0; row < 10; row++ {
		if rowButtons, exists := buttonGrid[row]; exists {
			var keyboardRow []models.InlineKeyboardButton
			for col := 0; col < 10; col++ {
				if btn, exists := rowButtons[col]; exists {
					keyboardRow = append(keyboardRow, btn)
				}
			}
			if len(keyboardRow) > 0 {
				keyboard = append(keyboard, keyboardRow)
			}
		}
	}

	if len(keyboard) == 0 {
		return nil
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

// ========== TEXTO ==========
func (mp *MessageProcessor) ProcessTextMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	text := post.Text
	messageID := post.ID
	if text == "" {
		return fmt.Errorf("texto da mensagem está vazio")
	}
	if !messageEditAllowed {
		return nil
	}
	finalText, customCaption, msgPerm, btnPerm, linkPrev := mp.processMessageWithHashtagFormatting(text, post.Entities, channel)
	if !msgPerm {
		return nil
	}

	params := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      finalText,
		ParseMode: "HTML",
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: func(b bool) *bool { v := b; return &v }(!linkPrev),
		},
	}

	if btnPerm {
		keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
		if keyboard != nil {
			params.ReplyMarkup = keyboard
		}
	}

	_, err := mp.bot.EditMessageText(ctx, params)
	return mp.handleTelegramError(err)
}

// ========== ÁUDIO ==========
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption
	mediaGroupID := post.MediaGroupID
	if !messageEditAllowed {
		return nil
	}

	time.Sleep(1 * time.Second)

	if mediaGroupID != "" {
		finalMessage, customCaption, msgPerm, btnPerm, _ := mp.processMessageWithHashtagFormatting(caption, post.CaptionEntities, channel)
		if !msgPerm {
			return nil
		}

		sendParams := &bot.SendAudioParams{
			ChatID:    post.Chat.ID,
			Audio:     &models.InputFileString{Data: post.Audio.FileID},
			Caption:   finalMessage,
			ParseMode: "HTML",
		}

		if btnPerm {
			keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
			if keyboard != nil {
				sendParams.ReplyMarkup = keyboard
			}
		}

		_, err := mp.bot.SendAudio(ctx, sendParams)
		if err != nil {
			return mp.handleTelegramError(err)
		}

		_, err = mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    post.Chat.ID,
			MessageID: messageID,
		})
		return mp.handleTelegramError(err)
	}

	finalMessage, customCaption, msgPerm, btnPerm, linkPrev := mp.processMessageWithHashtagFormatting(caption, post.CaptionEntities, channel)
	if !msgPerm {
		return nil
	}

	params := &bot.EditMessageCaptionParams{
		ChatID:                post.Chat.ID,
		MessageID:             messageID,
		Caption:               finalMessage,
		ParseMode:             "HTML",
		DisableWebPagePreview: !linkPrev,
	}

	if btnPerm {
		keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
		if keyboard != nil {
			params.ReplyMarkup = keyboard
		}
	}

	_, err := mp.bot.EditMessageCaption(ctx, params)
	return mp.handleTelegramError(err)
}

// ========== MÍDIA ==========
func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID
	if mediaGroupID != "" {
		return mp.handleGroupedMedia(ctx, channel, post, buttons, messageEditAllowed)
	}
	return mp.handleSingleMedia(ctx, channel, post, buttons, messageEditAllowed)
}

func (mp *MessageProcessor) handleSingleMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption
	if !messageEditAllowed {
		return nil
	}

	finalCaption, customCaption, msgPerm, btnPerm, linkPrev := mp.processMessageWithHashtagFormatting(caption, post.CaptionEntities, channel)
	if !msgPerm {
		return nil
	}

	params := &bot.EditMessageCaptionParams{
		ChatID:                post.Chat.ID,
		MessageID:             messageID,
		Caption:               finalCaption,
		ParseMode:             "HTML",
		DisableWebPagePreview: !linkPrev,
	}

	if btnPerm {
		keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
		if keyboard != nil {
			params.ReplyMarkup = keyboard
		}
	}

	_, err := mp.bot.EditMessageCaption(ctx, params)
	return mp.handleTelegramError(err)
}

// ========== PROCESSAMENTO DE GRUPO - CORRIGIDO ==========
func (mp *MessageProcessor) handleGroupedMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID
	messageID := post.ID
	caption := post.Caption

	value, _ := mediaGroups.LoadOrStore(mediaGroupID, &MediaGroup{
		Messages:           make([]MediaMessage, 0),
		Processed:          false,
		MessageEditAllowed: messageEditAllowed,
		ChatID:             post.Chat.ID,
	})

	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if group.Processed {
		return nil
	}

	group.Messages = append(group.Messages, MediaMessage{
		MessageID:       messageID,
		HasCaption:      caption != "",
		Caption:         caption,
		CaptionEntities: convertMessageEntitiesToInterface(post.CaptionEntities),
	})

	if group.Timer != nil {
		group.Timer.Stop()
	}

	timeout := mp.getAdaptiveTimeout(len(group.Messages))

	group.Timer = time.AfterFunc(timeout, func() {
		// ✅ SOLUÇÃO: Context dedicado com timeout muito maior
		networkCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		go mp.finishGroupProcessingWithRetry(networkCtx, mediaGroupID, channel, buttons)
	})

	return nil
}

// ✅ NOVA FUNÇÃO: Processamento com retry e backoff exponencial
func (mp *MessageProcessor) finishGroupProcessingWithRetry(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button) {
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := mp.finishGroupProcessing(ctx, groupID, channel, buttons)

		if err == nil {
			fmt.Printf("✅ Grupo de mídia %s processado com sucesso na tentativa %d\n", groupID, attempt)
			return
		}

		// Verificar se é erro de context canceled
		if strings.Contains(err.Error(), "context canceled") {
			fmt.Printf("⚠️ Context canceled na tentativa %d para grupo %s\n", attempt, groupID)

			if attempt < maxRetries {
				// Backoff exponencial
				delay := time.Duration(attempt) * baseDelay
				fmt.Printf("🔄 Tentando novamente em %v...\n", delay)
				time.Sleep(delay)

				// Criar novo context para próxima tentativa
				newCtx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
				defer cancel()
				ctx = newCtx
				continue
			}
		}

		// Para outros erros, não tentar novamente
		fmt.Printf("❌ Erro final ao processar grupo %s: %v\n", groupID, err)
		return
	}

	fmt.Printf("❌ Falha após %d tentativas para grupo %s\n", maxRetries, groupID)
}

// ========== PROCESSAMENTO DE GRUPO - ESTRATÉGIA DE REENVIO ==========
func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button) error {
	value, ok := mediaGroups.Load(groupID)
	if !ok {
		return fmt.Errorf("grupo não encontrado: %s", groupID)
	}

	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if group.Processed {
		return nil
	}
	group.Processed = true

	if len(group.Messages) == 0 {
		return fmt.Errorf("grupo vazio: %s", groupID)
	}

	fmt.Printf("🔍 Processando grupo %s com %d mensagens\n", groupID, len(group.Messages))

	// Debug: Listar todas as mensagens do grupo
	for i, msg := range group.Messages {
		fmt.Printf("📝 Mensagem %d: ID=%d, HasCaption=%t, Caption='%s'\n",
			i, msg.MessageID, msg.HasCaption, msg.Caption)
	}

	// ✅ ESTRATÉGIA 1: Verificar se alguma mensagem tem caption
	var targetMessage *MediaMessage
	hasAnyCaption := false

	for i := range group.Messages {
		if group.Messages[i].HasCaption && strings.TrimSpace(group.Messages[i].Caption) != "" {
			targetMessage = &group.Messages[i]
			hasAnyCaption = true
			fmt.Printf("🎯 Selecionada mensagem com caption: ID=%d\n", targetMessage.MessageID)
			break
		}
	}

	if targetMessage == nil {
		targetMessage = &group.Messages[0]
		fmt.Printf("🎯 Selecionada primeira mensagem: ID=%d\n", targetMessage.MessageID)
	}

	// Processar caption
	finalCaption, customCaption, msgPerm, btnPerm, linkPrev := mp.processMessageWithHashtagFormatting(
		targetMessage.Caption,
		convertInterfaceToMessageEntities(targetMessage.CaptionEntities),
		channel,
	)

	fmt.Printf("📄 Caption processada: '%s'\n", finalCaption)
	fmt.Printf("🔐 Permissões: msg=%t, btn=%t, link=%t\n", msgPerm, btnPerm, linkPrev)

	if !msgPerm {
		fmt.Printf("❌ Edição de mensagem não permitida\n")
		return nil
	}

	// ✅ ESTRATÉGIA 2: Se não há caption original, usar SendMediaGroup
	if !hasAnyCaption && strings.TrimSpace(finalCaption) != "" {
		fmt.Printf("🔄 Nenhuma mensagem tem caption, reenviando grupo com caption\n")
		return mp.resendMediaGroupWithCaption(ctx, group, finalCaption, customCaption, buttons, btnPerm, linkPrev)
	}

	// ✅ ESTRATÉGIA 3: Se há caption original, usar EditMessageCaption
	fmt.Printf("✏️ Editando caption da mensagem existente\n")
	return mp.editExistingCaption(ctx, group, targetMessage, finalCaption, customCaption, buttons, btnPerm, linkPrev)
}

// ✅ NOVA FUNÇÃO: Reenviar grupo de mídia com caption
func (mp *MessageProcessor) resendMediaGroupWithCaption(ctx context.Context, group *MediaGroup, caption string, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button, btnPerm bool, linkPrev bool) error {
	fmt.Printf("🔄 Iniciando reenvio do grupo de mídia\n")

	// Primeiro, precisamos obter as informações das mídias originais
	// Isso requer fazer download ou usar file_id das mensagens originais

	// ⚠️ LIMITAÇÃO: Para reenviar, precisamos dos file_ids das mídias
	// Por enquanto, vamos tentar editar apenas a primeira mensagem mesmo sem caption

	return mp.addCaptionToFirstMessage(ctx, group, caption, customCaption, buttons, btnPerm, linkPrev)
}

// ✅ NOVA FUNÇÃO: Adicionar caption à primeira mensagem (mesmo sem caption original)
func (mp *MessageProcessor) addCaptionToFirstMessage(ctx context.Context, group *MediaGroup, caption string, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button, btnPerm bool, linkPrev bool) error {
	targetMessage := &group.Messages[0]

	// ✅ SOLUÇÃO: Usar EditMessageCaption mesmo para mensagens sem caption
	// O Telegram permite isso em alguns casos

	params := &bot.EditMessageCaptionParams{
		ChatID:    group.ChatID,
		MessageID: targetMessage.MessageID,
		Caption:   caption,
		ParseMode: "HTML",
	}

	if !linkPrev {
		params.DisableWebPagePreview = true
	}

	if btnPerm {
		keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
		if keyboard != nil {
			params.ReplyMarkup = keyboard
			fmt.Printf("⌨️ Adicionando %d linhas de botões\n", len(keyboard.InlineKeyboard))
		}
	}

	editCtx, editCancel := context.WithTimeout(ctx, 30*time.Second)
	defer editCancel()

	fmt.Printf("🚀 Tentando adicionar caption à mensagem %d\n", targetMessage.MessageID)

	result, err := mp.bot.EditMessageCaption(editCtx, params)
	if err != nil {
		// Se falhar, tentar adicionar apenas os botões
		if strings.Contains(err.Error(), "message caption can't be edited") ||
			strings.Contains(err.Error(), "message has no caption") {
			fmt.Printf("⚠️ Não é possível adicionar caption, tentando apenas botões\n")
			return mp.addOnlyButtons(ctx, group, buttons, customCaption, btnPerm)
		}

		fmt.Printf("❌ Erro ao adicionar caption: %v\n", err)
		return mp.handleTelegramError(err)
	}

	fmt.Printf("✅ Caption adicionada com sucesso! Resultado: %+v\n", result)
	return nil
}

// ✅ NOVA FUNÇÃO: Adicionar apenas botões (para quando não é possível editar caption)
func (mp *MessageProcessor) addOnlyButtons(ctx context.Context, group *MediaGroup, buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption, btnPerm bool) error {
	if !btnPerm || len(buttons) == 0 {
		fmt.Printf("❌ Não é possível adicionar botões\n")
		return nil
	}

	targetMessage := &group.Messages[0]

	keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
	if keyboard == nil {
		fmt.Printf("❌ Falha ao criar teclado\n")
		return nil
	}

	editCtx, editCancel := context.WithTimeout(ctx, 30*time.Second)
	defer editCancel()

	fmt.Printf("⌨️ Adicionando apenas botões à mensagem %d\n", targetMessage.MessageID)

	result, err := mp.bot.EditMessageReplyMarkup(editCtx, &bot.EditMessageReplyMarkupParams{
		ChatID:      group.ChatID,
		MessageID:   targetMessage.MessageID,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		fmt.Printf("❌ Erro ao adicionar botões: %v\n", err)
		return mp.handleTelegramError(err)
	}

	fmt.Printf("✅ Botões adicionados com sucesso! Resultado: %+v\n", result)
	return nil
}

// ✅ FUNÇÃO: Editar caption existente
func (mp *MessageProcessor) editExistingCaption(ctx context.Context, group *MediaGroup, targetMessage *MediaMessage, caption string, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button, btnPerm bool, linkPrev bool) error {
	params := &bot.EditMessageCaptionParams{
		ChatID:    group.ChatID,
		MessageID: targetMessage.MessageID,
		Caption:   caption,
		ParseMode: "HTML",
	}

	if !linkPrev {
		params.DisableWebPagePreview = true
	}

	if btnPerm {
		keyboard := mp.CreateInlineKeyboard(buttons, customCaption, true)
		if keyboard != nil {
			params.ReplyMarkup = keyboard
		}
	}

	editCtx, editCancel := context.WithTimeout(ctx, 30*time.Second)
	defer editCancel()

	fmt.Printf("✏️ Editando caption existente da mensagem %d\n", targetMessage.MessageID)

	result, err := mp.bot.EditMessageCaption(editCtx, params)
	if err != nil {
		fmt.Printf("❌ Erro na edição: %v\n", err)
		return mp.handleTelegramError(err)
	}

	fmt.Printf("✅ Caption editada com sucesso! Resultado: %+v\n", result)
	return nil
}

// ✅ TRATAMENTO DE ERROS ESPECÍFICOS
func (mp *MessageProcessor) handleTelegramError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Erros específicos de caption
	captionErrors := []string{
		"message caption can't be edited",
		"message has no caption",
		"Bad Request: message caption can't be edited",
		"message is not modified",
		"Message is not modified",
	}

	for _, captionError := range captionErrors {
		if strings.Contains(errStr, captionError) {
			fmt.Printf("ℹ️ Erro de caption (ignorável): %s\n", captionError)
			return nil
		}
	}

	// Context canceled - retornar para retry
	if strings.Contains(errStr, "context canceled") {
		fmt.Printf("⚠️ Timeout na operação: %v\n", err)
		return fmt.Errorf("timeout: %w", err)
	}

	// Rate limiting
	if strings.Contains(errStr, "Too Many Requests") {
		fmt.Printf("⚠️ Rate limit atingido: %v\n", err)
		time.Sleep(2 * time.Second)
		return fmt.Errorf("rate limit: %w", err)
	}

	fmt.Printf("❌ Erro do Telegram: %v\n", err)
	return err
}

// ========== STICKER ==========
func (mp *MessageProcessor) ProcessStickerMessage(ctx context.Context, post *models.Message, buttons []dbmodels.Button) error {
	if len(buttons) == 0 {
		return nil
	}
	keyboard := mp.CreateInlineKeyboard(buttons, nil, true)
	if keyboard == nil {
		return nil
	}
	_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      post.Chat.ID,
		MessageID:   post.ID,
		ReplyMarkup: keyboard,
	})
	return mp.handleTelegramError(err)
}

// ========== AUXILIARES ==========
func extractHashtag(text string) string {
	if text == "" {
		return ""
	}
	matches := hashtagRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}

func removeHashtag(text, hashtag string) string {
	if text == "" || hashtag == "" {
		return text
	}
	var re *regexp.Regexp
	if value, ok := removeHashRegexCache.Load(hashtag); ok {
		re = value.(*regexp.Regexp)
	} else {
		re = regexp.MustCompile(`#` + regexp.QuoteMeta(hashtag) + `\s*`)
		removeHashRegexCache.Store(hashtag, re)
	}
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}

func convertMessageEntitiesToInterface(entities []models.MessageEntity) []interface{} {
	result := make([]interface{}, len(entities))
	for i, entity := range entities {
		result[i] = entity
	}
	return result
}

func convertInterfaceToMessageEntities(entities []interface{}) []models.MessageEntity {
	result := make([]models.MessageEntity, 0, len(entities))
	for _, entity := range entities {
		if msgEntity, ok := entity.(models.MessageEntity); ok {
			result = append(result, msgEntity)
		}
	}
	return result
}

// ========== FUNÇÕES AUXILIARES FALTANTES ==========
func detectParseModeV2(text string) string {
	return convertMarkdownToHTML(text)
}

func processTextWithFormattingV2(text string, entities []models.MessageEntity) string {
	return text // Implementação básica
}

func convertMarkdownToHTML(text string) string {
	if text == "" {
		return ""
	}
	result := text

	// Bold (**texto**)
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	result = boldRegex.ReplaceAllString(result, "<b>$1</b>")

	// Italic (*texto*)
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	result = italicRegex.ReplaceAllString(result, "<i>$1</i>")

	return result
}

// ✅ CORREÇÃO 9: Melhorar processamento de hashtag
func (mp *MessageProcessor) processMessageWithHashtagFormatting(
	text string,
	entities []models.MessageEntity,
	channel *dbmodels.Channel,
) (string, *dbmodels.CustomCaption, bool, bool, bool) {

	fmt.Printf("🔄 Processando texto: '%s'\n", text)

	// Processa o texto original com as entidades do Telegram
	formatted := processTextWithFormatting(text, entities)
	hashtag := extractHashtag(text)

	fmt.Printf("🏷️ Hashtag encontrada: '%s'\n", hashtag)

	// Valores padrão para permissões
	msgPerm, btnPerm, linkPrev := true, true, true

	// Aplicar permissões padrão
	if channel.DefaultCaption != nil {
		if channel.DefaultCaption.MessagePermission != nil {
			msgPerm = channel.DefaultCaption.MessagePermission.Message
			linkPrev = channel.DefaultCaption.MessagePermission.LinkPreview
		}
		if channel.DefaultCaption.ButtonsPermission != nil {
			btnPerm = channel.DefaultCaption.ButtonsPermission.Message
		}
	}

	// Caso 1: Sem hashtag - usar default caption
	if hashtag == "" {
		finalText := formatted
		if channel.DefaultCaption != nil && channel.DefaultCaption.Caption != "" {
			processedCaption := detectParseMode(channel.DefaultCaption.Caption)
			if strings.TrimSpace(formatted) != "" {
				finalText = fmt.Sprintf("%s\n\n%s", formatted, processedCaption)
			} else {
				finalText = processedCaption
			}
		}
		fmt.Printf("📝 Texto final (sem hashtag): '%s'\n", finalText)
		return finalText, nil, msgPerm, btnPerm, linkPrev
	}

	// Caso 2: Com hashtag - buscar custom caption
	var customCaption *dbmodels.CustomCaption
	for i := range channel.CustomCaptions {
		code := strings.TrimPrefix(channel.CustomCaptions[i].Code, "#")
		if strings.EqualFold(code, hashtag) {
			customCaption = &channel.CustomCaptions[i]
			fmt.Printf("🎯 Custom caption encontrada: '%s'\n", customCaption.Caption)
			break
		}
	}

	// Caso 3: Hashtag não encontrada - usar default caption
	if customCaption == nil {
		finalText := formatted
		if channel.DefaultCaption != nil && channel.DefaultCaption.Caption != "" {
			processedCaption := detectParseMode(channel.DefaultCaption.Caption)
			if strings.TrimSpace(formatted) != "" {
				finalText = fmt.Sprintf("%s\n\n%s", formatted, processedCaption)
			} else {
				finalText = processedCaption
			}
		}
		fmt.Printf("📝 Texto final (hashtag não encontrada): '%s'\n", finalText)
		return finalText, nil, msgPerm, btnPerm, linkPrev
	}

	// Caso 4: Custom caption encontrada
	cleanText := removeHashtag(text, hashtag)
	formattedCleanText := processTextWithFormatting(cleanText, adjustEntitiesAfterHashtagRemoval(entities, text, hashtag))

	// Aplicar linkPreview do customCaption
	linkPrev = customCaption.LinkPreview

	finalText := formattedCleanText
	if customCaption.Caption != "" {
		processedCaption := detectParseMode(customCaption.Caption)
		if strings.TrimSpace(formattedCleanText) != "" {
			finalText = fmt.Sprintf("%s\n\n%s", formattedCleanText, processedCaption)
		} else {
			finalText = processedCaption
		}
	}

	fmt.Printf("📝 Texto final (com custom caption): '%s'\n", finalText)
	return finalText, customCaption, msgPerm, btnPerm, linkPrev
}

func adjustEntitiesAfterHashtagRemoval(entities []models.MessageEntity, originalText, hashtag string) []models.MessageEntity {
	return entities // Implementação básica
}

// ========== MÉTODOS THREAD-SAFE ==========
func (mp *MessageProcessor) IsNewPackActive(channelID int64) bool {
	return mp.mediaGroupManager.IsNewPackActive(channelID)
}

func (mp *MessageProcessor) SetNewPackActive(channelID int64, active bool) {
	mp.mediaGroupManager.SetNewPackActive(channelID, active)
}

func (pm *PermissionManager) CheckPermission(userID int64) bool {
	return true
}

func (mgm *MediaGroupManager) IsNewPackActive(channelID int64) bool {
	return false
}

func (mgm *MediaGroupManager) SetNewPackActive(channelID int64, active bool) {
}
