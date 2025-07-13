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

// ✅ ESTRUTURA PARA VERIFICAÇÃO DE PERMISSÕES
type PermissionCheckResult struct {
	CanEdit           bool
	CanAddButtons     bool
	CanEditButtons    bool
	CanUseLinkPreview bool
	Reason            string
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

// ✅ CORRIGIDO: Verificar permissões usando a estrutura correta dos models
func (mp *MessageProcessor) CheckPermissions(channel *dbmodels.Channel, messageType MessageType) *PermissionCheckResult {
	result := &PermissionCheckResult{
		CanEdit:           true,
		CanAddButtons:     true,
		CanEditButtons:    true,
		CanUseLinkPreview: true,
	}

	if channel == nil {
		result.CanEdit = false
		result.Reason = "Canal não encontrado"
		return result
	}

	// ✅ VERIFICAR SE EXISTE DefaultCaption E SUAS PERMISSÕES
	if channel.DefaultCaption == nil {
		log.Printf("⚠️ Canal %d não tem DefaultCaption configurado - permitindo todas as operações", channel.ID)
		return result
	}

	// ✅ VERIFICAR MessagePermission
	if channel.DefaultCaption.MessagePermission != nil {
		messagePermission := channel.DefaultCaption.MessagePermission

		// ✅ VERIFICAR LinkPreview APENAS PARA TEXTO
		if messageType == MessageTypeText && !messagePermission.LinkPreview {
			result.CanUseLinkPreview = false
			log.Printf("🔗 Link preview desabilitado para canal %d (MessagePermission.LinkPreview = false)", channel.ID)
		}

		// Verificar permissão específica por tipo de mensagem
		switch messageType {
		case MessageTypeText:
			if !messagePermission.Message {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de texto desabilitada"
			}
		case MessageTypeAudio:
			if !messagePermission.Audio {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de áudio desabilitada"
			}
		case MessageTypeVideo:
			if !messagePermission.Video {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de vídeo desabilitada"
			}
		case MessageTypePhoto:
			if !messagePermission.Photo {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de foto desabilitada"
			}
		case MessageTypeSticker:
			if !messagePermission.Sticker {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de sticker desabilitada"
			}
		case MessageTypeAnimation:
			if !messagePermission.GIF {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de GIF desabilitada"
			}
		}
	}

	// ✅ VERIFICAR ButtonsPermission
	if channel.DefaultCaption.ButtonsPermission != nil {
		buttonsPermission := channel.DefaultCaption.ButtonsPermission

		switch messageType {
		case MessageTypeText:
			if !buttonsPermission.Message {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAudio:
			if !buttonsPermission.Audio {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeVideo:
			if !buttonsPermission.Video {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypePhoto:
			if !buttonsPermission.Photo {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeSticker:
			if !buttonsPermission.Sticker {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAnimation:
			if !buttonsPermission.GIF {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		}
	}

	if !result.CanEdit {
		log.Printf("❌ Edição bloqueada para canal %d, tipo %s: %s", channel.ID, messageType, result.Reason)
	}

	if !result.CanAddButtons {
		log.Printf("🔘 Botões padrão bloqueados para canal %d, tipo %s", channel.ID, messageType)
	}

	return result
}

// ✅ CORRIGIDO: CheckCustomCaptionPermissions com messageType
func (mp *MessageProcessor) CheckCustomCaptionPermissions(channel *dbmodels.Channel, customCaption *dbmodels.CustomCaption, messageType MessageType) *PermissionCheckResult {
	result := &PermissionCheckResult{
		CanEdit:           true,
		CanAddButtons:     true,
		CanEditButtons:    true,
		CanUseLinkPreview: true,
	}

	if channel == nil {
		result.CanEdit = false
		result.Reason = "Canal não encontrado"
		return result
	}

	// ✅ USAR APENAS OS CAMPOS QUE EXISTEM NOS MODELS
	if channel.DefaultCaption == nil {
		log.Printf("⚠️ Canal %d não tem DefaultCaption configurado", channel.ID)
		return result
	}

	// ✅ VERIFICAR MessagePermission
	if channel.DefaultCaption.MessagePermission != nil {
		messagePermission := channel.DefaultCaption.MessagePermission

		// ✅ LinkPreview APENAS para texto
		if messageType == MessageTypeText && !messagePermission.LinkPreview {
			result.CanUseLinkPreview = false
			log.Printf("🔗 Link preview desabilitado por MessagePermission para canal %d", channel.ID)
		}

		// Verificar permissão por tipo
		switch messageType {
		case MessageTypeText:
			if !messagePermission.Message {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de texto desabilitada"
			}
		case MessageTypeAudio:
			if !messagePermission.Audio {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de áudio desabilitada"
			}
		case MessageTypeVideo:
			if !messagePermission.Video {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de vídeo desabilitada"
			}
		case MessageTypePhoto:
			if !messagePermission.Photo {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de foto desabilitada"
			}
		case MessageTypeSticker:
			if !messagePermission.Sticker {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de sticker desabilitada"
			}
		case MessageTypeAnimation:
			if !messagePermission.GIF {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de GIF desabilitada"
			}
		}
	}

	// ✅ VERIFICAR ButtonsPermission APENAS para botões padrão
	if channel.DefaultCaption.ButtonsPermission != nil {
		buttonsPermission := channel.DefaultCaption.ButtonsPermission

		switch messageType {
		case MessageTypeText:
			if !buttonsPermission.Message {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAudio:
			if !buttonsPermission.Audio {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeVideo:
			if !buttonsPermission.Video {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypePhoto:
			if !buttonsPermission.Photo {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeSticker:
			if !buttonsPermission.Sticker {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAnimation:
			if !buttonsPermission.GIF {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		}
	}

	// ✅ VERIFICAR CustomCaption.LinkPreview APENAS para texto
	if customCaption != nil && messageType == MessageTypeText {
		if !customCaption.LinkPreview {
			result.CanUseLinkPreview = false
			log.Printf("🔗 Link preview desabilitado por CustomCaption %s para canal %d", customCaption.Code, channel.ID)
		}
	}

	// ✅ LOGS ESPECÍFICOS PARA CUSTOM CAPTION
	if customCaption != nil {
		log.Printf("✅ Custom caption %s: %d botões (sempre permitidos)", customCaption.Code, len(customCaption.Buttons))
		log.Printf("✅ Permissões verificadas - Edit=%v, BotõesPadrão=%v, LinkPreview=%v",
			result.CanEdit, result.CanAddButtons, result.CanUseLinkPreview)
	}

	return result
}

// ✅ CORRIGIDO: ApplyPermissions sem afetar botões de custom captions
func (mp *MessageProcessor) ApplyPermissions(channel *dbmodels.Channel, messageType MessageType, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button) (bool, []dbmodels.Button, *dbmodels.CustomCaption) {
	permissions := mp.CheckCustomCaptionPermissions(channel, customCaption, messageType)

	if !permissions.CanEdit {
		log.Printf("❌ Edição de mensagem bloqueada: %s", permissions.Reason)
		return false, nil, nil
	}

	// ✅ VERIFICAR ButtonsPermissions APENAS para botões padrão do canal
	if !permissions.CanAddButtons {
		log.Printf("⚠️ Botões padrão do canal removidos devido a ButtonsPermissions")
		buttons = nil
		// ✅ NÃO REMOVER botões da custom caption - eles são independentes
		log.Printf("✅ Botões de custom caption mantidos (independentes de ButtonsPermissions)")
	}

	return true, buttons, customCaption
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

// ✅ CORRIGIDO: Priorizar botões do custom caption SEM verificação de ButtonsPermissions
func (mp *MessageProcessor) CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption, channel *dbmodels.Channel, messageType MessageType) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button

	// ✅ PRIORIDADE: Se tem custom caption, usar APENAS seus botões (SEM verificar ButtonsPermissions)
	if customCaption != nil && len(customCaption.Buttons) > 0 {
		log.Printf("🔘 Usando botões do custom caption: %s (%d botões) - IGNORANDO ButtonsPermissions", customCaption.Code, len(customCaption.Buttons))
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
		// ✅ FALLBACK: Usar botões padrão do canal (COM verificação de ButtonsPermissions)
		permissions := mp.CheckPermissions(channel, messageType)
		if !permissions.CanAddButtons {
			log.Printf("🔘 Botões padrão bloqueados: ButtonsPermissions para canal %d", channel.ID)
			return nil
		}
		log.Printf("🔘 Usando botões padrão do canal (%d botões)", len(buttons))
		finalButtons = buttons
	}

	if len(finalButtons) == 0 {
		log.Printf("🔘 Nenhum botão disponível")
		return nil
	}

	// Criar grid de botões
	buttonGrid := make(map[int]map[int]models.InlineKeyboardButton)
	for i, button := range finalButtons {
		if button.NameButton == "" || button.ButtonURL == "" {
			log.Printf("⚠️ Botão inválido ignorado: %+v", button)
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

// ✅ CORRIGIDO: ProcessTextMessage com formatação HTML e LinkPreview
func (mp *MessageProcessor) ProcessTextMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	text := post.Text
	messageID := post.ID
	messageType := MessageTypeText

	if text == "" {
		return fmt.Errorf("texto da mensagem está vazio")
	}

	// ✅ VERIFICAR PERMISSÕES
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Edição de texto bloqueada para canal %d: %s", channel.ID, permissions.Reason)
		return fmt.Errorf("permissão de edição de texto desabilitada")
	}

	if !messageEditAllowed {
		if len(buttons) == 0 || !permissions.CanAddButtons {
			return nil
		}

		keyboard := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
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

	// ✅ APLICAR FORMATAÇÃO HTML
	formattedText := processTextWithFormatting(text, post.Entities)

	message, customCaption := mp.processMessageWithHashtag(formattedText, channel)

	// ✅ APLICAR VERIFICAÇÕES DE PERMISSÃO
	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		return fmt.Errorf("permissões insuficientes para editar mensagem")
	}

	keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)

	editParams := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      message,
		ParseMode: "HTML",
	}

	// ✅ VERIFICAR LINK PREVIEW: MessagePermission.LinkPreview E CustomCaption.LinkPreview
	disableLinkPreview := false

	// 1. Verificar MessagePermission.LinkPreview
	if !permissions.CanUseLinkPreview {
		disableLinkPreview = true
		log.Printf("🔗 Link preview desabilitado por MessagePermission para canal %d", channel.ID)
	}

	// 2. Verificar CustomCaption.LinkPreview (se existe custom caption)
	if customCaption != nil && !customCaption.LinkPreview {
		disableLinkPreview = true
		log.Printf("🔗 Link preview desabilitado por CustomCaption %s para canal %d", customCaption.Code, channel.ID)
	}

	// ✅ USAR LinkPreviewOptions ao invés de DisableWebPagePreview
	if disableLinkPreview {
		val := true
		editParams.LinkPreviewOptions = &models.LinkPreviewOptions{
			IsDisabled: &val,
		}
		log.Printf("🔗 Link preview DESABILITADO para mensagem de texto no canal %d", channel.ID)
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageText(ctx, editParams)
	return err
}

// ✅ CORRIGIDO: ProcessAudioMessage - SUBSTITUIR caption por legenda do banco
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	messageID := post.ID
	caption := post.Caption
	mediaGroupID := post.MediaGroupID
	messageType := MessageTypeAudio

	// ✅ VERIFICAR PERMISSÕES
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Edição de áudio bloqueada para canal %d: %s", channel.ID, permissions.Reason)
		return fmt.Errorf("permissão de edição de áudio desabilitada")
	}

	if !messageEditAllowed {
		if len(buttons) == 0 || !permissions.CanAddButtons {
			return nil
		}

		keyboard := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
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

	// ✅ APLICAR FORMATAÇÃO HTML ANTES DE QUALQUER PROCESSAMENTO
	formattedCaption := processTextWithFormatting(caption, post.CaptionEntities)

	// Para grupos de mídia: REENVIAR + DELETAR
	if mediaGroupID != "" {
		var finalMessage string
		var customCaption *dbmodels.CustomCaption

		// ✅ VERIFICAR SE TEM HASHTAG NA CAPTION ORIGINAL
		hashtag := extractHashtag(formattedCaption)
		if hashtag != "" {
			customCaption = findCustomCaption(channel, hashtag)
			if customCaption != nil {
				// ✅ SUBSTITUIR COMPLETAMENTE pela custom caption
				finalMessage = detectParseMode(customCaption.Caption)
			} else {
				// ✅ SUBSTITUIR COMPLETAMENTE pela default caption
				if channel.DefaultCaption != nil {
					finalMessage = detectParseMode(channel.DefaultCaption.Caption)
				}
			}
		} else {
			// ✅ SEM HASHTAG: SUBSTITUIR pela default caption
			if channel.DefaultCaption != nil {
				finalMessage = detectParseMode(channel.DefaultCaption.Caption)
			}
		}

		// ✅ APLICAR VERIFICAÇÕES DE PERMISSÃO
		canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
		if !canEdit {
			return fmt.Errorf("permissões insuficientes para editar mensagem")
		}

		keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)

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

	// ✅ VERIFICAR SE TEM HASHTAG NA CAPTION ORIGINAL
	hashtag := extractHashtag(formattedCaption)
	if hashtag != "" {
		customCaption = findCustomCaption(channel, hashtag)
		if customCaption != nil {
			// ✅ SUBSTITUIR COMPLETAMENTE pela custom caption
			finalMessage = detectParseMode(customCaption.Caption)
		} else {
			// ✅ SUBSTITUIR COMPLETAMENTE pela default caption
			if channel.DefaultCaption != nil {
				finalMessage = detectParseMode(channel.DefaultCaption.Caption)
			}
		}
	} else {
		// ✅ SEM HASHTAG: SUBSTITUIR pela default caption
		if channel.DefaultCaption != nil {
			finalMessage = detectParseMode(channel.DefaultCaption.Caption)
		}
	}

	// ✅ APLICAR VERIFICAÇÕES DE PERMISSÃO
	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		return fmt.Errorf("permissões insuficientes para editar mensagem")
	}

	keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)

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

// ✅ CORRIGIDO: ProcessMediaMessage com verificação de permissões
func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	var messageType MessageType
	if post.Photo != nil {
		messageType = MessageTypePhoto
	} else if post.Video != nil {
		messageType = MessageTypeVideo
	} else if post.Animation != nil {
		messageType = MessageTypeAnimation
	} else {
		return fmt.Errorf("tipo de mídia não suportado")
	}

	// ✅ VERIFICAR PERMISSÕES
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Edição de mídia bloqueada para canal %d: %s", channel.ID, permissions.Reason)
		return fmt.Errorf("permissão de edição de mídia desabilitada")
	}

	mediaGroupID := post.MediaGroupID
	if mediaGroupID != "" {
		return mp.handleGroupedMedia(ctx, channel, post, buttons, messageEditAllowed, messageType)
	}
	return mp.handleSingleMedia(ctx, channel, post, buttons, messageEditAllowed, messageType)
}

// ✅ CORRIGIDO: handleSingleMedia com formatação HTML
func (mp *MessageProcessor) handleSingleMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool, messageType MessageType) error {
	messageID := post.ID
	caption := post.Caption

	permissions := mp.CheckPermissions(channel, messageType)

	if !messageEditAllowed {
		if len(buttons) == 0 || !permissions.CanAddButtons {
			return nil
		}

		keyboard := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
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
	isShowCaptionAboveMedia := false
	if post.ShowCaptionAboveMedia {
		isShowCaptionAboveMedia = true
	}

	// ✅ APLICAR FORMATAÇÃO HTML NA CAPTION
	formattedCaption := processTextWithFormatting(caption, post.CaptionEntities)

	message, customCaption := mp.processMessageWithHashtag(formattedCaption, channel)

	// ✅ APLICAR VERIFICAÇÕES DE PERMISSÃO
	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		return fmt.Errorf("permissões insuficientes para editar mensagem")
	}

	keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)

	editParams := &bot.EditMessageCaptionParams{
		ChatID:                post.Chat.ID,
		MessageID:             messageID,
		Caption:               message,
		ParseMode:             "HTML",
		ShowCaptionAboveMedia: isShowCaptionAboveMedia,
	}

	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}

	_, err := mp.bot.EditMessageCaption(ctx, editParams)
	return err
}

// ✅ CORRIGIDO: handleGroupedMedia com verificação de permissões
func (mp *MessageProcessor) handleGroupedMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool, messageType MessageType) error {
	mediaGroupID := post.MediaGroupID
	messageID := post.ID
	caption := post.Caption

	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		return fmt.Errorf("permissão de edição de grupo de mídia desabilitada")
	}

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

	// ✅ TIMEOUT ADAPTATIVO
	timeout := time.Duration(800+len(group.Messages)*200) * time.Millisecond
	if timeout > 2*time.Second {
		timeout = 2 * time.Second
	}

	// ✅ CRIAR TIMER
	group.Timer = time.AfterFunc(timeout, func() {
		mp.finishGroupProcessing(ctx, mediaGroupID, channel, buttons, messageType)
	})

	return nil
}

// ✅ CORRIGIDO: finishGroupProcessing com formatação HTML
func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button, messageType MessageType) {
	log.Printf("📸 Iniciando processamento final do grupo: %s", groupID)

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
	log.Printf("📸 Marcando grupo como processado: %s com %d mensagens", groupID, len(group.Messages))

	if len(group.Messages) == 0 {
		return
	}

	// ✅ VERIFICAR PERMISSÕES
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		mp.cleanupGroup(groupID)
		return
	}

	// ✅ ENCONTRAR A MENSAGEM IDEAL PARA EDITAR
	var targetMessage *MediaMessage
	// Prioridade 1: Mensagem com caption
	for i := range group.Messages {
		if group.Messages[i].HasCaption {
			targetMessage = &group.Messages[i]
			break
		}
	}

	// Prioridade 2: Primeira mensagem se não houver caption
	if targetMessage == nil {
		targetMessage = &group.Messages[0]
	}

	// ✅ SE NÃO PODE EDITAR MENSAGEM, APENAS ADICIONAR BOTÕES (se permitido)
	if !group.MessageEditAllowed {
		if len(buttons) > 0 && permissions.CanAddButtons {
			keyboard := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
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

	// ✅ PROCESSAR CAPTION COM FORMATAÇÃO HTML
	var finalMessage string
	var customCaption *dbmodels.CustomCaption

	if targetMessage.HasCaption {
		// ✅ APLICAR FORMATAÇÃO HTML ANTES DE PROCESSAR
		entities := convertInterfaceToMessageEntities(targetMessage.CaptionEntities)
		formattedCaption := processTextWithFormatting(targetMessage.Caption, entities)
		// ✅ PROCESSAR HASHTAG E OBTER CUSTOM CAPTION COM CAPTION FORMATADA
		finalMessage, customCaption = mp.processMessageWithHashtag(formattedCaption, channel)
		if customCaption != nil {
			log.Printf("📸 Custom caption encontrado: %s", customCaption.Code)
		}
	} else {
		// ✅ USAR CAPTION PADRÃO FORMATADO se não houver caption na mensagem
		if channel.DefaultCaption != nil {
			finalMessage = detectParseMode(channel.DefaultCaption.Caption)
		}
	}

	// ✅ APLICAR VERIFICAÇÕES DE PERMISSÃO
	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		mp.cleanupGroup(groupID)
		return
	}

	// ✅ CRIAR KEYBOARD COM CUSTOM CAPTION BUTTONS
	keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)

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
		if customCaption != nil {
			log.Printf("✅ Custom caption aplicado: %s com %d botões", customCaption.Code, len(customCaption.Buttons))
		}
	}

	// ✅ CLEANUP
	mp.cleanupGroup(groupID)
}

// ✅ FUNÇÃO PARA LIMPEZA DO GRUPO
func (mp *MessageProcessor) cleanupGroup(groupID string) {
	time.AfterFunc(10*time.Second, func() {
		mediaGroups.Delete(groupID)
	})
}

// ✅ CORRIGIDO: ProcessStickerMessage com verificação de permissões
func (mp *MessageProcessor) ProcessStickerMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button) error {
	messageType := MessageTypeSticker
	permissions := mp.CheckPermissions(channel, messageType)

	if len(buttons) == 0 || !permissions.CanAddButtons {
		return nil
	}

	keyboard := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
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
		hashtag := strings.ToLower(matches[1])
		return hashtag
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

	log.Printf("📝 ❌ Custom caption não encontrado para: #%s", hashtag)
	customCaptionCache.Store(cacheKey, (*dbmodels.CustomCaption)(nil))
	return nil
}

// ✅ CORRIGIDO: processMessageWithHashtag com formatação das legendas do banco
func (mp *MessageProcessor) processMessageWithHashtaga(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
	hashtag := extractHashtag(text)

	if hashtag == "" {
		defaultCaption := ""
		if channel.DefaultCaption != nil {
			// ✅ FORMATAR CAPTION PADRÃO DO BANCO
			defaultCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
	}

	customCaption := findCustomCaption(channel, hashtag)
	if customCaption == nil {
		defaultCaption := ""
		if channel.DefaultCaption != nil {
			defaultCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
	}

	cleanText := removeHashtag(text, hashtag)

	// ✅ FORMATAR CUSTOM CAPTION DO BANCO
	formattedCustomCaption := detectParseMode(customCaption.Caption)

	return fmt.Sprintf("%s\n\n%s", cleanText, formattedCustomCaption), customCaption
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
