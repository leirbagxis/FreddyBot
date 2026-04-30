package repositories

import (
	"context"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type ServerConfig struct {
	db           *gorm.DB
	cacheService *cache.Service
}

func NewServerConfigRepository(db *gorm.DB, cacheService *cache.Service) *ServerConfig {
	return &ServerConfig{db: db, cacheService: cacheService}
}

func (r *ServerConfig) GetMaintence(ctx context.Context) (bool, error) {
	// Try cache first
	if r.cacheService != nil {
		var maintenance bool
		err := r.cacheService.Get(ctx, "server:maintenance", &maintenance)
		if err == nil {
			return maintenance, nil
		}
	}

	var config models.ServerConfig
	if err := r.db.WithContext(ctx).First(&config, 1).Error; err != nil {
		return false, err
	}

	// Set cache
	if r.cacheService != nil {
		_ = r.cacheService.Set(ctx, "server:maintenance", config.Maintence, 30*time.Minute)
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

	// Update cache
	if r.cacheService != nil {
		_ = r.cacheService.Set(ctx, "server:maintenance", config.Maintence, 30*time.Minute)
	}

	return config.Maintence, nil
}

func (r *ServerConfig) GetConfig(ctx context.Context) (*models.ServerConfig, error) {
	var config models.ServerConfig
	if err := r.db.WithContext(ctx).First(&config, 1).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ServerConfig) UpdateConfig(ctx context.Context, maintenance, forceJoin bool) (*models.ServerConfig, error) {
	var config models.ServerConfig
	if err := r.db.WithContext(ctx).First(&config, 1).Error; err != nil {
		return nil, err
	}

	config.Maintence = maintenance
	config.ForceJoin = forceJoin

	if err := r.db.WithContext(ctx).Save(&config).Error; err != nil {
		return nil, err
	}

	// Update cache
	if r.cacheService != nil {
		_ = r.cacheService.Set(ctx, "server:maintenance", config.Maintence, 30*time.Minute)
	}

	return &config, nil
}
