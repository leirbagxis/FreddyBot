package database

import (
	"context"

	"github.com/glebarez/sqlite"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	customLogger "github.com/leirbagxis/FreddyBot/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	var dialector gorm.Dialector

	if config.AppEnv == "dev" {
		customLogger.DB("📦 Usando banco de dados SQLite (modo dev)")
		dialector = sqlite.Open(config.DatabaseFile)
	} else {
		customLogger.DB("🐘 Usando banco de dados PostgreSQL (modo prod)")
		dialector = postgres.Open(config.DatabaseFile)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
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

	customLogger.DB("✔️ ServerConfig iniciado criadas com sucesso.")
	return nil
}
