package telegram

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"gorm.io/gorm"
)

func StartBot(db *gorm.DB) (http.Handler, *telego.Bot) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	cache.GetRedisClient()

	// Inicializar telego
	tb, err := telego.NewBot(config.TelegramBotToken)
	if err != nil {
		panic(err)
	}

	app := container.NewAppContainer(db, tb)

	botInfo, _ := tb.GetMe(context.Background())
	logger.Bot("🤖 Bot iniciado (Telego): %s", botInfo.Username)

	// Updates channel
	updates := make(chan telego.Update, 1000)
	bh, _ := telegohandler.NewBotHandler(tb, updates)

	// Custom HTTP Handler for Webhook
	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("BOT", "❌ Erro ao ler body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var update telego.Update
		if err := json.Unmarshal(body, &update); err != nil {
			logger.Error("BOT", "❌ Erro ao deserealizar update: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		updates <- update
		w.WriteHeader(http.StatusOK)
		logger.Bot("✅ Webhook processado com sucesso")
	})

	go func() {
		<-ctx.Done()
		logger.Bot("🔻 Shutting down gracefully...")
		if err := cache.CloseRedis(); err != nil {
			logger.Error("SYS", "❌ Error closing Redis: %v", err)
		}
		cancel()
	}()

	// Load Handlers
	LoadHandlersTelegoWithBH(bh, app)

	webhookUrl := config.WebhookURL
	if webhookUrl != "" {
		logger.Bot("🔗 Bot configurado para modo webhook: %s", webhookUrl)

		_ = tb.SetWebhook(context.Background(), &telego.SetWebhookParams{
			URL:            webhookUrl,
			AllowedUpdates: []string{"message", "edited_message", "callback_query", "inline_query", "chosen_inline_result", "my_chat_member", "channel_post", "edited_channel_post"},
		})

		logger.Bot("✅ Webhook configurado com sucesso")

		webhookInfo, err := tb.GetWebhookInfo(context.Background())
		if err == nil {
			logger.Bot("📊 Webhook Info - URL: %s, Pending: %d",
				webhookInfo.URL, webhookInfo.PendingUpdateCount)
		}

		logger.Bot("🚀 Iniciando processamento de updates...")
		go bh.Start()

	} else {
		logger.Bot("🔄 Bot iniciado em modo polling")
		_ = tb.DeleteWebhook(context.Background(), &telego.DeleteWebhookParams{})
		
		// Iniciar Long Polling em paralelo para alimentar o channel de updates
		pollingUpdates, _ := tb.UpdatesViaLongPolling(context.Background(), nil)
		go func() {
			for u := range pollingUpdates {
				updates <- u
			}
		}()
		
		go bh.Start()
	}

	return webhookHandler, tb
}
