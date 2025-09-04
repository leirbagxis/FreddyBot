package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type SessionManager struct {
	cache *Service
}

func NewSessionManager(cache *Service) *SessionManager {
	return &SessionManager{
		cache: cache,
	}
}

func (sm *SessionManager) CreateChannelSession(ctx context.Context, channelID, ownerID int64, title string) (*Session, error) {
	payload := ChannelPayload{
		ChannelID: channelID,
		Title:     title,
		OwnerID:   ownerID,
	}

	return sm.cache.CreateSession(ctx, payload)
}

func (sm *SessionManager) CreateClaimSession(ctx context.Context, channelID, ownerID, newOwnerId int64) (*Session, error) {
	payload := ChannelPayload{
		ChannelID:  channelID,
		OwnerID:    ownerID,
		NewOwnerID: newOwnerId,
	}

	return sm.cache.CreateSession(ctx, payload)
}

func (sm *SessionManager) GetChannelSession(ctx context.Context, key string) (*ChannelPayload, error) {
	return sm.cache.GetSession(ctx, key)
}

func (sm *SessionManager) DeleteChannelSession(ctx context.Context, key string) error {
	return sm.cache.DeleteSession(ctx, key)
}

// ### PLAN CAPTION SESSION ## \\

func (sm *SessionManager) SetPlainCaptionSession(ctx context.Context, userID, channelID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("ask_plain_caption:%d", userID)
	return client.Set(ctx, key, channelID, 5*time.Minute).Err()
}

func (sm *SessionManager) GetPlainCaptionSession(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()

	key := fmt.Sprintf("ask_plain_caption:%d", userID)
	data, err := client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, fmt.Errorf("session not found or expired")
		}
		return 0, fmt.Errorf("failed to get from cache: %w", err)
	}

	channelID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0, err
	}

	return channelID, nil
}

func (sm *SessionManager) DeletePlainCaptionSession(ctx context.Context, userID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("ask_plain_caption:%d", userID)

	return client.Del(ctx, key).Err()
}
