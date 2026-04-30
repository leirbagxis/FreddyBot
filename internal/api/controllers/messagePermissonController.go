package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type PermissionController struct {
	container *container.AppContainer
}

func NewPermissionController(container *container.AppContainer) *PermissionController {
	return &PermissionController{
		container: container,
	}
}

func (ctrl *PermissionController) UpdateMessagePermissionController(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	bodyRaw, _ := io.ReadAll(c.Request.Body)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	var body types.UpdateMessagePermissionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.UpdateMessagePermissionService(c, channelID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, result)

}

func (ctrl *PermissionController) UpdateButtonsPermissionController(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	bodyRaw, _ := io.ReadAll(c.Request.Body)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	var body types.UpdateButtonsPermissionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.UpdateButtonsPermissionService(c, channelID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, result)

}

func (ctrl *PermissionController) UpdateReactionsActiveController(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
		return
	}

	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	appService := (*service.AppContainerLocal)(ctrl.container)
	result, err := appService.UpdateReactionsActiveService(c, channelID, body.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

