package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

type ChannelController struct {
	container *container.AppContainer
}

func NewChannelController(container *container.AppContainer) *ChannelController {
	return &ChannelController{
		container: container,
	}
}

func (c *ChannelController) GetAllChannelsController(ctx *gin.Context) {
	channels, err := c.container.ChannelRepo.GetAllChannels(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao buscar todos os canais"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"channels": channels,
	})
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
	if err != nil {
		logger.Error("API", "Erro ao encontrar canal %d: %v", channelID, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ Erro ao encontrar canal!",
		})
		return
	}

	err = c.container.DisconnectChannel(ctx, channel.OwnerID, channelID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ Erro ao deletar canal!",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
