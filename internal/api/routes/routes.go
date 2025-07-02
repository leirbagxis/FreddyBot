package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/handlers"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func RegisterRoutes(r *gin.Engine, c *container.AppContainer) {
	api := r.Group("/api")
	{
		api.GET("/ping", handlers.PingHandler(c))
	}
}
