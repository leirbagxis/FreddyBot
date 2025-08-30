package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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

// ### SELECTED CHANNEL ## \\

func (s *Service) SetSelectedChannel(ctx context.Context, userID, channelID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("selected_channel:%d", userID)
	return client.Set(ctx, key, channelID, 43200*time.Minute).Err()
}

func (s *Service) GetSelectedChannel(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()

	key := fmt.Sprintf("selected_channel:%d", userID)
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

// ### SEPARATOR CHANNEL ## \\

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

// ### DELETE CHANNEL ## \\

func (s *Service) SetDeleteChannel(ctx context.Context, userID, channelID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("delete_channel:%d", userID)
	return client.Set(ctx, key, channelID, 5*time.Minute).Err()
}

func (s *Service) GetDeleteChannel(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()

	key := fmt.Sprintf("delete_channel:%d", userID)
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

// ### TRANSFER CHANNEL ## \\

func (s *Service) SetTransferChannel(ctx context.Context, userID, channelID int64) error {
	client := GetRedisClient()

	key := fmt.Sprintf("transfer_channel:%d", userID)
	return client.Set(ctx, key, channelID, 5*time.Minute).Err()
}

func (s *Service) GetTransferChannel(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()

	key := fmt.Sprintf("transfer_channel:%d", userID)
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

// ### DELETE ALL SESSIONS ### \\\
func (s *Service) DeleteAllUserSessionsBySuffix(ctx context.Context, userID int64) (int64, error) {
	client := GetRedisClient()
	pattern := "*:" + strconv.FormatInt(userID, 10)

	var totalDeleted int64
	var cursor uint64
	const page = 1000
	const batch = 500

	for {
		keys, next, err := client.Scan(ctx, cursor, pattern, page).Result()
		if err != nil {
			return totalDeleted, err
		}
		cursor = next

		for i := 0; i < len(keys); i += batch {
			end := i + batch
			if end > len(keys) {
				end = len(keys)
			}
			chunk := keys[i:end]
			n, err := unlinkOrDel(ctx, client, chunk)
			if err != nil {
				return totalDeleted, err
			}
			totalDeleted += n
		}

		if cursor == 0 {
			break
		}
	}

	return totalDeleted, nil
}

func unlinkOrDel(ctx context.Context, client *redis.Client, keys []string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	pipe := client.Pipeline()
	unlink := pipe.Unlink(ctx, keys...)
	_, err := pipe.Exec(ctx)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unknown command") ||
			strings.Contains(strings.ToLower(err.Error()), "unlink") {
			pipe2 := client.Pipeline()
			for _, k := range keys {
				pipe2.Del(ctx, k)
			}
			_, err2 := pipe2.Exec(ctx)
			if err2 != nil {
				return 0, err2
			}
			return int64(len(keys)), nil
		}
		return 0, err
	}
	return unlink.Val(), nil
}

func generateSessionKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
