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
	removeHashRegexCache = sync.Map{}
	customCaptionCache   = sync.Map{}
)

type MessageProcessor struct {
	bot               *bot.Bot
	permissionManager *PermissionManager
	mediaGroupManager *MediaGroupManager
}

var groupCreationMutex sync.Mutex

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

func (mp *MessageProcessor) CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button

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

// ✅ ÁUDIO: SUBSTITUIÇÃO TOTAL + REENVIO PARA GRUPOS
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption
	mediaGroupID := post.MediaGroupID

	log.Printf("🎵 Processando áudio - ID: %d, Grupo: %s, Caption: %s", messageID, mediaGroupID, caption)

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

	// Aguardar 1 segundo (igual JS)
	time.Sleep(1 * time.Second)

	// Gerar nova legenda (SUBSTITUIÇÃO TOTAL)
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

	// Para grupos de mídia: REENVIAR + DELETAR
	if mediaGroupID != "" {
		log.Printf("🎵 Reenviando áudio do grupo: %s", mediaGroupID)

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
			log.Printf("❌ Erro ao reenviar áudio: %v", err)
			return err
		}

		// Deletar original
		_, err = mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    post.Chat.ID,
			MessageID: messageID,
		})
		if err != nil {
			log.Printf("❌ Erro ao deletar áudio original: %v", err)
		}

		log.Printf("✅ Áudio reenviado e original deletado")
		return err
	}

	// Para áudios individuais: EDITAR CAPTION
	log.Printf("🎵 Editando caption do áudio individual")

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
	if err != nil {
		log.Printf("❌ Erro ao editar caption do áudio: %v", err)
	} else {
		log.Printf("✅ Caption do áudio editado com sucesso")
	}

	return err
}

// ✅ MÍDIA (FOTOS/VÍDEOS): Verificar se é grupo ou individual
func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID

	log.Printf("📸 Processando mídia - ID: %d, Grupo: %s", post.ID, mediaGroupID)

	if mediaGroupID != "" {
		return mp.handleGroupedMedia(ctx, channel, post, buttons, messageEditAllowed)
	}

	return mp.handleSingleMedia(ctx, channel, post, buttons, messageEditAllowed)
}

// ✅ CORRIGIDO: Texto com formatação preservada
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

	// ✅ APLICAR FORMATAÇÃO COMPLETA
	formattedText := processTextWithFormatting(text, post.Entities)
	message, customCaption := mp.processMessageWithHashtagPreservingFormat(formattedText, channel)
	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      message,
		ParseMode: "HTML", // ✅ IMPORTANTE: HTML para formatação
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageText(ctx, editParams)
	return err
}

// ✅ CORRIGIDO: Mídia individual com formatação preservada
func (mp *MessageProcessor) handleSingleMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption

	log.Printf("📸 Processando mídia individual - ID: %d", messageID)

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

	// ✅ APLICAR FORMATAÇÃO COMPLETA PARA CAPTION
	formattedCaption := processTextWithFormatting(caption, post.CaptionEntities)
	message, customCaption := mp.processMessageWithHashtagPreservingFormat(formattedCaption, channel)
	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Caption:   message,
		ParseMode: "HTML", // ✅ IMPORTANTE: HTML para formatação
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageCaption(ctx, editParams)
	if err != nil {
		log.Printf("❌ Erro ao editar mídia individual: %v", err)
	} else {
		log.Printf("✅ Mídia individual processada com formatação")
	}

	return err
}

// ✅ NOVA FUNÇÃO: Processar preservando formatação (SEM entities, já aplicadas)
func (mp *MessageProcessor) processMessageWithHashtagPreservingFormat(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
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

	if text == "" {
		return defaultCaption, nil
	}

	return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
}

// ✅ MUTEX MAIS ESPECÍFICO POR GRUPO
var groupMutexes = sync.Map{} // string -> *sync.Mutex

// ✅ FUNÇÃO PARA OBTER MUTEX ESPECÍFICO DO GRUPO
func getGroupMutex(groupID string) *sync.Mutex {
	value, _ := groupMutexes.LoadOrStore(groupID, &sync.Mutex{})
	return value.(*sync.Mutex)
}

// ✅ CORRIGIDO: Mídia em grupo com mutex específico por grupo
func (mp *MessageProcessor) handleGroupedMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	mediaGroupID := post.MediaGroupID
	messageID := post.ID
	caption := post.Caption

	log.Printf("📸 Processando mídia do grupo: %s, ID: %d", mediaGroupID, messageID)

	// ✅ USAR MUTEX ESPECÍFICO PARA ESTE GRUPO
	groupMutex := getGroupMutex(mediaGroupID)
	groupMutex.Lock()
	defer groupMutex.Unlock()

	// ✅ VERIFICAR SE GRUPO JÁ FOI PROCESSADO
	if value, ok := mp.mediaGroupManager.groups.Load(mediaGroupID); ok {
		groupInfo := value.(*MediaGroupInfo)
		if groupInfo.Processed {
			log.Printf("📸 Grupo já processado, ignorando: %s", mediaGroupID)
			return nil
		}
	}

	// ✅ CRIAR OU OBTER GRUPO (agora thread-safe com mutex específico)
	value, loaded := mp.mediaGroupManager.groups.LoadOrStore(mediaGroupID, &MediaGroupInfo{
		Messages:           make([]MediaMessage, 0),
		Processed:          false,
		MessageEditAllowed: messageEditAllowed,
	})

	groupInfo := value.(*MediaGroupInfo)

	if !loaded {
		log.Printf("📸 Novo grupo criado: %s", mediaGroupID)
	} else {
		log.Printf("📸 Usando grupo existente: %s", mediaGroupID)
	}

	// ✅ VERIFICAR NOVAMENTE SE FOI PROCESSADO (dentro do mutex)
	if groupInfo.Processed {
		log.Printf("📸 Grupo já processado (double-check): %s", mediaGroupID)
		return nil
	}

	// ✅ ADICIONAR MENSAGEM
	groupInfo.Messages = append(groupInfo.Messages, MediaMessage{
		MessageID:       messageID,
		HasCaption:      caption != "",
		Caption:         caption,
		CaptionEntities: convertMessageEntitiesToInterface(post.CaptionEntities),
	})

	// ✅ CANCELAR TIMER ANTERIOR
	if groupInfo.Timer != nil {
		groupInfo.Timer.Stop()
		log.Printf("📸 Timer anterior cancelado para grupo: %s", mediaGroupID)
	}

	// ✅ TIMEOUT BASEADO NO TAMANHO REAL DO GRUPO
	timeout := time.Duration(2000+len(groupInfo.Messages)*500) * time.Millisecond
	if timeout > 5*time.Second {
		timeout = 5 * time.Second
	}

	log.Printf("📸 Grupo %s: %d mensagens, timeout: %v", mediaGroupID, len(groupInfo.Messages), timeout)

	// ✅ CRIAR TIMER (apenas um por grupo)
	groupInfo.Timer = time.AfterFunc(timeout, func() {
		log.Printf("📸 Timer disparado para grupo: %s", mediaGroupID)
		processCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		mp.finishGroupProcessing(processCtx, mediaGroupID, channel, buttons, post.Chat.ID)
	})

	return nil
}

// ✅ CORRIGIDO: Processar grupo com verificação mais rigorosa
func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button, chatID int64) {
	log.Printf("📸 Iniciando finishGroupProcessing para: %s", groupID)

	// ✅ USAR O MESMO MUTEX DO GRUPO
	groupMutex := getGroupMutex(groupID)
	groupMutex.Lock()
	defer groupMutex.Unlock()

	value, ok := mp.mediaGroupManager.groups.Load(groupID)
	if !ok {
		log.Printf("❌ Grupo não encontrado: %s", groupID)
		return
	}

	groupInfo := value.(*MediaGroupInfo)

	// ✅ VERIFICAR SE JÁ FOI PROCESSADO
	if groupInfo.Processed {
		log.Printf("📸 Grupo já processado: %s", groupID)
		return
	}

	// ✅ MARCAR COMO PROCESSADO IMEDIATAMENTE
	groupInfo.Processed = true
	log.Printf("📸 Marcando grupo como processado: %s", groupID)

	log.Printf("📸 Finalizando processamento do grupo: %s com %d mensagens", groupID, len(groupInfo.Messages))

	// ✅ ENCONTRAR MENSAGEM ALVO
	var targetMessage *MediaMessage
	for i := range groupInfo.Messages {
		if groupInfo.Messages[i].HasCaption {
			targetMessage = &groupInfo.Messages[i]
			log.Printf("📸 Usando mensagem com caption: %d", targetMessage.MessageID)
			break
		}
	}

	if targetMessage == nil && len(groupInfo.Messages) > 0 {
		targetMessage = &groupInfo.Messages[0]
		log.Printf("📸 Usando primeira mensagem: %d", targetMessage.MessageID)
	}

	if targetMessage == nil {
		log.Printf("❌ Nenhuma mensagem encontrada no grupo: %s", groupID)
		return
	}

	// ✅ PROCESSAR APENAS BOTÕES SE NÃO PODE EDITAR
	if !groupInfo.MessageEditAllowed {
		if len(buttons) == 0 {
			log.Printf("📸 Sem botões para adicionar ao grupo: %s", groupID)
			return
		}
		keyboard := mp.CreateInlineKeyboard(buttons, nil)
		if keyboard == nil {
			return
		}

		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   targetMessage.MessageID,
			ReplyMarkup: keyboard,
		})
		if err != nil {
			log.Printf("❌ Erro ao editar markup do grupo: %v", err)
		} else {
			log.Printf("✅ Markup editado para grupo: %s, mensagem: %d", groupID, targetMessage.MessageID)
		}
		return
	}

	// ✅ PROCESSAR CAPTION DA MENSAGEM ALVO
	var finalMessage string
	var customCaption *dbmodels.CustomCaption

	if targetMessage.HasCaption {
		entities := convertInterfaceToMessageEntities(targetMessage.CaptionEntities)
		formattedCaption := processTextWithFormatting(targetMessage.Caption, entities)
		finalMessage, customCaption = mp.processMessageWithHashtagPreservingFormat(formattedCaption, channel)
		log.Printf("📸 Processando com caption formatado: %s", targetMessage.Caption)
	} else {
		if channel.DefaultCaption != nil {
			finalMessage = channel.DefaultCaption.Caption
		}
		log.Printf("📸 Usando caption padrão")
	}

	keyboard := mp.CreateInlineKeyboard(buttons, customCaption)

	editParams := &bot.EditMessageCaptionParams{
		ChatID:    chatID,
		MessageID: targetMessage.MessageID,
		Caption:   finalMessage,
		ParseMode: "HTML",
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	// ✅ EDITAR APENAS A MENSAGEM ALVO
	_, err := mp.bot.EditMessageCaption(ctx, editParams)
	if err != nil {
		log.Printf("❌ Erro ao editar caption do grupo: %v", err)
	} else {
		log.Printf("✅ SUCESSO: Grupo %s processado - APENAS mensagem %d editada", groupID, targetMessage.MessageID)
	}

	// ✅ CLEANUP após 15 segundos
	time.AfterFunc(15*time.Second, func() {
		mp.mediaGroupManager.groups.Delete(groupID)
		groupMutexes.Delete(groupID) // ✅ Limpar mutex também
		log.Printf("🧹 Grupo removido da memória: %s", groupID)
	})
}

// ✅ NOVA FUNÇÃO: Retry para edit caption
func (mp *MessageProcessor) editMessageCaptionWithRetry(ctx context.Context, params *bot.EditMessageCaptionParams) error {
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Context com timeout para cada tentativa
		attemptCtx, cancel := context.WithTimeout(ctx, 15*time.Second)

		_, err := mp.bot.EditMessageCaption(attemptCtx, params)
		cancel()

		if err == nil {
			return nil
		}

		// Verificar tipos específicos de erro
		if strings.Contains(err.Error(), "context canceled") {
			if attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(attempt+1)
				log.Printf("Context canceled, retry %d/%d after %v", attempt+1, maxRetries, delay)
				time.Sleep(delay)
				continue
			}
		}

		if strings.Contains(err.Error(), "Message is not modified") {
			log.Printf("Caption not modified, skipping edit")
			return nil
		}

		if strings.Contains(err.Error(), "Bad Request") {
			log.Printf("Bad request error: %v", err)
			return err // Não retry para bad requests
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(attempt+1)
			log.Printf("Caption edit failed, retry %d/%d after %v: %v", attempt+1, maxRetries, delay, err)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("failed to edit caption after %d attempts", maxRetries)
}

// ✅ FUNÇÃO AUXILIAR: Converter MessageEntity para interface{}
func convertMessageEntitiesToInterface(entities []models.MessageEntity) []interface{} {
	result := make([]interface{}, len(entities))
	for i, entity := range entities {
		result[i] = entity
	}
	return result
}

func convertToInterfaceSlice[T any](s []T) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
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

// ✅ FUNÇÃO AUXILIAR: Converter interface{} para MessageEntity
func convertInterfaceToMessageEntities(entities []interface{}) []models.MessageEntity {
	result := make([]models.MessageEntity, 0, len(entities))
	for _, entity := range entities {
		if msgEntity, ok := entity.(models.MessageEntity); ok {
			result = append(result, msgEntity)
		}
	}
	return result
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

// ✅ CORRIGIDO: Preservar formatação original + adicionar legenda do banco
func (mp *MessageProcessor) processMessageWithHashtag(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
	hashtag := extractHashtag(text)
	var customCaption *dbmodels.CustomCaption

	if hashtag != "" {
		customCaption = findCustomCaption(channel, hashtag)
		cleanText := removeHashtag(text, hashtag)

		if customCaption != nil {
			// ✅ PRESERVAR formatação: usar \n\n para separar sem quebrar entidades
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

func (mp *MessageProcessor) IsNewPackActive(channelID int64) bool {
	return mp.mediaGroupManager.IsNewPackActive(channelID)
}

func (mp *MessageProcessor) SetNewPackActive(channelID int64, active bool) {
	mp.mediaGroupManager.SetNewPackActive(channelID, active)
}
