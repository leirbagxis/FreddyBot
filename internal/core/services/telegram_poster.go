package services

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

type TelegramPosterService struct {
	clientSvc  *TelegramClientService
	channelRepo *repositories.ChannelRepository
}

func NewTelegramPosterService(clientSvc *TelegramClientService, channelRepo *repositories.ChannelRepository) *TelegramPosterService {
	return &TelegramPosterService{
		clientSvc:  clientSvc,
		channelRepo: channelRepo,
	}
}

func (s *TelegramPosterService) SendText(ctx context.Context, userID, channelID int64, htmlText string) (int32, error) {
	peer, err := s.resolvePeer(ctx, userID, channelID)
	if err != nil {
		return 0, fmt.Errorf("resolve peer: %w", err)
	}

	client, err := s.clientSvc.ConnectUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("connect user: %w", err)
	}

	plain, entities := HTMLToEntities(htmlText)
	api := client.API()

	result, err := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:     peer,
		Message:  plain,
		Entities: entities,
	})
	if err != nil {
		return 0, fmt.Errorf("send message: %w", err)
	}

	msgID := extractMessageID(result)
	logger.Info("MTPOST", "User %d sent text to channel %d (msgID=%d)", userID, channelID, msgID)
	return msgID, nil
}

func (s *TelegramPosterService) SendWithKeyboard(ctx context.Context, userID, channelID int64, htmlText string, kbRows [][]tg.KeyboardButtonClass) (int32, error) {
	peer, err := s.resolvePeer(ctx, userID, channelID)
	if err != nil {
		return 0, fmt.Errorf("resolve peer: %w", err)
	}

	client, err := s.clientSvc.ConnectUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("connect user: %w", err)
	}

	plain, entities := HTMLToEntities(htmlText)
	var replyMarkup tg.ReplyMarkupClass
	if len(kbRows) > 0 {
		replyMarkup = ConvertKeyboardToMTProto(kbRows)
	}

	api := client.API()
	result, err := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:        peer,
		Message:     plain,
		Entities:    entities,
		ReplyMarkup: replyMarkup,
	})
	if err != nil {
		return 0, fmt.Errorf("send with keyboard: %w", err)
	}

	msgID := extractMessageID(result)
	logger.Info("MTPOST", "User %d sent text+kb to channel %d (msgID=%d)", userID, channelID, msgID)
	return msgID, nil
}

func (s *TelegramPosterService) EditMessage(ctx context.Context, userID, channelID int64, messageID int32, htmlText string, kbRows [][]tg.KeyboardButtonClass) error {
	peer, err := s.resolvePeer(ctx, userID, channelID)
	if err != nil {
		return fmt.Errorf("resolve peer: %w", err)
	}

	client, err := s.clientSvc.ConnectUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("connect user: %w", err)
	}

	plain, entities := HTMLToEntities(htmlText)
	var replyMarkup tg.ReplyMarkupClass
	if len(kbRows) > 0 {
		replyMarkup = ConvertKeyboardToMTProto(kbRows)
	}

	api := client.API()
	_, err = api.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        peer,
		ID:          int(messageID),
		Message:     plain,
		Entities:    entities,
		ReplyMarkup: replyMarkup,
	})
	if err != nil {
		return fmt.Errorf("edit message: %w", err)
	}

	logger.Info("MTPOST", "User %d edited msg %d in channel %d", userID, messageID, channelID)
	return nil
}

func (s *TelegramPosterService) EditText(ctx context.Context, userID, channelID int64, messageID int32, htmlText string) error {
	return s.EditMessage(ctx, userID, channelID, messageID, htmlText, nil)
}

// EditEntitiesOnly updates message entities (custom emoji, bold, italic, etc.)
// without changing the message text. Sends ONLY the Entities field (no Message),
// so the API only updates entities and avoids MESSAGE_NOT_MODIFIED.
func (s *TelegramPosterService) EditEntitiesOnly(ctx context.Context, userID, channelID int64, messageID int32, htmlText string) error {
	peer, err := s.resolvePeer(ctx, userID, channelID)
	if err != nil {
		return fmt.Errorf("resolve peer: %w", err)
	}

	client, err := s.clientSvc.ConnectUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("connect user: %w", err)
	}

	_, entities := HTMLToEntities(htmlText)
	if len(entities) == 0 {
		return nil
	}

	api := client.API()
	_, err = api.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:     peer,
		ID:       int(messageID),
		Entities: entities,
	})
	if err != nil {
		return fmt.Errorf("edit message: %w", err)
	}

	logger.Info("MTPOST", "User %d edited entities of msg %d in channel %d (no text change, %d entities)", userID, messageID, channelID, len(entities))
	return nil
}

func (s *TelegramPosterService) IsChannelAdmin(ctx context.Context, userID int64, channelID int64) (bool, error) {
	ch, err := s.channelRepo.GetChannelByIDLight(ctx, channelID)
	if err != nil {
		return false, nil
	}

	client, err := s.clientSvc.ConnectUser(ctx, userID)
	if err != nil {
		return false, nil
	}

	api := client.API()
	_, err = api.ChannelsGetParticipant(ctx, &tg.ChannelsGetParticipantRequest{
		Channel: &tg.InputChannel{
			ChannelID:  mtprotoChannelID(channelID),
			AccessHash: ch.AccessHash,
		},
		Participant: &tg.InputPeerSelf{},
	})
	if err != nil {
		logger.Info("MTPOST", "User %d is NOT admin of channel %d: %v", userID, channelID, err)
		return false, nil
	}

	logger.Info("MTPOST", "User %d is admin of channel %d", userID, channelID)
	return true, nil
}

func (s *TelegramPosterService) resolvePeer(ctx context.Context, userID, channelID int64) (*tg.InputPeerChannel, error) {
	ch, err := s.channelRepo.GetChannelByIDLight(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	peer := &tg.InputPeerChannel{
		ChannelID:  mtprotoChannelID(channelID),
		AccessHash: ch.AccessHash,
	}

	if ch.AccessHash != 0 {
		return peer, nil
	}

	// Try resolving by username first (most reliable)
	if ch.Username != "" {
		resolved, err := s.clientSvc.ResolvePeerByUsername(ctx, userID, ch.Username)
		if err == nil {
			_ = s.channelRepo.SetAccessHash(ctx, channelID, resolved.AccessHash)
			return resolved, nil
		}
	}

	// Fallback: resolve by channel ID
	accessHash, err := s.clientSvc.ResolveChannelAccessHash(ctx, userID, channelID)
	if err != nil {
		return nil, fmt.Errorf("resolve access hash: %w", err)
	}
	_ = s.channelRepo.SetAccessHash(ctx, channelID, accessHash)
	peer.AccessHash = accessHash
	return peer, nil
}

func extractMessageID(updates tg.UpdatesClass) int32 {
	switch u := updates.(type) {
	case *tg.Updates:
		for _, msg := range u.Updates {
			if m, ok := msg.(*tg.UpdateNewMessage); ok {
				if msg2, ok := m.Message.(*tg.Message); ok {
					return int32(msg2.ID)
				}
			}
			if m, ok := msg.(*tg.UpdateNewChannelMessage); ok {
				if msg2, ok := m.Message.(*tg.Message); ok {
					return int32(msg2.ID)
				}
			}
		}
	case *tg.UpdatesCombined:
		for _, msg := range u.Updates {
			if m, ok := msg.(*tg.UpdateNewMessage); ok {
				if msg2, ok := m.Message.(*tg.Message); ok {
					return int32(msg2.ID)
				}
			}
			if m, ok := msg.(*tg.UpdateNewChannelMessage); ok {
				if msg2, ok := m.Message.(*tg.Message); ok {
					return int32(msg2.ID)
				}
			}
		}
	}
	return 0
}

func ConvertKeyboardToMTProto(kb [][]tg.KeyboardButtonClass) *tg.ReplyInlineMarkup {
	if len(kb) == 0 {
		return nil
	}
	rows := make([]tg.KeyboardButtonRow, len(kb))
	for i, row := range kb {
		rows[i] = tg.KeyboardButtonRow{Buttons: row}
	}
	return &tg.ReplyInlineMarkup{Rows: rows}
}

func ConvertTelegoKeyboardToMTProto(ikb *telego.InlineKeyboardMarkup) [][]tg.KeyboardButtonClass {
	if ikb == nil || len(ikb.InlineKeyboard) == 0 {
		return nil
	}
	rows := make([][]tg.KeyboardButtonClass, len(ikb.InlineKeyboard))
	for i, row := range ikb.InlineKeyboard {
		btns := make([]tg.KeyboardButtonClass, len(row))
		for j, btn := range row {
			switch {
			case btn.URL != "":
				btns[j] = &tg.KeyboardButtonURL{
					Text: btn.Text,
					URL:  btn.URL,
				}
			case btn.CallbackData != "":
				btns[j] = &tg.KeyboardButtonCallback{
					Text: btn.Text,
					Data: []byte(btn.CallbackData),
				}
			default:
				btns[j] = &tg.KeyboardButtonURL{
					Text: btn.Text,
				}
			}
		}
		rows[i] = btns
	}
	return rows
}
