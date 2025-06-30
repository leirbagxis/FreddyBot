package database

import (
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DatabaseFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Channel{},
		&models.DefaultCaption{},
		&models.MessagePermission{},
		&models.ButtonsPermission{},
		&models.Button{},
		&models.Separator{},
		&models.CustomCaption{},
		&models.CustomCaptionButton{},
	)
	if err != nil {
		panic(err)
	}

	return db
}
