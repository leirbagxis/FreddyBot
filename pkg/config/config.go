package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// isTestMode retorna true se o binario foi compilado com testes.
func isTestMode() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

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
	JWTIssuer        string
	CORSAllowOrigins []string

	// MTProto
	MTProtoAppID    int
	MTProtoAppHash  string
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
	JWTIssuer = getEnvDefault("JWT_ISSUER", "t.me/legendasbrbot")
	CORSAllowOrigins = parseOrigins(os.Getenv("CORS_ALLOW_ORIGINS"), WebAppURL)

	// MTProto credentials (optional for now, required when MTProto client is active)
	mtprotoAppIDStr := os.Getenv("MTPROTO_APP_ID")
	if mtprotoAppIDStr != "" {
		if id, err := strconv.Atoi(mtprotoAppIDStr); err == nil {
			MTProtoAppID = id
		}
	}
	MTProtoAppHash = os.Getenv("MTPROTO_APP_HASH")
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		if isTestMode() {
			return "test_" + key
		}
		log.Fatalf("Environment variable %s is required", key)
	}
	return v
}

func mustGetEnvInt64(key string) int64 {
	v := mustGetEnv(key)
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		if isTestMode() {
			return 0
		}
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

// GetMTProtoAppID retorna o App ID para MTProto ou 0 se nao configurado.
func GetMTProtoAppID() int {
	return MTProtoAppID
}

// GetMTProtoAppHash retorna o App Hash para MTProto ou string vazia.
func GetMTProtoAppHash() string {
	return MTProtoAppHash
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
