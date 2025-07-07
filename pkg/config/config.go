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
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if TelegramBotToken == "" {
		log.Fatalf(`You need to set the "TELEGRAM_BOT_TOKEN" in the .env file!`)
	}

	RedisAddr = os.Getenv("REDIS_HOST")
	if RedisAddr == "" {
		log.Fatalf(`You need to set the "REDIS_HOST" in the .env file!`)
	}

	DatabaseFile = os.Getenv("DATABASE_FILE")

	OwnerID, _ = strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
	if OwnerID == 0 {
		log.Fatalf(`You need to set the "OWNER_ID" in the .env file!`)
	}

	SecreteKey = os.Getenv("SECRETE_KEY")
	if SecreteKey == "" {
		log.Fatalf(`You need to set the "SECRETE_KEY" in the .env file!`)
	}

	WebAppURL = os.Getenv("WEPAPP_URL")
	if WebAppURL == "" {
		log.Fatalf(`You need to set the "WEPAPP_URL" in the .env file!`)
	}
}
