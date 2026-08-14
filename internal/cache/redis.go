package cache

import (
	"context"
	"sync"
	"time"

	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
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
			logger.Error("REDIS", "URL do Redis inválida: %v", err)
			return
		}
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
			logger.Error("REDIS", "Falha ao conectar no Redis: %v", err)
			return
		}

		logger.Bot("✅ Redis conectado com sucesso")

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
