package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/leirbagxis/FreddyBot/internal/api"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database"
	"github.com/leirbagxis/FreddyBot/internal/telegram"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// Send any text message to the bot after the bot has been started

func main() {

	db := database.InitDB()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	webhookHandler, _, app, err := telegram.StartBot(db)
	if err != nil {
		logger.Error("APP", "Erro ao iniciar bot: %v", err)
		return
	}
	defer func() {
		if err := cache.CloseRedis(); err != nil {
			logger.Error("APP", "Erro ao fechar Redis: %v", err)
		}
	}()

	app.StartBackground(ctx)

	go func() {
		if err := api.StartApi(ctx, app, webhookHandler); err != nil {
			logger.Error("APP", "Erro ao iniciar API: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("APP", "🧹 Encerrando app com segurança...")

}
