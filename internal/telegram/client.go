package telegram

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands"
	"github.com/leirbagxis/FreddyBot/internal/telegram/events"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/gorm"
)

func StartBot(db *gorm.DB) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(
			middleware.SaveUserMiddleware(db),
		),
	}

	app := container.NewAppContainer(db)

	b, err := bot.New(config.TelegramBotToken, opts...)
	if err != nil {
		panic(err)
	}

	commands.LoadCommandHandlers(b)
	events.LoadEvents(b)
	callbacks.LoadCallbacksHandlers(b, app)

	log.Println("Bot iniciado...")
	b.Start(ctx)
	return nil
}
