package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type ButtonsController struct {
	container *container.AppContainer
}

func NewButtonsController(container *container.AppContainer) *ButtonsController {
	return &ButtonsController{
		container: container,
	}
}

func (c *ButtonsController) CreateDefaultButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var buttonData types.ButtonCreateRequest
	if err := ctx.ShouldBindJSON(&buttonData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.CreateButtonService(ctx, channelId, buttonData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !result.Success {
		ctx.JSON(http.StatusOK, result)
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (c *ButtonsController) DeleteDefaultButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	buttonID := ctx.Param("buttonId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	err = appService.DeleteDefaulfButtonService(ctx, channelId, buttonID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, "")

}

func (c *ButtonsController) UpdateDefaultButtonController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	buttonID := ctx.Param("buttonId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var buttonData types.ButtonCreateRequest
	if err := ctx.ShouldBindJSON(&buttonData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateButtonService(ctx, channelId, buttonID, buttonData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !result.Success {
		ctx.JSON(http.StatusOK, result)
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (c *ButtonsController) UpdateLayoutDefaultButtons(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")

	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	var body types.UpdateLayoutRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(c.container)
	result, err := appService.UpdateButtonsLayoutService(ctx, channelId, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !result.Success {
		ctx.JSON(http.StatusOK, result)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
