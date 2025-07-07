package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/handlers"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func RegisterRoutes(r *gin.Engine, c *container.AppContainer) {
	loginRouter := r.Group("/auth")
	loginRouter.GET("/verify-token", handlers.VerifyJWTHandler())
	loginRouter.POST("/generate-token", handlers.GenerateJWTHandler(c))

	api := r.Group("/api")
	api.Use(auth.AuthMiddlewareJWT())
	{
		api.GET("/ping", handlers.PingHandler(c))
		api.GET("/channel/:channelId", handlers.GetChannelHandler(c))

	}
}
