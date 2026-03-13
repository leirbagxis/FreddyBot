package adminrepositories

import (
	"context"
	"fmt"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type AdminRepositories struct {
	db *gorm.DB
}

func NewAdminRepositories(db *gorm.DB) *AdminRepositories {
	return &AdminRepositories{
		db: db,
	}
}

func (r *AdminRepositories) GetAllUsersAdminRepository(ctx context.Context) ([]models.User, error) {
	var users []models.User

	result := r.db.WithContext(ctx).
		Preload("Channels").
		Order("updated_at DESC").
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", result.Error)
	}

	fmt.Println("✅ Usuários listados com sucesso")

	return users, nil
}
