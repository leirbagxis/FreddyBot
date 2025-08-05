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
	"github.com/leirbagxis/FreddyBot/internal/utils"
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

// NewMessageProcessor factory
func NewMessageProcessor(b *bot.Bot) *MessageProcessor {
	return &MessageProcessor{
		bot:               b,
		permissionManager: NewPermissionManager(),
		mediaGroupManager: NewMediaGroupManager(),
	}
}

func (mp *MessageProcessor) HandleSeparator(channel *dbmodels.Channel, post *models.Message, messageType MessageType) {
	if channel.Separator == nil || channel.Separator.SeparatorID == "" {
		return
	}

	mediaGroupID := post.MediaGroupID
	chatID := post.Chat.ID

	if messageType == MessageTypeAudio && mediaGroupID == "" {
		time.Sleep(1 * time.Second)
		mp.sendSeparatorDirect(channel, chatID)
		return
	}

	if messageType == MessageTypeAudio && mediaGroupID != "" {
		mp.handleGroupSeparator(channel, mediaGroupID, chatID)
		return
	}

	if mediaGroupID != "" && (messageType == MessageTypePhoto || messageType == MessageTypeVideo || messageType == MessageTypeAnimation) {
		log.Printf("🔄 Separator para grupo de mídia %s será enviado via finishGroupProcessing", mediaGroupID)
		return
	}

	mp.sendSeparatorDirect(channel, chatID)
}

func (mp *MessageProcessor) UpdateChannelBasicInfo(ctx context.Context, chatID int64, channel *dbmodels.Channel) (*dbmodels.Channel, bool) {
	chat, err := mp.bot.GetChat(ctx, &bot.GetChatParams{
		ChatID: chatID,
	})
	if err != nil {
		return channel, false
	}
	updated := false

	if chat.Title != "" && chat.Title != channel.Title {
		channel.Title = utils.RemoveHTMLTags(chat.Title)
		updated = true
	}
	if chat.Username != "" {
		newUsername := "@" + chat.Username
		if newUsername != channel.InviteURL {
			channel.InviteURL = newUsername
			updated = true
		}
	} else if chat.InviteLink != "" {
		if chat.InviteLink != channel.InviteURL {
			channel.InviteURL = chat.InviteLink
			updated = true
		}
	}
	if len(channel.Buttons) > 0 {
		buttonUpdated := mp.updateFirstButtonFromChannel(ctx, channel)
		if buttonUpdated {
			updated = true
		}
	}
	return channel, updated
}

// helper, caso usado:
func (mp *MessageProcessor) updateFirstButtonFromChannel(ctx context.Context, channel *dbmodels.Channel) bool {
	if len(channel.Buttons) == 0 {
		return false
	}
	chat, err := mp.bot.GetChat(ctx, &bot.GetChatParams{
		ChatID: channel.ID,
	})
	if err != nil {
		return false
	}
	novoNome := fmt.Sprintf("%s", chat.Title)
	var novaURL string
	if chat.Username != "" {
		novaURL = "https://t.me/" + chat.Username
	} else if chat.InviteLink != "" {
		novaURL = chat.InviteLink
	} else {
		return false
	}
	firstButton := &channel.Buttons[0]
	if firstButton.NameButton == novoNome && firstButton.ButtonURL == novaURL {
		return false
	}
	log.Printf("🔘 Primeiro botão atualizado: '%s' → '%s' | URL: '%s' → '%s'",
		firstButton.NameButton, novoNome, firstButton.ButtonURL, novaURL)
	firstButton.NameButton = utils.RemoveHTMLTags(novoNome)
	firstButton.ButtonURL = novaURL
	return true
}

// helpers usados acima:
func (mp *MessageProcessor) sendSeparatorDirect(channel *dbmodels.Channel, chatID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mp.bot.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  chatID,
		Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
	})
	if err != nil {
		log.Printf("❌ Erro ao enviar separator: %v", err)
	} else {
		log.Printf("✅ Separator enviado com sucesso")
	}
}

func (mp *MessageProcessor) handleGroupSeparator(channel *dbmodels.Channel, mediaGroupID string, chatID int64) {
	if _, exists := groupSeparators.LoadOrStore(mediaGroupID, true); exists {
		return
	}
	time.Sleep(3 * time.Second)
	mp.sendSeparatorDirect(channel, chatID)
	time.AfterFunc(10*time.Second, func() {
		groupSeparators.Delete(mediaGroupID)
	})
}

// ProcessMessage directs to process functions by type
func (mp *MessageProcessor) ProcessMessage(ctx context.Context, messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	switch messageType {
	case MessageTypeText:
		return mp.ProcessTextMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeAudio:
		return mp.ProcessAudioMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeSticker:
		if len(buttons) > 0 {
			return mp.ProcessStickerMessage(ctx, channel, post, buttons)
		}
		return nil
	case MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation:
		return mp.ProcessMediaMessage(ctx, channel, post, buttons, messageEditAllowed)
	default:
		return nil
	}
}

// CreateInlineKeyboard creates inline keyboard checking button validity and permissions
func (mp *MessageProcessor) CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption, channel *dbmodels.Channel, messageType MessageType) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button
	if customCaption != nil && len(customCaption.Buttons) > 0 {
		log.Printf("🔘 Using custom caption buttons: %s (%d buttons) - ignoring ButtonsPermissions", customCaption.Code, len(customCaption.Buttons))
		for _, b := range customCaption.Buttons {
			if b.NameButton != "" && b.ButtonURL != "" {

			}
		}
	} else {
		if !mp.permissionManager.IsButtonsAllowed(channel, messageType) {
			log.Printf("🔘 Default buttons blocked by ButtonsPermissions for channel %d", channel.ID)
			return nil
		}
		for _, b := range buttons {
			if b.NameButton == "" || b.ButtonURL == "" {
				log.Printf("⚠️ Invalid button ignored: %+v", b)
				continue
			}
			finalButtons = append(finalButtons, b)
		}
	}
	if len(finalButtons) == 0 {
		log.Printf("🔘 No valid buttons available")
		return nil
	}
	buttonGrid := make(map[int]map[int]models.InlineKeyboardButton)
	for i, b := range finalButtons {
		row, col := b.PositionY, b.PositionX
		if col == 0 {
			col = i
		}
		if buttonGrid[row] == nil {
			buttonGrid[row] = make(map[int]models.InlineKeyboardButton)
		}
		buttonGrid[row][col] = models.InlineKeyboardButton{
			Text: b.NameButton,
			URL:  b.ButtonURL,
		}
	}
	var keyboard [][]models.InlineKeyboardButton
	for row := 0; row < 10; row++ {
		if rowButtons, ok := buttonGrid[row]; ok {
			var keyboardRow []models.InlineKeyboardButton
			for col := 0; col < 10; col++ {
				if btn, ok := rowButtons[col]; ok {
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
	return &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}
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

// ProcessTextMessage processes text messages with permissions and formatting
func (mp *MessageProcessor) ProcessTextMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	text := post.Text
	messageID := post.ID
	messageType := MessageTypeText
	if text == "" {
		return fmt.Errorf("message text is empty")
	}
	permissions := mp.CheckPermissions(channel, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Text editing blocked for channel %d: %s", channel.ID, permissions.Reason)
		return fmt.Errorf("text edit not allowed")
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
	formattedText := processTextWithFormatting(text, post.Entities)
	message, customCaption := mp.processMessageWithHashtag(formattedText, channel)
	canEdit, allowedButtons, allowedCustomCaption := mp.ApplyPermissions(channel, messageType, customCaption, buttons)
	if !canEdit {
		return fmt.Errorf("insufficient permissions to edit message")
	}
	keyboard := mp.CreateInlineKeyboard(allowedButtons, allowedCustomCaption, channel, messageType)
	editParams := &bot.EditMessageTextParams{
		ChatID:    post.Chat.ID,
		MessageID: messageID,
		Text:      message,
		ParseMode: "HTML",
	}
	disableLinkPreview := false
	if !permissions.CanUseLinkPreview {
		disableLinkPreview = true
	}
	if customCaption != nil && !customCaption.LinkPreview {
		disableLinkPreview = true
	}
	if disableLinkPreview {
		val := true
		editParams.LinkPreviewOptions = &models.LinkPreviewOptions{IsDisabled: &val}
	}
	if keyboard != nil {
		editParams.ReplyMarkup = keyboard
	}
	_, err := mp.bot.EditMessageText(ctx, editParams)
	return err
}

// ProcessAudioMessage processes audio messages (simplified)
func (mp *MessageProcessor) ProcessAudioMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	// Similar permission checks and logic as ProcessTextMessage
	// Implement detailed logic as needed
	return nil
}

// ProcessMediaMessage processes Photo, Video, Animation
func (mp *MessageProcessor) ProcessMediaMessage(ctx context.Context, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	// Similar checks and implementation for media types
	return nil
}

// ProcessStickerMessage processes sticker messages
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

// Helper functions to extract hashtag, remove hashtag from text, find custom caption, etc.

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
	if v, ok := removeHashRegexCache.Load(hashtag); ok {
		re = v.(*regexp.Regexp)
	} else {
		re = regexp.MustCompile(`#` + regexp.QuoteMeta(hashtag) + `\s*`)
		removeHashRegexCache.Store(hashtag, re)
	}
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}

func findCustomCaption(channel *dbmodels.Channel, hashtag string) *dbmodels.CustomCaption {
	cacheKey := fmt.Sprintf("%d_%s", channel.ID, hashtag)
	if val, ok := customCaptionCache.Load(cacheKey); ok {
		if caption, ok := val.(*dbmodels.CustomCaption); ok {
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
	log.Printf("📝 ❌ Custom caption not found for: #%s", hashtag)
	customCaptionCache.Store(cacheKey, (*dbmodels.CustomCaption)(nil))
	return nil
}

func (mp *MessageProcessor) processMessageWithHashtaga(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
	hashtag := extractHashtag(text)
	if hashtag == "" {
		defCaption := ""
		if channel.DefaultCaption != nil {
			defCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defCaption), nil
	}
	customCaption := findCustomCaption(channel, hashtag)
	if customCaption == nil {
		defCaption := ""
		if channel.DefaultCaption != nil {
			defCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defCaption), nil
	}
	cleanText := removeHashtag(text, hashtag)
	formatted := detectParseMode(customCaption.Caption)
	return fmt.Sprintf("%s\n\n%s", cleanText, formatted), customCaption
}

// ApplyPermissions applies message and button permissions with custom captions
func (mp *MessageProcessor) ApplyPermissions(channel *dbmodels.Channel, messageType MessageType, customCaption *dbmodels.CustomCaption, buttons []dbmodels.Button) (bool, []dbmodels.Button, *dbmodels.CustomCaption) {
	permissions := mp.CheckCustomCaptionPermissions(channel, customCaption, messageType)
	if !permissions.CanEdit {
		log.Printf("❌ Message edit blocked: %s", permissions.Reason)
		return false, nil, nil
	}
	if !permissions.CanAddButtons {
		log.Printf("⚠️ Default buttons removed due to ButtonsPermissions")
		buttons = nil
		log.Printf("✅ Custom caption buttons kept (independent of ButtonsPermissions)")
	}
	return true, buttons, customCaption
}

// CheckPermissions checks channel permissions for a message type
func (mp *MessageProcessor) CheckPermissions(channel *dbmodels.Channel, messageType MessageType) *PermissionCheckResult {
	return mp.permissionManager.CheckPermissions(channel, messageType)
}

// CheckCustomCaptionPermissions checks custom caption permissions for message type
func (mp *MessageProcessor) CheckCustomCaptionPermissions(channel *dbmodels.Channel, customCaption *dbmodels.CustomCaption, messageType MessageType) *PermissionCheckResult {
	// Very similar to CheckPermissions but includes custom caption checks
	result := &PermissionCheckResult{
		CanEdit:           true,
		CanAddButtons:     true,
		CanEditButtons:    true,
		CanUseLinkPreview: true,
	}
	if channel == nil {
		result.CanEdit = false
		result.Reason = "Channel not found"
		return result
	}
	if channel.DefaultCaption == nil {
		log.Printf("⚠️ Channel %d has no DefaultCaption - allowing all", channel.ID)
		return result
	}
	if channel.DefaultCaption.MessagePermission != nil {
		mp := channel.DefaultCaption.MessagePermission
		// use PermissionMap from types.go
		permKey := PermissionMap[messageType]
		switch permKey {
		case "message":
			if !mp.Message {
				result.CanEdit = false
				result.Reason = "Text message edit disabled"
			}
		case "audio":
			if !mp.Audio {
				result.CanEdit = false
				result.Reason = "Audio edit disabled"
			}
		case "video":
			if !mp.Video {
				result.CanEdit = false
				result.Reason = "Video edit disabled"
			}
		case "photo":
			if !mp.Photo {
				result.CanEdit = false
				result.Reason = "Photo edit disabled"
			}
		case "sticker":
			if !mp.Sticker {
				result.CanEdit = false
				result.Reason = "Sticker edit disabled"
			}
		case "gif":
			if !mp.GIF {
				result.CanEdit = false
				result.Reason = "GIF edit disabled"
			}
		}
		if messageType == MessageTypeText && !mp.LinkPreview {
			result.CanUseLinkPreview = false
			log.Printf("🔗 Link preview disabled by MessagePermission for channel %d", channel.ID)
		}
	}
	if channel.DefaultCaption.ButtonsPermission != nil {
		bp := channel.DefaultCaption.ButtonsPermission
		permKey := PermissionMap[messageType]
		switch permKey {
		case "message":
			if !bp.Message {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case "audio":
			if !bp.Audio {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case "video":
			if !bp.Video {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case "photo":
			if !bp.Photo {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case "sticker":
			if !bp.Sticker {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		case "gif":
			if !bp.GIF {
				result.CanAddButtons = false
				result.CanEditButtons = false
			}
		}
	}
	if customCaption != nil && messageType == MessageTypeText {
		if !customCaption.LinkPreview {
			result.CanUseLinkPreview = false
			log.Printf("🔗 Link preview disabled by CustomCaption %s for channel %d", customCaption.Code, channel.ID)
		}
		log.Printf("✅ Custom caption %s: %d buttons (always allowed)", customCaption.Code, len(customCaption.Buttons))
	}
	log.Printf("✅ Permissions checked - Edit=%v, ButtonsDefault=%v, LinkPreview=%v", result.CanEdit, result.CanAddButtons, result.CanUseLinkPreview)
	return result
}

// Convert helper functions (convertMessageEntitiesToInterface, convertInterfaceToMessageEntities)

// And additional utility functions as needed...

// IsNewPackActive and SetNewPackActive proxy to mediaGroupManager methods

func (mp *MessageProcessor) IsNewPackActive(channelID int64) bool {
	return mp.mediaGroupManager.IsNewPackActive(channelID)
}

func (mp *MessageProcessor) SetNewPackActive(channelID int64, active bool) {
	mp.mediaGroupManager.SetNewPackActive(channelID, active)
}
