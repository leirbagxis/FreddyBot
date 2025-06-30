package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         config.RedisAddr,
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   5,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}

		fmt.Println("✅ Redis connected successfully")

	})
	return redisClient
}

func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

func HealthCheck(ctx context.Context) error {
	client := GetRedisClient()
	return client.Ping(ctx).Err()
}
