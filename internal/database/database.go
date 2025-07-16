package database

import (
	//"github.com/glebarez/sqlite"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DatabaseFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Config.Logger = logger.Default.LogMode(logger.Silent)

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
