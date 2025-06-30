package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func generateSessionKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return "channel_session:" + hex.EncodeToString(bytes), nil
}
