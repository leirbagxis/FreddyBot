package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func GetChannelByTwoID(app *container.AppContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		channelIdStr := c.Param("channelId")
		channelId, err := strconv.Atoi(channelIdStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
			return
		}

		user, err := app.UserRepo.GetUserById(&gin.Context{}, int64(userId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}

		channel, err := app.ChannelRepo.GetChannelByTwoID(&gin.Context{}, int64(userId), int64(channelId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canal não encontrado"})
			return
		}

		data := map[string]any{
			"user":    user,
			"channel": channel,
		}

		c.JSON(http.StatusOK, data)

	}
}
