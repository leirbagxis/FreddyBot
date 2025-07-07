package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/internal/api"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/database"
	"github.com/leirbagxis/FreddyBot/internal/telegram"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

// Send any text message to the bot after the bot has been started

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado")
	}

	token := auth.GenerateSignature("7595607953", "-1002676384505", config.SecreteKey)
	fmt.Println(token)

	db := database.InitDB()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Iniciar Bot
	go func() {
		if err := telegram.StartBot(db); err != nil {
			log.Printf("Erro ao iniciar Bot: %v", err)
			stop()
		}
	}()

	// Iniciar API
	go func() {
		if err := api.StartApi(db); err != nil {
			log.Printf("Erro ao iniciar API: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("🧹 Encerrando app com segurança...")

}
