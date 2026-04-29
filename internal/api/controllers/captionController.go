package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type CaptionController struct {
	container *container.AppContainer
}

func NewCaptionController(container *container.AppContainer) *CaptionController {
	return &CaptionController{
		container: container,
	}
}

func (c *CaptionController) UpdateDefaultCaptionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var captionData types.CaptionDefaultUpdateRequest
	if err := ctx.ShouldBindJSON(&captionData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateDefaultCaptionService(ctx, channelId, captionData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *CaptionController) UpdateNewPackCaptionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var captionData types.CaptionDefaultUpdateRequest
	if err := ctx.ShouldBindJSON(&captionData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateNewPackCaptionService(ctx, channelId, captionData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *CaptionController) UpdateReactionsController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var reactionsData types.ReactionsUpdateRequest
	if err := ctx.ShouldBindJSON(&reactionsData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateReactionsService(ctx, channelId, reactionsData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *CaptionController) UpdateReactionPositionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var posData types.ReactionPositionUpdateRequest
	if err := ctx.ShouldBindJSON(&posData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateReactionPositionService(ctx, channelId, posData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
