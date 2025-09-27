package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	TelegramBotToken string
	DatabaseFile     string
	RedisAddr        string
	OwnerID          int64
	SecreteKey       string
	WebAppURL        string
	WebhookUrl       string
)

func init() {
	// Carrega .env apenas se existir, para desenvolvimento local.
	_ = godotenv.Load()

	TelegramBotToken = mustGetEnv("TELEGRAM_BOT_TOKEN")
	RedisAddr        = mustGetEnv("REDIS_HOST")
	DatabaseFile     = os.Getenv("DATABASE_FILE") // pode ser opcional
	OwnerID          = mustGetEnvInt64("OWNER_ID")
	SecreteKey       = mustGetEnv("SECRET_KEY")
	WebAppURL        = mustGetEnv("WEBAPP_URL")   // Corrigido o nome da var
	WebhookUrl       = os.Getenv("WEBHOOK_URL")   // opcional
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return v
}

func mustGetEnvInt64(key string) int64 {
	v := mustGetEnv(key)
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer: %v", key, err)
	}
	return n
}
