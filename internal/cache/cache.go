package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Service struct{}

func NewService() *Service {
	GetRedisClient()
	return &Service{}
}

func (s *Service) CreateSession(ctx context.Context, payload ChannelPayload) (*Session, error) {
	client := GetRedisClient()

	key, err := generateSessionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session key: %w", err)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = client.Set(ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store in cache: %w", err)
	}

	return &Session{
		Key:     key,
		Payload: payload,
	}, nil
}

func (s *Service) GetSession(ctx context.Context, key string) (*ChannelPayload, error) {
	client := GetRedisClient()

	data, err := client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var payload ChannelPayload
	err = json.Unmarshal([]byte(data), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &payload, nil
}

func (s *Service) DeleteSession(ctx context.Context, key string) error {
	client := GetRedisClient()
	return client.Del(ctx, key).Err()
}

func (s *Service) SetAwaitingStickerSeparator(ctx context.Context, userID, channelID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("awaiting_sticker:%d", userID)
	return client.Set(ctx, key, channelID, 5*time.Minute).Err()
}

func (s *Service) GetAwaitingStickerSeparator(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()

	key := fmt.Sprintf("awaiting_sticker:%d", userID)
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

func (s *Service) DeleteAwaitingStickerSeparator(ctx context.Context, userID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("awaiting_sticker:%d", userID)

	return client.Del(ctx, key).Err()
}

func generateSessionKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
