package channelpost

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// ✅ REGEX E CACHES GLOBAIS
var (
	hashtagRegex         = regexp.MustCompile(`#(\w+)`)
	removeHashRegexCache = sync.Map{} // string -> *regexp.Regexp
	customCaptionCache   = sync.Map{} // string -> *dbmodels.CustomCaption
	mediaGroups          = sync.Map{} // string -> *MediaGroup
)

// ✅ ESTRUTURA ÚNICA PARA GRUPOS DE MÍDIA
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

func (mp *MessageProcessor) GetMessageType(post *models.Message) MessageType {
	if post.Text != "" {
		return MessageTypeText
	}
	if post.Audio != nil {
		return MessageTypeAudio
	}
	if post.Sticker != nil {
		return MessageTypeSticker
	}
	if post.Photo != nil {
		return MessageTypePhoto
	}
	if post.Video != nil {
		return MessageTypeVideo
	}
	if post.Animation != nil {
		return MessageTypeAnimation
	}
	return ""
}

// ✅ CRIAR KEYBOARD SIMPLES
func (mp *MessageProcessor) CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button

	// Usar custom caption buttons se existirem
	if customCaption != nil && len(customCaption.Buttons) > 0 {
		finalButtons = make([]dbmodels.Button, 0, len(customCaption.Buttons))
		for _, cb := range customCaption.Buttons {
			finalButtons = append(finalButtons, dbmodels.Button{
				NameButton: cb.NameButton,
				ButtonURL:  cb.ButtonURL,
				PositionY:  cb.PositionY,
				PositionX:  cb.PositionX,
			})
		}
	} else {
		finalButtons = buttons
	}

	if len(finalButtons) == 0 {
		return nil
	}

	// Criar grid de botões
	buttonGrid := make(map[int]map[int]models.InlineKeyboardButton)

	for i, button := range finalButtons {
		if button.NameButton == "" || button.ButtonURL == "" {
			continue
		}

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

	// Construir keyboard final
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

// ✅ PROCESSAR TEXTO COM FORMATAÇÃO
func (mp *MessageProcessor) ProcessTextMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	text := post.Text
	messageID := post.ID

	if text == "" {
		return fmt.Errorf("texto da mensagem está vazio")
	}

	if !messageEditAllowed {
		if len(buttons) == 0 {
			return nil
		}
		keyboard := mp.CreateInlineKeyboard(buttons, nil)
		if keyboard == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   messageID,
			ReplyMarkup: keyboard,
		})
		return err
	}

	// ✅ APLICAR FORMATAÇÃO
	formattedText := processTextWithFormatting(text, post.Entities)
	message, customCaption := mp.processMessageWithHashtag(formattedText, channel)
	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      message,
		ParseMode: "HTML",
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageText(ctx, editParams)
	return err
}

// ✅ PROCESSAR ÁUDIO (SUBSTITUIÇÃO TOTAL)
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption
	mediaGroupID := post.MediaGroupID

	if !messageEditAllowed {
		if len(buttons) == 0 {
			return nil
		}
		keyboard := mp.CreateInlineKeyboard(buttons, nil)
		if keyboard == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   messageID,
			ReplyMarkup: keyboard,
		})
		return err
	}

	// Aguardar 1 segundo
	time.Sleep(1 * time.Second)

	// Para grupos de mídia: REENVIAR + DELETAR
	if mediaGroupID != "" {
		var finalMessage string
		var customCaption *dbmodels.CustomCaption

		hashtag := extractHashtag(caption)
		if hashtag != "" {
			customCaption = findCustomCaption(channel, hashtag)
			if customCaption != nil {
				finalMessage = customCaption.Caption
			} else {
				if channel.DefaultCaption != nil {
					finalMessage = channel.DefaultCaption.Caption
				}
			}
		} else {
			if channel.DefaultCaption != nil {
				finalMessage = channel.DefaultCaption.Caption
			}
		}

		keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

		// Reenviar áudio
		sendParams := &bot.SendAudioParams{
			ChatID:    post.Chat.ID,
			Audio:     &models.InputFileString{Data: post.Audio.FileID},
			Caption:   finalMessage,
			ParseMode: "HTML",
		}

		if keyboard != nil {
			sendParams.ReplyMarkup = keyboard
		}

		_, err := mp.bot.SendAudio(ctx, sendParams)
		if err != nil {
			return err
		}

		// Deletar original
		_, err = mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    post.Chat.ID,
			MessageID: messageID,
		})
		return err
	}

	// Para áudios individuais: SUBSTITUIÇÃO TOTAL
	var finalMessage string
	var customCaption *dbmodels.CustomCaption

	hashtag := extractHashtag(caption)
	if hashtag != "" {
		customCaption = findCustomCaption(channel, hashtag)
		if customCaption != nil {
			finalMessage = customCaption.Caption
		} else {
			if channel.DefaultCaption != nil {
				finalMessage = channel.DefaultCaption.Caption
			}
		}
	} else {
		if channel.DefaultCaption != nil {
			finalMessage = channel.DefaultCaption.Caption
		}
	}

	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Caption:   finalMessage,
		ParseMode: "HTML",
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageCaption(ctx, editParams)
	return err
}

// ✅ PROCESSAR MÍDIA
func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID

	if mediaGroupID != "" {
		return mp.handleGroupedMedia(ctx, channel, post, buttons, messageEditAllowed)
	}

	return mp.handleSingleMedia(ctx, channel, post, buttons, messageEditAllowed)
}

// ✅ MÍDIA INDIVIDUAL COM FORMATAÇÃO
func (mp *MessageProcessor) handleSingleMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption

	if !messageEditAllowed {
		if len(buttons) == 0 {
			return nil
		}
		keyboard := mp.CreateInlineKeyboard(buttons, nil)
		if keyboard == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   messageID,
			ReplyMarkup: keyboard,
		})
		return err
	}

	// ✅ APLICAR FORMATAÇÃO NA CAPTION
	formattedCaption := processTextWithFormatting(caption, post.CaptionEntities)
	message, customCaption := mp.processMessageWithHashtag(formattedCaption, channel)
	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Caption:   message,
		ParseMode: "HTML",
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageCaption(ctx, editParams)
	return err
}

// ✅ CORRIGIDO: Grupo de mídia com estrutura unificada
func (mp *MessageProcessor) handleGroupedMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID
	messageID := post.ID
	caption := post.Caption

	log.Printf("📸 Processando mídia do grupo: %s, ID: %d, Caption: %q", mediaGroupID, messageID, caption)

	// ✅ USAR LoadOrStore ATÔMICO
	value, loaded := mediaGroups.LoadOrStore(mediaGroupID, &MediaGroup{
		Messages:           make([]MediaMessage, 0),
		Processed:          false,
		MessageEditAllowed: messageEditAllowed,
		ChatID:             post.Chat.ID,
	})

	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if !loaded {
		log.Printf("📸 Novo grupo criado: %s", mediaGroupID)
	} else {
		log.Printf("📸 Usando grupo existente: %s", mediaGroupID)
	}

	// ✅ VERIFICAR SE JÁ FOI PROCESSADO
	if group.Processed {
		log.Printf("📸 Grupo já processado: %s", mediaGroupID)
		return nil
	}

	// ✅ ADICIONAR MENSAGEM
	group.Messages = append(group.Messages, MediaMessage{
		MessageID:       messageID,
		HasCaption:      caption != "",
		Caption:         caption,
		CaptionEntities: convertMessageEntitiesToInterface(post.CaptionEntities),
	})

	// ✅ CANCELAR TIMER ANTERIOR
	if group.Timer != nil {
		group.Timer.Stop()
	}

	// ✅ TIMEOUT ADAPTATIVO (reduzido para ser mais responsivo)
	timeout := time.Duration(800+len(group.Messages)*200) * time.Millisecond
	if timeout > 2*time.Second {
		timeout = 2 * time.Second
	}

	log.Printf("📸 Grupo %s: %d mensagens, timeout: %v", mediaGroupID, len(group.Messages), timeout)

	// ✅ CRIAR TIMER
	group.Timer = time.AfterFunc(timeout, func() {
		mp.finishGroupProcessing(ctx, mediaGroupID, channel, buttons)
	})

	return nil
}

// ✅ CORRIGIDO: Finalizar processamento de grupo - EDITA APENAS UMA MENSAGEM
func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button) {
	log.Printf("📸 Iniciando processamento final do grupo: %s", groupID)

	value, ok := mediaGroups.Load(groupID)
	if !ok {
		log.Printf("❌ Grupo não encontrado: %s", groupID)
		return
	}

	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if group.Processed {
		log.Printf("📸 Grupo já processado: %s", groupID)
		return
	}

	group.Processed = true
	log.Printf("📸 Marcando grupo como processado: %s com %d mensagens", groupID, len(group.Messages))

	if len(group.Messages) == 0 {
		log.Printf("❌ Nenhuma mensagem no grupo: %s", groupID)
		return
	}

	// ✅ ENCONTRAR A MENSAGEM IDEAL PARA EDITAR
	var targetMessage *MediaMessage

	// Prioridade 1: Mensagem com caption
	for i := range group.Messages {
		if group.Messages[i].HasCaption {
			targetMessage = &group.Messages[i]
			log.Printf("📸 Usando mensagem com caption: %d (caption: %q)", targetMessage.MessageID, targetMessage.Caption)
			break
		}
	}

	// Prioridade 2: Primeira mensagem se não houver caption
	if targetMessage == nil {
		targetMessage = &group.Messages[0]
		log.Printf("📸 Usando primeira mensagem (sem caption): %d", targetMessage.MessageID)
	}

	// ✅ SE NÃO PODE EDITAR MENSAGEM, APENAS ADICIONAR BOTÕES
	if !group.MessageEditAllowed {
		if len(buttons) > 0 {
			keyboard := mp.CreateInlineKeyboard(buttons, nil)
			if keyboard != nil {
				editCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := mp.bot.EditMessageReplyMarkup(editCtx, &bot.EditMessageReplyMarkupParams{
					ChatID:      group.ChatID,
					MessageID:   targetMessage.MessageID,
					ReplyMarkup: keyboard,
				})
				if err != nil {
					log.Printf("❌ Erro ao editar markup do grupo %s: %v", groupID, err)
				} else {
					log.Printf("✅ Markup editado para grupo: %s, mensagem: %d", groupID, targetMessage.MessageID)
				}
			}
		}
		mp.cleanupGroup(groupID)
		return
	}

	// ✅ PROCESSAR CAPTION E EDITAR APENAS UMA MENSAGEM
	var finalMessage string
	var customCaption *dbmodels.CustomCaption

	if targetMessage.HasCaption {
		// Aplicar formatação se tiver entities
		entities := convertInterfaceToMessageEntities(targetMessage.CaptionEntities)
		formattedCaption := processTextWithFormatting(targetMessage.Caption, entities)
		finalMessage, customCaption = mp.processMessageWithHashtag(formattedCaption, channel)
		log.Printf("📸 Processando com caption formatado: %s -> %s", targetMessage.Caption, finalMessage)
	} else {
		// Usar caption padrão se não houver caption na mensagem
		if channel.DefaultCaption != nil {
			finalMessage = channel.DefaultCaption.Caption
		}
		log.Printf("📸 Usando caption padrão: %s", finalMessage)
	}

	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	// ✅ EDITAR APENAS A MENSAGEM ALVO
	editCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	editParams := &bot.EditMessageCaptionParams{
		ChatID:    group.ChatID,
		MessageID: targetMessage.MessageID,
		Caption:   finalMessage,
		ParseMode: "HTML",
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageCaption(editCtx, editParams)
	if err != nil {
		log.Printf("❌ Erro ao editar caption do grupo %s, mensagem %d: %v", groupID, targetMessage.MessageID, err)
	} else {
		log.Printf("✅ SUCESSO: Grupo %s processado - APENAS mensagem %d editada com caption: %q", groupID, targetMessage.MessageID, finalMessage)
	}

	// ✅ CLEANUP
	mp.cleanupGroup(groupID)
}

// ✅ FUNÇÃO PARA LIMPEZA DO GRUPO
func (mp *MessageProcessor) cleanupGroup(groupID string) {
	time.AfterFunc(10*time.Second, func() {
		mediaGroups.Delete(groupID)
		log.Printf("🧹 Grupo removido da memória: %s", groupID)
	})
}

func (mp *MessageProcessor) ProcessStickerMessage(ctx context.Context, post *models.Message, buttons []dbmodels.Button) error {
	if len(buttons) == 0 {
		return nil
	}

	keyboard := mp.CreateInlineKeyboard(buttons, nil)
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

// ✅ FUNÇÕES AUXILIARES
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
		if strings.EqualFold(channel.CustomCaptions[i].Code, hashtag) {
			customCaptionCache.Store(cacheKey, &channel.CustomCaptions[i])
			return &channel.CustomCaptions[i]
		}
	}

	customCaptionCache.Store(cacheKey, (*dbmodels.CustomCaption)(nil))
	return nil
}

// ✅ PROCESSAR HASHTAG (CONCATENAÇÃO)
func (mp *MessageProcessor) processMessageWithHashtag(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
	hashtag := extractHashtag(text)
	var customCaption *dbmodels.CustomCaption

	if hashtag != "" {
		customCaption = findCustomCaption(channel, hashtag)
		cleanText := removeHashtag(text, hashtag)

		if customCaption != nil {
			return fmt.Sprintf("%s\n\n%s", cleanText, customCaption.Caption), customCaption
		}

		defaultCaption := ""
		if channel.DefaultCaption != nil {
			defaultCaption = channel.DefaultCaption.Caption
		}
		return fmt.Sprintf("%s\n\n%s", cleanText, defaultCaption), nil
	}

	defaultCaption := ""
	if channel.DefaultCaption != nil {
		defaultCaption = channel.DefaultCaption.Caption
	}

	return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
}

// ✅ FUNÇÕES DE CONVERSÃO
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

// ✅ MÉTODOS BÁSICOS
func (mp *MessageProcessor) IsNewPackActive(channelID int64) bool {
	return mp.mediaGroupManager.IsNewPackActive(channelID)
}

func (mp *MessageProcessor) SetNewPackActive(channelID int64, active bool) {
	mp.mediaGroupManager.SetNewPackActive(channelID, active)
}
