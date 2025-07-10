package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type CustomCaptionController struct {
	container *container.AppContainer
}

func NewCustomCaptionController(container *container.AppContainer) *CustomCaptionController {
	return &CustomCaptionController{
		container: container,
	}
}

func (ctrl *CustomCaptionController) Create(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelID, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	var body types.CreateCustomCaptionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido", "details": err.Error()})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.CreateCustomCaptionService(ctx, channelID, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)

}

func (c *CustomCaptionController) DeleteCustomCaptionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	captionID := ctx.Param("captionId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	err = appService.DeleteCustomCaptionService(ctx, channelId, captionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, "")

}
