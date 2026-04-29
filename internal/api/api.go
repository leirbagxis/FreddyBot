package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/leirbagxis/FreddyBot/internal/api/routes"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"gorm.io/gorm"
)

func StartApi(db *gorm.DB, webhookHandler http.Handler, bot *bot.Bot) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := container.NewAppContainer(db, bot)
	router := gin.Default() // Usar Default para ter Logger e Recovery

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "x-telegram-init-data"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(router, app)

	// Rota de Webhook real (usa handler do bot)
	if config.WebhookURL != "" && webhookHandler != nil {
		logger.API("🔗 Registrando endpoint do webhook")
		router.POST("/webhook", gin.WrapH(webhookHandler))
	}

	// 1. Arquivos estáticos PRECISAM vir antes das rotas dinâmicas
	router.Static("/dashboard/assets", "./dashboard/dist/assets")
	router.StaticFile("/dashboard/favicon.svg", "./dashboard/dist/favicon.svg")
	router.StaticFile("/dashboard/icons.svg", "./dashboard/dist/icons.svg")
	router.StaticFile("/favicon.svg", "./dashboard/dist/favicon.svg")
	router.StaticFile("/icons.svg", "./dashboard/dist/icons.svg")

	// 2. Handler para o Dashboard
	dashboardHandler := func(c *gin.Context) {
		// Proteção contra caminhos curtos para evitar panic no slice
		path := c.Request.URL.Path
		const assetsPrefix = "/dashboard/assets/"
		
		if len(path) >= len(assetsPrefix) && path[:len(assetsPrefix)] == assetsPrefix {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.File("./dashboard/dist/index.html")
	}

	// 3. Rotas de Dashboard (mais específicas)
	router.GET("/dashboard/:channelID", dashboardHandler)
	router.GET("/dashboard/:channelID/captions", dashboardHandler)
	router.GET("/dashboard/:channelID/buttons", dashboardHandler)
	router.GET("/dashboard/:channelID/admin", dashboardHandler)
	router.GET("/admin/dash", dashboardHandler)
	router.GET("/me/channels", dashboardHandler)
	
	// Fallback para qualquer rota de dashboard não mapeada
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		logger.API("⚠️ Rota não encontrada: %s", path)
		
		// Se começar com /api ou /webhook e não foi encontrada, é 404 real
		if (len(path) >= 4 && path[:4] == "/api") || (len(path) >= 8 && path[:8] == "/webhook") {
			c.JSON(404, gin.H{"code": "ENDPOINT_NOT_FOUND", "message": "Endpoint de API não encontrado"})
			return
		}

		// Para todo o resto, servimos o Dashboard (React handles the routes)
		dashboardHandler(c)
	})

	port := utils.NormalizePort(config.AppPort)

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.API("🌐 API REST rodando em http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API", "Erro ao iniciar servidor: %v", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.API("🔻 Encerrando API...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
