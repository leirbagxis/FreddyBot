package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands"
	"github.com/leirbagxis/FreddyBot/internal/telegram/events"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/gorm"
)

func CreateBot(db *gorm.DB) (*bot.Bot, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opts := []bot.Option{
		bot.WithMiddlewares(
			middleware.SaveUserMiddleware(db),
		),
	}

	cache.GetRedisClient()
	app := container.NewAppContainer(db)

	b, err := bot.New(config.TelegramBotToken, opts...)
	if err != nil {
		return nil, err
	}

	commands.LoadCommandHandlers(b)
	events.LoadEvents(b, app)
	callbacks.LoadCallbacksHandlers(b, app)

	botInfo, _ := b.GetMe(context.Background())
	log.Println("Bot criado...", botInfo.Username)

	return b, nil
}

func SetupWebhook(b *bot.Bot, webhookURL string) error {
	ctx := context.Background()

	// Configurar webhook
	_, err := b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: webhookURL,
	})
	if err != nil {
		return err
	}

	log.Printf("Webhook configurado com sucesso: %s", webhookURL)
	return nil
}

func CleanupWebhook(b *bot.Bot) error {
	ctx := context.Background()

	// Remover webhook
	_, err := b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{})
	if err != nil {
		log.Printf("Erro ao remover webhook: %v", err)
		return err
	}

	log.Println("Webhook removido com sucesso")
	return nil
}
