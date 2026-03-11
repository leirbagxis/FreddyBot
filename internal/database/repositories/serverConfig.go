package repositories

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type ServerConfig struct {
	db *gorm.DB
}

func NewServerConfigRepository(db *gorm.DB) *ServerConfig {
	return &ServerConfig{db: db}
}

func (r *ServerConfig) GetMaintence(ctx context.Context) (bool, error) {
	var config models.ServerConfig

	if err := r.db.WithContext(ctx).First(&config, 1).Error; err != nil {
		return false, err
	}

	return config.Maintence, nil
}

func (r *ServerConfig) ToggleMaintence(ctx context.Context) (bool, error) {
	var config models.ServerConfig

	if err := r.db.WithContext(ctx).First(&config, 1).Error; err != nil {
		return false, err
	}

	config.Maintence = !config.Maintence

	if err := r.db.WithContext(ctx).Save(&config).Error; err != nil {
		return false, err
	}

	return config.Maintence, nil
}
