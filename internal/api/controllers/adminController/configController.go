package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
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
		logger.Error("API", "Erro ao buscar configurações no ServerService: %v", err)
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao buscar configurações"))
		return
	}

	logger.Bot("Configuração recuperada com sucesso: %+v", config)
	ctx.JSON(http.StatusOK, types.NewSuccessResponse(config))
}

func (ctrl *ConfigController) UpdateConfig(ctx *gin.Context) {
	var body struct {
		Maintenance          bool   `json:"maintence"`
		ForceJoin            bool   `json:"forceJoin"`
		GlobalDefaultCaption string `json:"globalDefaultCaption"`
		GlobalNewPackCaption string `json:"globalNewPackCaption"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("API", "Erro ao fazer bind do JSON no UpdateConfig: %v", err)
		ctx.Error(errors.BadRequest("Dados inválidos"))
		return
	}

	logger.Bot("Recebendo atualização de configuração: %+v", body)

	config, err := ctrl.container.ServerService.UpdateConfig(ctx, body.Maintenance, body.ForceJoin, body.GlobalDefaultCaption, body.GlobalNewPackCaption)
	if err != nil {
		logger.Error("API", "Erro ao atualizar configuração no ServerService: %v", err)
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao atualizar configurações"))
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(config, "Configurações atualizadas com sucesso"))
}
