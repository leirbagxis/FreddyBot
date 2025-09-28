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

var (
	hashtagRegex         = regexp.MustCompile(`#(\w+)`)
	removeHashRegexCache = sync.Map{}
	customCaptionCache   = sync.Map{}
	mediaGroups          = sync.Map{}
)

type MediaGroup struct {
	Messages           []MediaMessage
	Processed          bool
	Timer              *time.Timer
	MessageEditAllowed bool
	ChatID             int64
	mu                 sync.Mutex
}

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

/*
	UTILS DE COMPOSIÇÃO
*/

// composeMessage combina o conteúdo original com uma legenda do banco.
// order: "append" -> original + sep + db; "prepend" -> db + sep + original
func composeMessage(original, fromDB, sep, order string) string {
	o := strings.TrimSpace(original)
	d := strings.TrimSpace(fromDB)
	if o == "" && d == "" {
		return ""
	}
	if o == "" {
		return d
	}
	if d == "" {
		return o
	}
	if sep == "" {
		sep = "\n\n"
	}
	if order == "prepend" {
		return d + sep + o
	}
	return o + sep + d
}

/*
	PERMISSÕES
*/

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

	if channel.DefaultCaption != nil && channel.DefaultCaption.MessagePermission != nil {
		mpPerm := channel.DefaultCaption.MessagePermission
		if messageType == MessageTypeText && !mpPerm.LinkPreview {
			result.CanUseLinkPreview = false
		}
		switch messageType {
		case MessageTypeText:
			if !mpPerm.Message {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de texto desabilitada"
			}
		case MessageTypeAudio:
			if !mpPerm.Audio {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de áudio desabilitada"
			}
		case MessageTypeVideo:
			if !mpPerm.Video {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de vídeo desabilitada"
			}
		case MessageTypePhoto:
			if !mpPerm.Photo {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de foto desabilitada"
			}
		case MessageTypeSticker:
			if !mpPerm.Sticker {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de sticker desabilitada"
			}
		case MessageTypeAnimation:
			if !mpPerm.GIF {
				result.CanEdit = false
				result.Reason = "Edição de mensagens de GIF desabilitada"
			}
		}
	}

	if channel.DefaultCaption != nil && channel.DefaultCaption.ButtonsPermission != nil {
		bp := channel.DefaultCaption.ButtonsPermission
		switch messageType {
		case MessageTypeText:
			if !bp.Message {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAudio:
			if !bp.Audio {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeVideo:
			if !bp.Video {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypePhoto:
			if !bp.Photo {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeSticker:
			if !bp.Sticker {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case MessageTypeAnimation:
			if !bp.GIF {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		}
	}

	return result
}

func (mp *MessageProcessor) CheckCustomCaptionPermissions(channel *dbmodels.Channel, customCaption *dbmodels.CustomCaption, messageType MessageType) *PermissionCheckResult {
	result := mp.CheckPermissions(channel, messageType)
	if customCaption != nil && messageType == MessageTypeText && !customCaption.LinkPreview {
		result.CanUseLinkPreview = false
	}
	return result
}

// Degradação: nunca bloqueia o fluxo, apenas filtra botões padrão.
func (mp *MessageProcessor) ApplyPermissions(channel *dbmodels.Channel, messageType MessageType, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button) (bool, []dbmodels.Button, *dbmodels.CustomCaption) {
	perms := mp.CheckCustomCaptionPermissions(channel, customCaption, messageType)
	if !perms.CanAddButtons {
		buttons = nil
	}
	return true, buttons, customCaption
}

/*
	KEYBOARD
*/

func (mp *MessageProcessor) CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption, channel *dbmodels.Channel, messageType MessageType) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button

	if customCaption != nil && len(customCaption.Buttons) > 0 {
		for _, cb := range customCaption.Buttons {
			finalButtons = append(finalButtons, dbmodels.Button{
				NameButton: cb.NameButton,
				ButtonURL:  cb.ButtonURL,
				PositionY:  cb.PositionY,
				PositionX:  cb.PositionX,
			})
		}
	} else {
		perms := mp.CheckPermissions(channel, messageType)
		if !perms.CanAddButtons {
			return nil
		}
		finalButtons = buttons
	}

	if len(finalButtons) == 0 {
		return nil
	}

	// Construção simples por linhas
	rows := map[int][]models.InlineKeyboardButton{}
	for _, b := range finalButtons {
		if b.NameButton == "" || b.ButtonURL == "" {
			continue
		}
		row := b.PositionY
		if row < 0 {
			row = 0
		}
		btn := models.InlineKeyboardButton{Text: b.NameButton, URL: b.ButtonURL}
		rows[row] = append(rows[row], btn)
	}

	// Ordenar por linha
	keyboard := make([][]models.InlineKeyboardButton, 0, len(rows))
	for r := 0; r < 20; r++ {
		if line, ok := rows[r]; ok && len(line) > 0 {
			keyboard = append(keyboard, line)
		}
	}
	if len(keyboard) == 0 {
		return nil
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

/*
	DISPATCH
*/

func (mp *MessageProcessor) GetMessageType(post *models.Message) MessageType {
	switch {
	case post.Text != "":
		return MessageTypeText
	case post.Audio != nil:
		return MessageTypeAudio
	case post.Sticker != nil:
		return MessageTypeSticker
	case post.Photo != nil:
		return MessageTypePhoto
	case post.Video != nil:
		return MessageTypeVideo
	case post.Animation != nil:
		return MessageTypeAnimation
	default:
		return ""
	}
}

func (mp *MessageProcessor) ProcessMessagea(ctx context.Context, messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool) error {
	switch messageType {
	case MessageTypeText:
		return mp.ProcessTextMessage(ctx, channel, post, buttons, true)
	case MessageTypeAudio:
		return mp.ProcessAudioMessage(ctx, channel, post, buttons, true)
	case MessageTypeSticker:
		if len(buttons) > 0 {
			return mp.ProcessStickerMessage(ctx, channel, post, buttons)
		}
		return nil
	case MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation:
		return mp.ProcessMediaMessage(ctx, channel, post, buttons, true)
	default:
		return nil
	}
}

/*
	TEXT
*/

func (mp *MessageProcessor) ProcessTextMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool) error {
	messageID := post.ID
	text := post.Text
	if text == "" {
		return nil
	}

	perms := mp.CheckPermissions(channel, MessageTypeText)
	// Degradação: sem editar texto
	if !perms.CanEdit {
		if len(buttons) == 0 || !perms.CanAddButtons {
			return nil
		}
		kb := mp.CreateInlineKeyboard(buttons, nil, channel, MessageTypeText)
		if kb == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   messageID,
			ReplyMarkup: kb,
		})
		return err
	}

	// Pode editar: compor original + db
	formatted := processTextWithFormatting(text, post.Entities)
	var custom *dbmodels.CustomCaption
	var dbCaption string
	// Detectar hashtag para custom
	if h := extractHashtag(formatted); h != "" {
		if cc := findCustomCaption(channel, h); cc != nil {
			custom = cc
			dbCaption = detectParseMode(cc.Caption)
		}
	}
	if custom == nil && channel.DefaultCaption != nil {
		dbCaption = detectParseMode(channel.DefaultCaption.Caption)
	}
	finalText := composeMessage(formatted, dbCaption, "\n\n", "append")

	_, filteredButtons, allowedCustom := mp.ApplyPermissions(channel, MessageTypeText, custom, buttons)
	kb := mp.CreateInlineKeyboard(filteredButtons, allowedCustom, channel, MessageTypeText)

	edit := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      finalText,
		ParseMode: "HTML",
	}
	disableLP := !perms.CanUseLinkPreview
	if custom != nil && !custom.LinkPreview {
		disableLP = true
	}
	if disableLP {
		val := true
		edit.LinkPreviewOptions = &models.LinkPreviewOptions{IsDisabled: &val}
	}
	if kb != nil {
		edit.ReplyMarkup = kb
	}
	_, err := mp.bot.EditMessageText(ctx, edit)
	return err
}

/*
	STICKER (apenas botões)
*/

func (mp *MessageProcessor) ProcessStickerMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button) error {
	perms := mp.CheckPermissions(channel, MessageTypeSticker)
	if len(buttons) == 0 || !perms.CanAddButtons {
		return nil
	}
	kb := mp.CreateInlineKeyboard(buttons, nil, channel, MessageTypeSticker)
	if kb == nil {
		return nil
	}
	_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      post.Chat.ID,
		MessageID:   post.ID,
		ReplyMarkup: kb,
	})
	return err
}

/*
	AUDIO
*/

func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool) error {
	perms := mp.CheckPermissions(channel, MessageTypeAudio)

	// Degradação: apenas markup
	if !perms.CanEdit {
		if len(buttons) == 0 || !perms.CanAddButtons {
			return nil
		}
		kb := mp.CreateInlineKeyboard(buttons, nil, channel, MessageTypeAudio)
		if kb == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   post.ID,
			ReplyMarkup: kb,
		})
		return err
	}

	// Pode editar
	time.Sleep(1500 * time.Millisecond)
	formattedCaption := processTextWithFormatting(post.Caption, post.CaptionEntities)

	if post.MediaGroupID != "" {
		return mp.processAudioInGroupInPlace(ctx, channel, post, buttons, formattedCaption, MessageTypeAudio)
	}
	return mp.processSingleAudio(ctx, channel, post, buttons, formattedCaption, MessageTypeAudio)
}

// Edição in-place para grupos de áudio
func (mp *MessageProcessor) processAudioInGroupInPlace(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, formattedCaption string, messageType MessageType) error {
	mediaGroupID := post.MediaGroupID

	value, _ := mediaGroups.LoadOrStore(mediaGroupID, &MediaGroup{
		Messages:           make([]MediaMessage, 0),
		Processed:          false,
		MessageEditAllowed: true,
		ChatID:             post.Chat.ID,
	})
	group := value.(*MediaGroup)

	group.mu.Lock()
	group.Messages = append(group.Messages, MediaMessage{
		MessageID:       post.ID,
		HasCaption:      post.Caption != "",
		Caption:         post.Caption,
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
		mp.finishGroupProcessingAudioInPlace(ctx, mediaGroupID, channel, buttons, messageType)
	})
	group.mu.Unlock()

	return nil
}

func (mp *MessageProcessor) finishGroupProcessingAudioInPlace(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button, messageType MessageType) {
	value, ok := mediaGroups.Load(groupID)
	if !ok {
		return
	}
	group := value.(*MediaGroup)

	group.mu.Lock()
	if group.Processed {
		group.mu.Unlock()
		return
	}
	group.Processed = true

	var baseCaption string
	var baseEntities []interface{}
	targetMessageID := 0
	for _, m := range group.Messages {
		if m.HasCaption && targetMessageID == 0 {
			targetMessageID = m.MessageID
			baseCaption = m.Caption
			baseEntities = m.CaptionEntities
			break
		}
	}
	if targetMessageID == 0 && len(group.Messages) > 0 {
		targetMessageID = group.Messages[0].MessageID
		baseCaption = group.Messages[0].Caption
		baseEntities = group.Messages[0].CaptionEntities
	}
	group.mu.Unlock()

	perms := mp.CheckPermissions(channel, messageType)
	if !perms.CanEdit {
		if len(buttons) == 0 || !perms.CanAddButtons {
			mediaGroups.Delete(groupID)
			return
		}
		kb := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
		if kb != nil {
			_, _ = mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
				ChatID:      group.ChatID,
				MessageID:   targetMessageID,
				ReplyMarkup: kb,
			})
		}
		mediaGroups.Delete(groupID)
		return
	}

	// Compor caption final (original + db)
	formatted := processTextWithFormatting(baseCaption, convertInterfaceToEntities(baseEntities))
	var dbCaption string
	var customCaption *dbmodels.CustomCaption
	if h := extractHashtag(formatted); h != "" {
		if cc := findCustomCaption(channel, h); cc != nil {
			customCaption = cc
			dbCaption = detectParseMode(cc.Caption)
		}
	}
	if customCaption == nil && channel.DefaultCaption != nil {
		dbCaption = detectParseMode(channel.DefaultCaption.Caption)
	}

	// Áudio em grupo: SUBSTITUIR se puder editar e houver dbCaption; caso contrário, manter original
	finalCaption := formatted
	if perms.CanEdit && dbCaption != "" {
		finalCaption = dbCaption
	}

	_, filteredButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	kb := mp.CreateInlineKeyboard(filteredButtons, allowedCustomCaption, channel, messageType)

	edit := &bot.EditMessageCaptionParams{
		ChatID:    group.ChatID,
		MessageID: targetMessageID,
		Caption:   finalCaption,
		ParseMode: "HTML",
	}
	if kb != nil {
		edit.ReplyMarkup = kb
	}
	_, _ = mp.bot.EditMessageCaption(ctx, edit)

	mediaGroups.Delete(groupID)
}

func (mp *MessageProcessor) processSingleAudio(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, formattedCaption string, messageType MessageType) error {
	// Determinar dbCaption (Custom > Default)
	var dbCaption string
	var customCaption *dbmodels.CustomCaption
	if h := extractHashtag(formattedCaption); h != "" {
		if cc := findCustomCaption(channel, h); cc != nil {
			customCaption = cc
			dbCaption = detectParseMode(cc.Caption)
		}
	}
	if customCaption == nil && channel.DefaultCaption != nil {
		dbCaption = detectParseMode(channel.DefaultCaption.Caption)
	}

	// Áudio: SUBSTITUIR se houver permissão e houver dbCaption; senão manter original
	finalCaption := formattedCaption
	perms := mp.CheckPermissions(channel, messageType)
	if perms.CanEdit && dbCaption != "" {
		finalCaption = dbCaption
	}

	_, filteredButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	kb := mp.CreateInlineKeyboard(filteredButtons, allowedCustomCaption, channel, messageType)

	edit := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: post.ID,
		Caption:   finalCaption,
		ParseMode: "HTML",
	}
	if kb != nil {
		edit.ReplyMarkup = kb
	}
	_, err := mp.bot.EditMessageCaption(ctx, edit)
	return err
}

/*
	MEDIA (PHOTO/VIDEO/ANIMATION)
*/

func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool) error {
	var messageType MessageType
	switch {
	case post.Photo != nil:
		messageType = MessageTypePhoto
	case post.Video != nil:
		messageType = MessageTypeVideo
	case post.Animation != nil:
		messageType = MessageTypeAnimation
	default:
		return fmt.Errorf("tipo de mídia não suportado")
	}

	mp.CheckPermissions(channel, messageType)
	if post.MediaGroupID != "" {
		return mp.handleGroupedMedia(ctx, channel, post, buttons, false, messageType)
	}
	return mp.handleSingleMedia(ctx, channel, post, buttons, false, messageType)
}

func (mp *MessageProcessor) handleSingleMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool, messageType MessageType) error {
	perms := mp.CheckPermissions(channel, messageType)
	// Degradação
	if !perms.CanEdit {
		if len(buttons) == 0 || !perms.CanAddButtons {
			return nil
		}
		kb := mp.CreateInlineKeyboard(buttons, nil, channel, messageType)
		if kb == nil {
			return nil
		}
		_, err := mp.bot.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   post.ID,
			ReplyMarkup: kb,
		})
		return err
	}

	// Compor caption: original + db
	formatted := processTextWithFormatting(post.Caption, post.CaptionEntities)
	var dbCaption string
	var customCaption *dbmodels.CustomCaption
	if h := extractHashtag(formatted); h != "" {
		if cc := findCustomCaption(channel, h); cc != nil {
			customCaption = cc
			dbCaption = detectParseMode(cc.Caption)
		}
	}
	if customCaption == nil && channel.DefaultCaption != nil {
		dbCaption = detectParseMode(channel.DefaultCaption.Caption)
	}
	finalCaption := composeMessage(formatted, dbCaption, "\n\n", "append")

	_, filteredButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	kb := mp.CreateInlineKeyboard(filteredButtons, allowedCustomCaption, channel, messageType)

	edit := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: post.ID,
		Caption:   finalCaption,
		ParseMode: "HTML",
	}
	if kb != nil {
		edit.ReplyMarkup = kb
	}
	_, err := mp.bot.EditMessageCaption(ctx, edit)
	return err
}

func (mp *MessageProcessor) handleGroupedMedia(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool, messageType MessageType) error {
	mediaGroupID := post.MediaGroupID
	messageID := post.ID
	caption := post.Caption

	value, loaded := mediaGroups.LoadOrStore(mediaGroupID, &MediaGroup{
		Messages:           make([]MediaMessage, 0),
		Processed:          false,
		MessageEditAllowed: true,
		ChatID:             post.Chat.ID,
	})
	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if !loaded {
		log.Printf("📸 Novo grupo criado: %s", mediaGroupID)
	}

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
		mp.finishGroupProcessing(ctx, mediaGroupID, channel, buttons, messageType)
	})
	return nil
}

func (mp *MessageProcessor) finishGroupProcessing(ctx context.Context, groupID string, channel *dbmodels.Channel, buttons []dbmodels.Button, messageType MessageType) {
	log.Printf("📸 Iniciando processamento final do grupo: %s", groupID)

	value, ok := mediaGroups.Load(groupID)
	if !ok {
		log.Printf("❌ Grupo %s não encontrado", groupID)
		return
	}

	group := value.(*MediaGroup)
	group.mu.Lock()
	defer group.mu.Unlock()

	if group.Processed {
		log.Printf("⚠️ Grupo %s já foi processado", groupID)
		return
	}

	group.Processed = true
	log.Printf("📸 Marcando grupo como processado: %s com %d mensagens", groupID, len(group.Messages))

	if len(group.Messages) == 0 {
		log.Printf("❌ Grupo %s não tem mensagens", groupID)
		return
	}

	// ✅ VERIFICAR PERMISSÕES
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Sem permissões para editar mensagens no grupo %s", groupID)
		mp.cleanupGroup(groupID)
		return
	}

	// ✅ ENCONTRAR A MENSAGEM IDEAL PARA EDITAR
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

	// ✅ PROCESSAR CAPTION E EDITAR MENSAGEM
	var finalMessage string
	var customCaption *dbmodels.CustomCaption

	if targetMessage.HasCaption {
		entities := convertInterfaceToEntities(targetMessage.CaptionEntities)
		formattedCaption := processTextWithFormatting(targetMessage.Caption, entities)
		finalMessage, customCaption = mp.processMessageWithHashtag(formattedCaption, channel)
	} else {
		if channel.DefaultCaption != nil {
			finalMessage = detectParseMode(channel.DefaultCaption.Caption)
		}
	}

	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		log.Printf("❌ Permissões insuficientes para editar grupo %s", groupID)
		mp.cleanupGroup(groupID)
		return
	}

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
		log.Printf("✅ Grupo %s processado - mensagem %d editada", groupID, targetMessage.MessageID)
	}

	// ✅ ENVIAR SEPARATOR APÓS EDITAR A MENSAGEM
	if channel.Separator != nil && (permissions.CanEdit || permissions.CanAddButtons) {
		time.Sleep(1 * time.Second) // Delay antes de enviar separator

		log.Printf("🔄 Tentando enviar separator para grupo %s (tipo: %s)", groupID, messageType)

		separatorCtx, separatorCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer separatorCancel()

		err := mp.ProcessSeparator(separatorCtx, channel, nil)
		if err != nil {
			log.Printf("❌ Erro ao processar separator para grupo %s: %v", groupID, err)
		} else {
			log.Printf("✅ Separator enviado com sucesso para grupo %s", groupID)
		}
	} else {
		if channel.Separator == nil {
			log.Printf("⚠️ Separator não configurado para canal %d", channel.ID)
		} else {
			log.Printf("⚠️ Sem permissões para enviar separator no grupo %s", groupID)
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

/*
	HELPERS
*/

// ✅ CORRIGIDO: ProcessStickerMessage com verificação de permissões
func (mp *MessageProcessor) ProcessStickerMessagea(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button) error {
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

// ✅ FUNÇÃO ProcessSeparator COM RETRY
func (mp *MessageProcessor) ProcessSeparator(ctx context.Context, channel *dbmodels.Channel, post *models.Message) error {
	if channel.Separator == nil || channel.Separator.SeparatorID == "" {
		log.Printf("⚠️ Separator não configurado para canal %d", channel.ID)
		return nil
	}

	var chatID int64
	if post != nil {
		chatID = post.Chat.ID
	} else {
		chatID = channel.ID // Fallback para grupos
	}

	log.Printf("🔄 Enviando separator para chat %d", chatID)

	// ✅ RETRY COM BACKOFF PARA SEPARATORS
	maxRetries := 2
	baseDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := mp.bot.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:  chatID,
			Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
		})

		if err == nil {
			log.Printf("✅ Separator enviado com sucesso para chat %d", chatID)
			return nil
		}

		// Verificar se é erro 429
		if strings.Contains(err.Error(), "Too Many Requests") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter == 0 {
				retryAfter = int(baseDelay.Seconds()) * (attempt + 1)
			}

			log.Printf("⏳ Rate limit no separator, aguardando %d segundos (tentativa %d/%d)", retryAfter, attempt+1, maxRetries)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		log.Printf("❌ Erro ao enviar separator: %v", err)
		return err
	}

	return fmt.Errorf("falha após %d tentativas no envio do separator", maxRetries)
}

func extractHashtag(text string) string {
	m := hashtagRegex.FindStringSubmatch(text)
	if len(m) > 1 {
		return strings.ToLower(m[1])
	}
	return ""
}

func findCustomCaption(channel *dbmodels.Channel, code string) *dbmodels.CustomCaption {
	if channel == nil || len(channel.CustomCaptions) == 0 {
		return nil
	}
	code = strings.ToLower(strings.TrimSpace(code))
	for i := range channel.CustomCaptions {
		if strings.ToLower(channel.CustomCaptions[i].Code) == code {
			return &channel.CustomCaptions[i]
		}
	}
	return nil
}

func convertMessageEntitiesToInterface(ents []models.MessageEntity) []interface{} {
	out := make([]interface{}, 0, len(ents))
	for _, e := range ents {
		out = append(out, e)
	}
	return out
}

func convertInterfaceToEntities(anys []interface{}) []models.MessageEntity {
	out := make([]models.MessageEntity, 0, len(anys))
	for _, v := range anys {
		if e, ok := v.(models.MessageEntity); ok {
			out = append(out, e)
		}
	}
	return out
}
