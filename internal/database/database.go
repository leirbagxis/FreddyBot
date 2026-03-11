package database

import (
	//"github.com/glebarez/sqlite"
	"context"
	"log"

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
		&models.ServerConfig{},
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

	if err := initServerConfig(db); err != nil {
		panic(err)
	}

	return db
}

func initServerConfig(db *gorm.DB) error {
	config := models.ServerConfig{
		ID:        1,
		Maintence: false,
	}

	if err := db.WithContext(context.Background()).FirstOrCreate(&config, models.ServerConfig{ID: 1}).Error; err != nil {
		return err
	}

	log.Println("✔️ ServerConfig iniciado criadas com sucesso.")
	return nil
}
