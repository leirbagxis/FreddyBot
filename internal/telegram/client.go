package telegram

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

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

func StartBot(db *gorm.DB) http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	opts := []bot.Option{
		bot.WithMiddlewares(
			middleware.SaveUserMiddleware(db),
		),
	}

	cache.GetRedisClient()
	app := container.NewAppContainer(db)

	b, err := bot.New(config.TelegramBotToken, opts...)
	if err != nil {
		panic(err)
	}

	botInfo, _ := b.GetMe(ctx)
	log.Println("🤖 Bot iniciado:", botInfo.Username)

	// Criar wrapper do handler com logs detalhados
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

	// Configurar shutdown graceful
	go func() {
		<-ctx.Done()
		log.Println("🔻 Shutting down gracefully...")
		if err := cache.CloseRedis(); err != nil {
			log.Printf("❌ Error closing Redis: %v", err)
		}
		cancel()
	}()

	webhookUrl := config.WebhookUrl
	if webhookUrl != "" {
		log.Printf("🔗 Bot configurado para modo webhook: %s", webhookUrl)

		commands.LoadCommandHandlers(b)
		events.LoadEvents(b, app)
		callbacks.LoadCallbacksHandlers(b, app)

		// 4. QUARTO: Configurar webhook
		_, err := b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL:            webhookUrl,
			AllowedUpdates: []string{"message", "callback_query", "inline_query"},
		})
		if err != nil {
			log.Fatalf("❌ Erro ao setar webhook: %v", err)
		}

		log.Println("✅ Webhook configurado com sucesso")

		// 5. QUINTO: Verificar status do webhook
		webhookInfo, err := b.GetWebhookInfo(ctx)
		if err == nil {
			log.Printf("📊 Webhook Info - URL: %s, Pending: %d",
				webhookInfo.URL, webhookInfo.PendingUpdateCount)
		}

		// 6. POR ÚLTIMO: Iniciar webhook
		log.Println("🚀 Iniciando webhook...")
		go b.StartWebhook(ctx)

	} else {
		log.Println("🔄 Bot iniciado em modo polling")

		// Carregar handlers para modo polling
		commands.LoadCommandHandlers(b)
		events.LoadEvents(b, app)
		callbacks.LoadCallbacksHandlers(b, app)

		b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{})
		go b.Start(ctx)
	}

	return debugHandler
}
