package adminservice

import (
	"context"
	"fmt"

	"github.com/leirbagxis/FreddyBot/internal/api/service"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
)

type AdminService struct {
	*service.AppContainerLocal
}

func (app *AdminService) GetAllUsersAdminService(ctx context.Context) ([]models.User, error) {
	var users []models.User

	result := app.DB.WithContext(ctx).Preload("Channels").
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", result.Error)
	}

	fmt.Println("✅ Usuários listados com sucesso")

	return users, nil
}
