package container

import (
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	UserRepo    *repositories.UserRepository
	ChannelRepo *repositories.ChannelRepository
	ButtonRepo  *repositories.ButtonRepository

	// ## CACHE ## \\
	CacheService   *cache.Service
	SessionManager *cache.SessionManager
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	cacheService := cache.NewService()
	return &AppContainer{
		DB:          db,
		UserRepo:    repositories.NewUserRepository(db),
		ChannelRepo: repositories.NewChannelRepository(db),
		ButtonRepo:  repositories.NewButtonRepository(db),

		CacheService:   cacheService,
		SessionManager: cache.NewSessionManager(cacheService),
	}
}
