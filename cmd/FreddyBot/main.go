package main

import (
	"log"

	"github.com/leirbagxis/FreddyBot/internal/database"
	"github.com/leirbagxis/FreddyBot/internal/telegram"
)

// Send any text message to the bot after the bot has been started

func main() {
	db := database.InitDB()
	err := telegram.StartBot(db)
	if err != nil {
		log.Fatal(err)
	}
}
