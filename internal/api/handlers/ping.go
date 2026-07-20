package handlers

import (
	"io"
	"log"
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

// ClientErrorHandler logs JavaScript errors sent from the frontend.
// This is a best-effort, fire-and-forget endpoint - no auth required
// since it's meant to capture errors even before login completes.
func ClientErrorHandler(c *container.AppContainer) gin.HandlerFunc {
	return func(g *gin.Context) {
		body, err := io.ReadAll(g.Request.Body)
		if err != nil {
			g.Status(http.StatusOK)
			return
		}

		log.Printf("[CLIENT-ERROR] %s", string(body))
		g.Status(http.StatusOK)
	}
}
