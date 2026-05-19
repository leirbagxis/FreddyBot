package database

import (
	"context"
	"time"

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

	// Habilitar Foreign Keys no SQLite
	if config.AppEnv == "dev" {
		db.Exec("PRAGMA foreign_keys = ON;")
	}

	// Configurar Pool de Conexões (Crucial para produção)
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		customLogger.DB("⚙️ Pool de conexões configurado (Idle: 10, Open: 100)")
	}

	// Forçar recriação de índices que mudaram de estrutura
	db.Exec("DROP INDEX IF EXISTS idx_vote_user")

	err = db.AutoMigrate(
		&models.User{},
		&models.ServerConfig{},
		&models.Channel{},
		&models.DefaultCaption{},
		&models.MessagePermission{},
		&models.ButtonsPermission{},
		&models.Button{},
		&models.Separator{},
		&models.CustomCaption{},
		&models.CustomCaptionButton{},
		&models.Vote{},
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
		ID:                   1,
		Maintence:            false,
		ForceJoin:            false,
		GlobalDefaultCaption: "🐈‍⠀៹ [t.me/legendasbot](https://t.me/{usernameBot})  ‹",
		GlobalNewPackCaption: `╔═━──━═༻✧༺═━──━═╗

        𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
        (｡•́︿•̀｡)っ✧.*ೃ༄
        ˗ˏˋ [$name]($link) ⋆｡˚ ☁︎
            彡♡ ₊˚

⋆｡˚ ❀ @LegendasBrBot ☽⁺₊

╚═━──━═༻✧༺═━──━═╝`,
	}

	if err := db.WithContext(context.Background()).FirstOrCreate(&config, models.ServerConfig{ID: 1}).Error; err != nil {
		return err
	}

	customLogger.DB("✔️ ServerConfig iniciado criadas com sucesso.")
	return nil
}
