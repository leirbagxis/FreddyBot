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
		var err error

		opt, err := redis.ParseURL(config.RedisAddr)
		if err != nil {
			panic(fmt.Sprintf("Invalid Redis URL: %v", err))
		}

		// ⚙️ sobrescreve configurações adicionais
		opt.PoolSize = 10
		opt.MinIdleConns = 5
		opt.MaxRetries = 5
		opt.DialTimeout = 5 * time.Second
		opt.ReadTimeout = 5 * time.Second
		opt.WriteTimeout = 3 * time.Second

		redisClient = redis.NewClient(opt)

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
