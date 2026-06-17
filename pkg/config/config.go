package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

var (
	TelegramBotToken string
	TelegramAPIID    int
	TelegramAPIHash  string
	DatabaseFile     string
	RedisAddr        string
	OwnerID          int64
	SecreteKey       string
	EncryptionKey    string
	WebAppURL        string
	WebhookURL       string
	AppPort          string
	AppEnv           string
	JWTIssuer        string
	CORSAllowOrigins []string
)

func init() {
	if err := godotenv.Overload(); err != nil {
		logger.Warn("CONFIG", ".env não encontrado — usando variáveis de ambiente do container/sistema")
	}

	TelegramBotToken = mustGetEnv("TELEGRAM_BOT_TOKEN")
	TelegramAPIID = mustGetEnvInt("TELEGRAM_API_ID")
	TelegramAPIHash = mustGetEnv("TELEGRAM_API_HASH")
	RedisAddr = mustGetEnv("REDIS_HOST")
	DatabaseFile = os.Getenv("DATABASE_FILE") // opcional
	OwnerID = mustGetEnvInt64("OWNER_ID")
	AppPort = os.Getenv("APP_PORT")
	SecreteKey = mustGetEnv("SECRET_KEY")
	EncryptionKey = mustGetEnv("ENC_KEY")
	WebAppURL = mustGetEnv("WEBAPP_URL")
	WebhookURL = os.Getenv("WEBHOOK_URL") // opcional
	AppEnv = os.Getenv("APP_ENV")         // dev ou prod
	JWTIssuer = getEnvDefault("JWT_ISSUER", "t.me/legendasbrbot")
	CORSAllowOrigins = parseOrigins(os.Getenv("CORS_ALLOW_ORIGINS"), WebAppURL)

	logger.Info("CONFIG", "WEBAPP_URL=%s", WebAppURL)
}

func mustGetEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return v
}

func mustGetEnvInt(key string) int {
	v := mustGetEnv(key)
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer: %v", key, err)
	}
	return n
}

func mustGetEnvInt64(key string) int64 {
	v := mustGetEnv(key)
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer: %v", key, err)
	}
	return n
}

func getEnvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseOrigins(raw string, fallback string) []string {
	var origins []string
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 && fallback != "" {
		origins = append(origins, fallback)
	}
	return origins
}
