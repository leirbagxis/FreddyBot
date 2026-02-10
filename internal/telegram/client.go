package telegram

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks"
	"github.com/leirbagxis/FreddyBot/internal/telegram/commands"
	"github.com/leirbagxis/FreddyBot/internal/telegram/events"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/gorm"
)

func StartBot(db *gorm.DB) (http.Handler, bot.Bot) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	opts := []bot.Option{
		bot.WithMiddlewares(
			middleware.SaveUserMiddleware(db),
		),
	}

	cache.GetRedisClient()

	b, err := bot.New(config.TelegramBotToken, opts...)
	if err != nil {
		panic(err)
	}

	app := container.NewAppContainer(db, b)

	botInfo, _ := b.GetMe(ctx)
	botUsername := fmt.Sprintf("@%s", botInfo.Username)
	log.Println("🤖 Bot iniciado:", botInfo.Username)

	originalHandler := b.WebhookHandler()
	debugHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("❌ Erro ao ler body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		originalHandler.ServeHTTP(w, r)
		log.Println("✅ Webhook processado com sucesso")
	})

	go func() {
		<-ctx.Done()
		log.Println("🔻 Shutting down gracefully...")
		if err := cache.CloseRedis(); err != nil {
			log.Printf("❌ Error closing Redis: %v", err)
		}
		cancel()
	}()

	webhookUrl := config.WebhookURL
	if webhookUrl != "" {
		log.Printf("🔗 Bot configurado para modo webhook: %s", webhookUrl)

		events.LoadEvents(b, app)
		commands.LoadCommandHandlers(b, app)
		callbacks.LoadCallbacksHandlers(b, app, botUsername)

		_, err := b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL: webhookUrl,
			//AllowedUpdates: []string{"message", "callback_query", "inline_query", "my_chat_member"},
		})
		if err != nil {
			log.Fatalf("❌ Erro ao setar webhook: %v", err)
		}

		log.Println("✅ Webhook configurado com sucesso")

		webhookInfo, err := b.GetWebhookInfo(ctx)
		if err == nil {
			log.Printf("📊 Webhook Info - URL: %s, Pending: %d",
				webhookInfo.URL, webhookInfo.PendingUpdateCount)
		}

		log.Println("🚀 Iniciando webhook...")
		go b.StartWebhook(ctx)

	} else {
		log.Println("🔄 Bot iniciado em modo polling")

		events.LoadEvents(b, app)
		commands.LoadCommandHandlers(b, app)
		callbacks.LoadCallbacksHandlers(b, app, botUsername)

		b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{})
		go b.Start(ctx)
	}

	return debugHandler, *b
}
