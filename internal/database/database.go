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

	db.AutoMigrate(
		&models.User{},
		&models.Channel{},
		&models.DefaultCaption{},
		&models.MessagePermission{},
		&models.ButtonPermission{},
		&models.Button{},
	)
	return db
}
