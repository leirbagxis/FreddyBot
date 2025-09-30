package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/routes"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/gorm"
)

func StartApi(db *gorm.DB, webhookHandler http.Handler) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := container.NewAppContainer(db)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(router, app)

	// Rota de Webhook real (usa handler do bot)
	if config.WebhookURL != "" && webhookHandler != nil {
		log.Println("🔗 Registrando endpoint do webhook")
		router.POST("/webhook", gin.WrapH(webhookHandler))
	}

	router.Static("/assets", "./webapp/assets")
	router.GET("/dashboard", func(c *gin.Context) {
		c.File("./webapp/index.html")
	})

	srv := &http.Server{
		Addr:    ":7000",
		Handler: router,
	}

	go func() {
		log.Println("🌐 API REST rodando em http://localhost:7000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("🔻 Encerrando API...")

	return srv.Shutdown(context.Background())
}
