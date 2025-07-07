package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func GetChannelHandler(app *container.AppContainer) gin.HandlerFunc {
	return func(c *gin.Context) {

		channelIdStr := c.Param("channelId")
		channelId, err := strconv.Atoi(channelIdStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
			return
		}

		channel, err := app.ChannelRepo.GetChannelByID(&gin.Context{}, int64(channelId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canal não encontrado"})
			return
		}

		user, err := app.UserRepo.GetUserById(&gin.Context{}, channel.OwnerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}

		data := map[string]any{
			"user":    user,
			"channel": channel,
		}

		c.JSON(http.StatusOK, data)

	}
}
