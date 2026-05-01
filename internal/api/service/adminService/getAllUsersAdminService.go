package adminservice

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

type AdminService struct {
	*service.AppContainerLocal
}

func (app *AdminService) GetAllUsersAdminService(ctx context.Context) ([]models.User, error) {
	var users []models.User

	result := app.DB.WithContext(ctx).Preload("Channels").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	logger.Bot("✅ Usuários listados com sucesso no Admin")

	return users, nil
}
