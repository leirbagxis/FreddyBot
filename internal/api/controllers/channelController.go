package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type ChannelController struct {
	container *container.AppContainer
}

func NewChannelController(container *container.AppContainer) *ChannelController {
	return &ChannelController{
		container: container,
	}
}

func (c *ChannelController) DisconectChannel(ctx *gin.Context) {
	channelIDStr, exists := ctx.Get("channelID")
	channelID, ok := channelIDStr.(int64)
	if !exists || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ channelID invalido no contexto!",
		})
		return
	}

	channel, err := c.container.ChannelRepo.GetChannelByID(ctx, channelID)
	fmt.Println(channel, err)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ Erro ao encontrar canal!",
		})
		return
	}

	err = c.container.ChannelRepo.DeleteChannelWithRelations(ctx, channel.OwnerID, channelID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ Erro ao deletar canal!",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
