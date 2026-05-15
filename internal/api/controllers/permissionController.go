package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
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
		c.Error(errors.BadRequest("channelId inválido"))
		return
	}

	bodyRaw, _ := io.ReadAll(c.Request.Body)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	var body types.UpdateMessagePermissionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	result, err := ctrl.container.PermissionsService.UpdateMessagePermission(c, channelID, body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, types.NewSuccessResponse(result))

}

func (ctrl *PermissionController) UpdateButtonsPermissionController(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("channelId inválido"))
		return
	}

	bodyRaw, _ := io.ReadAll(c.Request.Body)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	var body types.UpdateButtonsPermissionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	result, err := ctrl.container.PermissionsService.UpdateButtonsPermission(c, channelID, body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, types.NewSuccessResponse(result))

}

func (ctrl *PermissionController) UpdateReactionsActiveController(c *gin.Context) {
	channelIDStr := c.Param("channelId")
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("channelId inválido"))
		return
	}

	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	result, err := ctrl.container.PermissionsService.UpdateReactionsActive(c, channelID, body.Active)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, types.NewSuccessResponse(result))
}
