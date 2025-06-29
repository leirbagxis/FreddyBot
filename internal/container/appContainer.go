package container

import (
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	UserRepo    *repositories.UserRepository
	ChannelRepo *repositories.ChannelRepository
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	return &AppContainer{
		DB:          db,
		UserRepo:    repositories.NewUserRepository(db),
		ChannelRepo: repositories.NewChannelRepository(db),
	}
}
