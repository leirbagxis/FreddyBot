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

// ✅ Áudio: individual edita no lugar; grupo reenvia + apaga; sem permissão tenta só botões
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, _ bool) error {
	perms := mp.CheckPermissions(channel, MessageTypeAudio)

	// Grupo de áudio? agrupar e finalizar depois
	if post.MediaGroupID != "" {
		return mp.handleGroupedAudio(ctx, channel, post, buttons)
	}

	// Áudio único
	if !perms.CanEdit {
		if len(buttons) == 0 || !perms.CanAddButtons {
			return nil
		}
		kb := mp.CreateInlineKeyboard(buttons, nil, channel, MessageTypeAudio)
		if kb == nil {
			return nil
		}
		ctxEdit, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		_, err := mp.bot.EditMessageReplyMarkup(ctxEdit, &bot.EditMessageReplyMarkupParams{
			ChatID:      post.Chat.ID,
			MessageID:   post.ID,
			ReplyMarkup: kb,
		})
		return err
	}

	// Pode editar: substituir legenda por banco se existir, senão manter original
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
	finalCaption := formatted
	if dbCaption != "" {
		finalCaption = dbCaption
	}

	_, filteredButtons, allowedCustom := mp.ApplyPermissions(channel, MessageTypeAudio, customCaption, buttons)
	kb := mp.CreateInlineKeyboard(filteredButtons, allowedCustom, channel, MessageTypeAudio)

	edit := &bot.EditMessageCaptionParams{
		ChatID:    post.Chat.ID,
		MessageID: post.ID,
		Caption:   finalCaption,
		ParseMode: "HTML",
	}
	if kb != nil {
		edit.ReplyMarkup = kb
	}

	ctxEdit, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(250+attempt*250) * time.Millisecond)
		}
		_, lastErr = mp.bot.EditMessageCaption(ctxEdit, edit)
		if lastErr == nil {
			break
		}
		if strings.Contains(strings.ToLower(lastErr.Error()), "too many requests") {
			time.Sleep(1 * time.Second)
			continue
		}
	}
	return lastErr
}

// ✅ Acumula itens do grupo de áudio e arma timer para finalizar
func (mp *MessageProcessor) handleGroupedAudio(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button) error {
	mediaGroupID := post.MediaGroupID
	if mediaGroupID == "" || post.Audio == nil {
		return nil
	}

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
		FileID:          post.Audio.FileID,
		HasCaption:      post.Caption != "",
		Caption:         post.Caption,
		CaptionEntities: convertMessageEntitiesToInterface(post.CaptionEntities),
	})
	// Reset/arma timer para consolidar após última peça chegar
	if group.Timer != nil {
		group.Timer.Stop()
	}
	timeout := time.Duration(900+len(group.Messages)*200) * time.Millisecond
	if timeout > 2500*time.Millisecond {
		timeout = 2500 * time.Millisecond
	}
	group.Timer = time.AfterFunc(timeout, func() {
		mp.finishGroupedAudioProcessing(channel, mediaGroupID, buttons)
	})
	group.mu.Unlock()

	return nil
}

// Finaliza álbum de áudio:
// - Se ButtonsPermissions.Audio == false e MessagePermissions.Audio == false: não faz nada.
// - Caso contrário (mesmo com MessagePermissions.Audio == false):
//   - Apaga cada mensagem original do grupo
//   - Reenvia cada áudio com a LEGENDA ORIGINAL do item
//   - Se ButtonsPermissions.Audio == true, inclui os botões
//   - Ao final, envia o separator (único envio)

func (mp *MessageProcessor) finishGroupedAudioProcessingg(channel *dbmodels.Channel, groupID string, buttons []dbmodels.Button) {
	value, ok := mediaGroups.Load(groupID)
	if !ok {
		return
	}
	group := value.(*MediaGroup)

	// Evita processamento duplicado
	group.mu.Lock()
	if group.Processed {
		group.mu.Unlock()
		return
	}
	group.Processed = true

	// Snapshot e libera lock
	messages := append([]MediaMessage(nil), group.Messages...)
	chatID := group.ChatID
	group.mu.Unlock()

	// Verificar permissões específicas (edicao/botoes)
	perms := mp.CheckPermissions(channel, MessageTypeAudio)

	// Se não pode nem editar nem adicionar botões, não faz nada
	if !perms.CanEdit && !perms.CanAddButtons {
		mediaGroups.Delete(groupID)
		return
	}

	// 1) Determinar baseCaption do grupo para resolver Custom/Default
	var baseCaption string
	var baseEntities []interface{}
	for _, m := range messages {
		if m.HasCaption {
			baseCaption = m.Caption
			baseEntities = m.CaptionEntities
			break
		}
	}
	if baseCaption == "" && len(messages) > 0 {
		baseCaption = messages[0].Caption
		baseEntities = messages[0].CaptionEntities
	}

	// Formatar original e extrair hashtag
	formattedBase := processTextWithFormatting(baseCaption, convertInterfaceToEntities(baseEntities))

	// 2) Resolver dbCaption (Custom > Default), que será aplicada a TODO o álbum
	var dbCaption string
	var customCaption *dbmodels.CustomCaption
	if h := extractHashtag(formattedBase); h != "" {
		if cc := findCustomCaption(channel, h); cc != nil {
			customCaption = cc
			dbCaption = detectParseMode(cc.Caption)
		}
	}
	if customCaption == nil && channel.DefaultCaption != nil {
		dbCaption = detectParseMode(channel.DefaultCaption.Caption)
	}
	// Observação: dbCaption pode ser "", o que implica APAGAR a legenda no reenvio

	// 3) Preparar teclado se ButtonsPermissions.Audio permitir
	var kb *models.InlineKeyboardMarkup
	if perms.CanAddButtons {
		_, filteredButtons, allowedCustom := mp.ApplyPermissions(channel, MessageTypeAudio, customCaption, buttons)
		kb = mp.CreateInlineKeyboard(filteredButtons, allowedCustom, channel, MessageTypeAudio)
		if kb != nil && len(kb.InlineKeyboard) == 0 {
			kb = nil
		}
	}

	// 4) Reenvio item-a-item usando SEMPRE dbCaption (vazia => apaga)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i, m := range messages {
		send := &bot.SendAudioParams{
			ChatID:    chatID,
			Audio:     &models.InputFileString{Data: m.FileID},
			Caption:   dbCaption, // substitui SEMPRE pela legenda do banco
			ParseMode: "HTML",
		}
		if kb != nil {
			send.ReplyMarkup = kb
		}

		// Backoff simples entre itens do álbum
		time.Sleep(time.Duration(200+i*150) * time.Millisecond)

		var sendErr error
		for attempt := 0; attempt < 3; attempt++ {
			if attempt > 0 {
				time.Sleep(time.Duration(280+attempt*320) * time.Millisecond)
			}
			_, sendErr = mp.bot.SendAudio(ctx, send)
			if sendErr == nil {
				break
			}
			if strings.Contains(strings.ToLower(sendErr.Error()), "too many requests") {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		if sendErr != nil {
			log.Printf("❌ Falha ao reenviar áudio do grupo %s (msg %d): %v", groupID, m.MessageID, sendErr)
			// Não apaga o original se falhou o reenvio
			continue
		}

		// Apagar original após reenviar
		time.Sleep(200 * time.Millisecond)
		_, delErr := mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: m.MessageID,
		})
		if delErr != nil {
			log.Printf("⚠️ Falha ao apagar áudio original (grupo %s msg %d): %v", groupID, m.MessageID, delErr)
		}
	}

	// 5) Enviar separator ao final
	if channel.Separator != nil && channel.Separator.SeparatorID != "" {
		time.Sleep(350 * time.Millisecond)
		sepCtx, cancelSep := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelSep()
		_, err := mp.bot.SendSticker(sepCtx, &bot.SendStickerParams{
			ChatID:  chatID,
			Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
		})
		if err != nil {
			log.Printf("⚠️ Falha ao enviar separator pós-álbum %s: %v", groupID, err)
		} else {
			log.Printf("✅ Separator enviado após processamento do grupo %s", groupID)
		}
	}

	mediaGroups.Delete(groupID)
}

func (mp *MessageProcessor) finishGroupedAudioProcessing(channel *dbmodels.Channel, groupID string, buttons []dbmodels.Button) {
	value, ok := mediaGroups.Load(groupID)
	if !ok {
		return
	}
	group := value.(*MediaGroup)

	// Evitar processamento duplicado
	group.mu.Lock()
	if group.Processed {
		group.mu.Unlock()
		return
	}
	group.Processed = true

	// Snapshot e libera lock
	messages := append([]MediaMessage(nil), group.Messages...)
	chatID := group.ChatID
	group.mu.Unlock()

	// Verificar permissões específicas
	perms := mp.CheckPermissions(channel, MessageTypeAudio)

	// Se não pode nem editar nem adicionar botões, não faz nada
	if !perms.CanEdit && !perms.CanAddButtons {
		mediaGroups.Delete(groupID)
		return
	}

	// 1) Determinar baseCaption do grupo (para resolver Custom/Default se edição estiver permitida)
	var baseCaption string
	var baseEntities []interface{}
	for _, m := range messages {
		if m.HasCaption {
			baseCaption = m.Caption
			baseEntities = m.CaptionEntities
			break
		}
	}
	if baseCaption == "" && len(messages) > 0 {
		baseCaption = messages[0].Caption
		baseEntities = messages[0].CaptionEntities
	}

	// 2) Resolver dbCaption (apenas se edição for permitida)
	var dbCaption string
	var customCaption *dbmodels.CustomCaption
	if perms.CanEdit {
		formattedBase := processTextWithFormatting(baseCaption, convertInterfaceToEntities(baseEntities))
		if h := extractHashtag(formattedBase); h != "" {
			if cc := findCustomCaption(channel, h); cc != nil {
				customCaption = cc
				dbCaption = detectParseMode(cc.Caption)
			}
		}
		if customCaption == nil && channel.DefaultCaption != nil {
			dbCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		// Observação: dbCaption pode ser "", o que implica APAGAR a legenda no reenvio
	}

	// 3) Preparar teclado se ButtonsPermissions.Audio permitir
	var kb *models.InlineKeyboardMarkup
	if perms.CanAddButtons {
		_, filteredButtons, allowedCustom := mp.ApplyPermissions(channel, MessageTypeAudio, customCaption, buttons)
		kb = mp.CreateInlineKeyboard(filteredButtons, allowedCustom, channel, MessageTypeAudio)
		if kb != nil && len(kb.InlineKeyboard) == 0 {
			kb = nil
		}
	}

	// 4) Reenvio item-a-item
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i, m := range messages {
		// Aplicar legenda somente se edição for permitida; caso contrário, enviar sem caption
		cap := ""
		if perms.CanEdit {
			cap = dbCaption // pode ser vazia para apagar
		}

		send := &bot.SendAudioParams{
			ChatID:    chatID,
			Audio:     &models.InputFileString{Data: m.FileID},
			Caption:   cap,
			ParseMode: "HTML",
		}
		if kb != nil {
			send.ReplyMarkup = kb
		}

		// Backoff simples entre itens
		time.Sleep(time.Duration(200+i*150) * time.Millisecond)

		var sendErr error
		for attempt := 0; attempt < 3; attempt++ {
			if attempt > 0 {
				time.Sleep(time.Duration(280+attempt*320) * time.Millisecond)
			}
			_, sendErr = mp.bot.SendAudio(ctx, send)
			if sendErr == nil {
				break
			}
			if strings.Contains(strings.ToLower(sendErr.Error()), "too many requests") {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		if sendErr != nil {
			log.Printf("❌ Falha ao reenviar áudio do grupo %s (msg %d): %v", groupID, m.MessageID, sendErr)
			// Não apaga o original se falhou o reenvio
			continue
		}

		// Apagar original após reenviar
		time.Sleep(200 * time.Millisecond)
		_, delErr := mp.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: m.MessageID,
		})
		if delErr != nil {
			log.Printf("⚠️ Falha ao apagar áudio original (grupo %s msg %d): %v", groupID, m.MessageID, delErr)
		}
	}

	// 5) Enviar separator ao final
	if channel.Separator != nil && channel.Separator.SeparatorID != "" {
		time.Sleep(350 * time.Millisecond)
		sepCtx, cancelSep := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelSep()
		_, err := mp.bot.SendSticker(sepCtx, &bot.SendStickerParams{
			ChatID:  chatID,
			Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
		})
		if err != nil {
			log.Printf("⚠️ Falha ao enviar separator pós-álbum %s: %v", groupID, err)
		} else {
			log.Printf("✅ Separator enviado após processamento do grupo %s", groupID)
		}
	}

	mediaGroups.Delete(groupID)
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

// Envia o separador independentemente de CanEdit/CanAddButtons; suprime apenas no início de álbum de áudio.
// Use esta função no Handler após enfileirar a mensagem; o finalizador de grupo enviará no fim do álbum.
func (mp *MessageProcessor) ProcessSeparator(ctx context.Context, channel *dbmodels.Channel, post *models.Message) error {
	if channel == nil || channel.Separator == nil || channel.Separator.SeparatorID == "" {
		log.Printf("⚠️ Separator não configurado para o canal")
		return nil
	}

	// Suprime no início do álbum de áudio: deixa para o finalizador do grupo
	if post != nil && post.MediaGroupID != "" && post.Audio != nil {
		log.Printf("ℹ️ Separator suprimido no início do álbum de áudio (groupID=%s)", post.MediaGroupID)
		return nil
	}

	// Determinar chat alvo
	var chatID int64
	if post != nil {
		chatID = post.Chat.ID
	} else {
		chatID = channel.ID
	}

	// Contexto próprio com timeout
	sendCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("🔄 Enviando separator para chat %d", chatID)

	maxRetries := 2
	baseDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := mp.bot.SendSticker(sendCtx, &bot.SendStickerParams{
			ChatID:  chatID,
			Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
		})
		if err == nil {
			log.Printf("✅ Separator enviado com sucesso para chat %d", chatID)
			return nil
		}

		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "too many requests") || strings.Contains(lower, "429") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter <= 0 {
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
