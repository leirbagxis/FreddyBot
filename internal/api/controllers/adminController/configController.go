package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
)

type ConfigController struct {
	container *container.AppContainer
}

func NewConfigController(container *container.AppContainer) *ConfigController {
	return &ConfigController{container: container}
}

func (ctrl *ConfigController) GetConfig(ctx *gin.Context) {
	config, err := ctrl.container.ServerService.GetConfig(ctx)
	if err != nil {
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao buscar configurações"))
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(config))
}

func (ctrl *ConfigController) UpdateConfig(ctx *gin.Context) {
	var body struct {
		Maintenance bool `json:"maintence"`
		ForceJoin   bool `json:"forceJoin"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(errors.BadRequest("Dados inválidos"))
		return
	}

	config, err := ctrl.container.ServerService.UpdateConfig(ctx, body.Maintenance, body.ForceJoin)
	if err != nil {
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao atualizar configurações"))
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(config, "Configurações atualizadas com sucesso"))
}
