package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
)

func (ctrl *CustomCaptionController) CreateCustomCaptionButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	captionID := ctx.Param("captionId")

	channelID, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	var body types.ButtonCreateRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido", "details": err.Error()})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.CreateCustomCaptionButtonService(ctx, channelID, captionID, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)

}

func (ctrl *CustomCaptionController) UpdateCustomCaptionButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	captionID := ctx.Param("captionId")
	buttonID := ctx.Param("buttonId")

	channelID, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	var body types.ButtonCreateRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido", "details": err.Error()})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.UpdateCustomCaptionButtonService(ctx, channelID, captionID, buttonID, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)

}

func (ctrl *CustomCaptionController) UpdateCustomCaptionLayoutController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	captionID := ctx.Param("captionId")

	channelID, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	var body types.UpdateCustomCaptionLayoutRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido", "details": err.Error()})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.UpdateCustomCaptionLayoutService(ctx, channelID, captionID, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)

}

func (ctrl *CustomCaptionController) DeleteCustomCaptionButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	captionID := ctx.Param("captionId")
	buttonID := ctx.Param("buttonId")

	channelID, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	err = appService.DeleteCustomCaptionButtonService(ctx, channelID, captionID, buttonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "")

}
