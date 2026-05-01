package adminrepositories

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
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
		return nil, result.Error
	}

	logger.Bot("✅ Usuários listados com sucesso no repositório Admin")

	return users, nil
}
