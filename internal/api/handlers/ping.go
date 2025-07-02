package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func PingHandler(c *container.AppContainer) gin.HandlerFunc {
	return func(g *gin.Context) {
		res := map[string]any{
			"ping": "pong",
		}
		g.JSON(http.StatusOK, res)
	}
}
