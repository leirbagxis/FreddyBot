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

func NewMessageProcessor(b *bot.Bot) *MessageProcessor {
	return &MessageProcessor{
		bot:               b,
		permissionManager: NewPermissionManager(),
		mediaGroupManager: NewMediaGroupManager(),
	}
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
	return err
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
			return err
		}
		_, err = mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    post.Chat.ID,
			MessageID: messageID,
		})
		return err
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
	return err
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
	return err
}

// ========== GRUPO DE MÍDIA ==========
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
	timeout := time.Duration(800+len(group.Messages)*200) * time.Millisecond
	if timeout > 2*time.Second {
		timeout = 2 * time.Second
	}
	group.Timer = time.AfterFunc(timeout, func() {
		mp.finishGroupProcessing(ctx, mediaGroupID, channel, buttons)
	})
	return nil
}

func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button) {
	value, ok := mediaGroups.Load(groupID)
	if !ok {
		return
	}
	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()
	if group.Processed {
		return
	}
	group.Processed = true
	if len(group.Messages) == 0 {
		return
	}
	var targetMessage *MediaMessage
	for i := range group.Messages {
		if group.Messages[i].HasCaption {
			targetMessage = &group.Messages[i]
			break
		}
	}
	if targetMessage == nil {
		targetMessage = &group.Messages[0]
	}
	finalCaption, customCaption, msgPerm, btnPerm, linkPrev := mp.processMessageWithHashtagFormatting(
		targetMessage.Caption,
		convertInterfaceToMessageEntities(targetMessage.CaptionEntities),
		channel,
	)
	if !msgPerm {
		return
	}

	params := &bot.EditMessageCaptionParams{
		ChatID:                group.ChatID,
		MessageID:             targetMessage.MessageID,
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

	_, _ = mp.bot.EditMessageCaption(ctx, params)
	time.AfterFunc(10*time.Second, func() {
		mediaGroups.Delete(groupID)
	})
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
	return err
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

func findCustomCaption(channel *dbmodels.Channel, hashtag string) *dbmodels.CustomCaption {
	cacheKey := fmt.Sprintf("%d_%s", channel.ID, hashtag)
	if value, ok := customCaptionCache.Load(cacheKey); ok {
		if caption, ok := value.(*dbmodels.CustomCaption); ok {
			return caption
		}
		return nil
	}
	for i := range channel.CustomCaptions {
		ccCode := strings.TrimPrefix(channel.CustomCaptions[i].Code, "#")
		if strings.EqualFold(ccCode, hashtag) {
			customCaptionCache.Store(cacheKey, &channel.CustomCaptions[i])
			return &channel.CustomCaptions[i]
		}
	}
	customCaptionCache.Store(cacheKey, (*dbmodels.CustomCaption)(nil))
	return nil
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

// ========== CONVERSÃO MARKDOWN PARA HTML ==========
func convertMarkdownToHTML(text string) string {
	if text == "" {
		return ""
	}

	result := text

	// Processar em ordem específica para evitar conflitos

	// 1. Code blocks primeiro (```código```)
	codeBlockRegex := regexp.MustCompile("```([\\s\\S]*?)```")
	result = codeBlockRegex.ReplaceAllString(result, "<pre>$1</pre>")

	// 2. Inline code (`código`)
	inlineCodeRegex := regexp.MustCompile("`([^`]+)`")
	result = inlineCodeRegex.ReplaceAllString(result, "<code>$1</code>")

	// 3. Bold (**texto**)
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	result = boldRegex.ReplaceAllString(result, "<b>$1</b>")

	// 4. Italic (*texto*)
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	result = italicRegex.ReplaceAllString(result, "<i>$1</i>")

	// 5. Underline (__texto__)
	underlineRegex := regexp.MustCompile(`__([^_]+)__`)
	result = underlineRegex.ReplaceAllString(result, "<u>$1</u>")

	// 6. Strikethrough (~~texto~~)
	strikeRegex := regexp.MustCompile(`~~([^~]+)~~`)
	result = strikeRegex.ReplaceAllString(result, "<s>$1</s>")

	// 7. Spoiler (||texto||)
	spoilerRegex := regexp.MustCompile(`\|\|([^|]+)\|\|`)
	result = spoilerRegex.ReplaceAllString(result, `<span class="tg-spoiler">$1</span>`)

	// 8. Links [texto](url) - fazer por último
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	result = linkRegex.ReplaceAllString(result, `<a href="$2">$1</a>`)

	return result
}

// ========== PROCESSAMENTO DE MENSAGEM COM HASHTAG ==========
func (mp *MessageProcessor) processMessageWithHashtagFormatting(
	text string,
	entities []models.MessageEntity,
	channel *dbmodels.Channel,
) (string, *dbmodels.CustomCaption, bool, bool, bool) {

	// Processa o texto original com as entidades do Telegram
	formatted := processTextWithFormatting(text, entities)

	hashtag := extractHashtag(text)

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
			processedCaption := convertMarkdownToHTML(channel.DefaultCaption.Caption)
			finalText = fmt.Sprintf("%s\n\n%s", formatted, processedCaption)
		}
		return finalText, nil, msgPerm, btnPerm, linkPrev
	}

	// Caso 2: Com hashtag - buscar custom caption
	var customCaption *dbmodels.CustomCaption
	for i := range channel.CustomCaptions {
		code := strings.TrimPrefix(channel.CustomCaptions[i].Code, "#")
		if strings.EqualFold(code, hashtag) {
			customCaption = &channel.CustomCaptions[i]
			break
		}
	}

	// Caso 3: Hashtag não encontrada - usar default caption
	if customCaption == nil {
		finalText := formatted
		if channel.DefaultCaption != nil && channel.DefaultCaption.Caption != "" {
			processedCaption := convertMarkdownToHTML(channel.DefaultCaption.Caption)
			finalText = fmt.Sprintf("%s\n\n%s", formatted, processedCaption)
		}
		return finalText, nil, msgPerm, btnPerm, linkPrev
	}

	// Caso 4: Custom caption encontrada
	// Remove hashtag do texto original
	cleanText := removeHashtag(text, hashtag)
	formattedCleanText := processTextWithFormatting(cleanText, adjustEntitiesAfterHashtagRemoval(entities, text, hashtag))

	// Aplicar linkPreview do customCaption
	linkPrev = customCaption.LinkPreview

	finalText := formattedCleanText
	if customCaption.Caption != "" {
		processedCaption := convertMarkdownToHTML(customCaption.Caption)
		finalText = fmt.Sprintf("%s\n\n%s", formattedCleanText, processedCaption)
	}

	return finalText, customCaption, msgPerm, btnPerm, linkPrev
}

// ========== FUNÇÕES AUXILIARES ==========
func processMarkdownText(text string) string {
	if text == "" {
		return ""
	}
	return detectParseMode(text)
}

func processCustomCaptionText(caption string) string {
	if caption == "" {
		return ""
	}
	return detectParseMode(caption)
}

// Função auxiliar para ajustar as entidades após remoção de hashtag
func adjustEntitiesAfterHashtagRemoval(entities []models.MessageEntity, originalText, hashtag string) []models.MessageEntity {
	hashtagPattern := "#" + hashtag
	hashtagIndex := strings.Index(strings.ToLower(originalText), strings.ToLower(hashtagPattern))

	if hashtagIndex == -1 {
		return entities
	}

	// Calcular o deslocamento após remoção da hashtag
	hashtagLength := len(hashtagPattern)
	endIndex := hashtagIndex + hashtagLength
	for endIndex < len(originalText) && (originalText[endIndex] == ' ' || originalText[endIndex] == '\n') {
		endIndex++
	}

	removedLength := endIndex - hashtagIndex

	var adjustedEntities []models.MessageEntity
	for _, entity := range entities {
		// Se a entidade está completamente antes da hashtag, manter como está
		if entity.Offset+entity.Length <= hashtagIndex {
			adjustedEntities = append(adjustedEntities, entity)
			continue
		}

		// Se a entidade está completamente depois da hashtag, ajustar offset
		if entity.Offset >= endIndex {
			newEntity := entity
			newEntity.Offset -= removedLength
			adjustedEntities = append(adjustedEntities, newEntity)
			continue
		}

		// Se a entidade se sobrepõe com a hashtag, ajustar ou pular
		if entity.Offset < hashtagIndex && entity.Offset+entity.Length > endIndex {
			// Entidade atravessa a hashtag - dividir ou ajustar
			newEntity := entity
			newEntity.Length -= removedLength
			if newEntity.Length > 0 {
				adjustedEntities = append(adjustedEntities, newEntity)
			}
		} else if entity.Offset < endIndex && entity.Offset+entity.Length > hashtagIndex {
			// Entidade se sobrepõe parcialmente - pode ser necessário ajustar
			// Implementar lógica específica conforme necessário
		}
	}

	return adjustedEntities
}

// ========== MÉTODOS THREAD-SAFE ==========
func (mp *MessageProcessor) IsNewPackActive(channelID int64) bool {
	return mp.mediaGroupManager.IsNewPackActive(channelID)
}

func (mp *MessageProcessor) SetNewPackActive(channelID int64, active bool) {
	mp.mediaGroupManager.SetNewPackActive(channelID, active)
}
