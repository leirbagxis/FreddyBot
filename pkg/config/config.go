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
	OwnerID          int64
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if TelegramBotToken == "" {
		log.Fatalf(`You need to set the "TELEGRAM_BOT_TOKEN" in the .env file!`)
	}

	DatabaseFile = os.Getenv("DATABASE_FILE")

	OwnerID, _ = strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
	if OwnerID == 0 {
		log.Fatalf(`You need to set the "OWNER_ID" in the .env file!`)
	}
}
