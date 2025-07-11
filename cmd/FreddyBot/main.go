package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/internal/api"
	"github.com/leirbagxis/FreddyBot/internal/database"
	"github.com/leirbagxis/FreddyBot/internal/telegram"
)

// Send any text message to the bot after the bot has been started

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado")
	}

	db := database.InitDB()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Iniciar Bot
	webhookHandler := telegram.StartBot(db)

	// Iniciar API
	go func() {
		if err := api.StartApi(db, webhookHandler); err != nil {
			log.Printf("Erro ao iniciar API: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("🧹 Encerrando app com segurança...")

}
