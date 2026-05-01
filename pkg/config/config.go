package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

var (
	TelegramBotToken string
	DatabaseFile     string
	RedisAddr        string
	OwnerID          int64
	SecreteKey       string
	WebAppURL        string
	WebhookURL       string
	AppPort          string
	AppEnv           string
)

func init() {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logger.Warn("CONFIG", "⚠️  .env não encontrado — usando variáveis de ambiente do container")
		}
	}

	TelegramBotToken = mustGetEnv("TELEGRAM_BOT_TOKEN")
	RedisAddr = mustGetEnv("REDIS_HOST")
	DatabaseFile = os.Getenv("DATABASE_FILE") // opcional
	OwnerID = mustGetEnvInt64("OWNER_ID")
	AppPort = os.Getenv("APP_PORT")
	SecreteKey = mustGetEnv("SECRET_KEY")
	WebAppURL = mustGetEnv("WEBAPP_URL")
	WebhookURL = os.Getenv("WEBHOOK_URL") // opcional
	AppEnv = os.Getenv("APP_ENV")         // dev ou prod
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
