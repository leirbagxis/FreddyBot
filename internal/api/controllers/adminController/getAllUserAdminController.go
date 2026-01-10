package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type UsersAdminController struct {
	container *container.AppContainer
}

func NewUsersAdminController(app *container.AppContainer) *UsersAdminController {
	return &UsersAdminController{
		container: app,
	}
}

func (c *UsersAdminController) GetAllUsersAdminController(ctx *gin.Context) {
	users, err := c.container.AdminService.GetAllUsersAdminRepository(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)

}
