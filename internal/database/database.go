package database

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	customLogger "github.com/leirbagxis/FreddyBot/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultGlobalDefaultCaption = "🐈‍⠀៹ [t.me/legendasbot](https://t.me/{usernameBot})  ‹"

const defaultGlobalNewPackCaption = `╔═━──━═༻✧༺═━──━═╗

        𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
        (｡•́︿•̀｡)っ✧.*ೃ༄
        ˗ˏˋ [$name]($link) ⋆｡˚ ☁︎
            彡♡ ₊˚

⋆｡˚ ❀ @{usernameBot} ☽⁺₊

╚═━──━═༻✧༺═━──━═╝`

const defaultFixedPostBuilderPayload = `{"media_type":"photo","media_file_id":"AgACAgEAAxkBAAIN1GoO7mINPlBGs_ydPnmkDPdxeQ8eAAKoC2sbf_d4RIZ9nu_0BSIiAQADAgADeAADOwQ","menu_message_id":0,"prompt_message_id":0,"title":"","body":"<tg-emoji emoji-id=\"5373026167722876724\">🤩</tg-emoji> Cansado de perder tempo editando postagens?\nO <a href=\"http://t.me/LegendasBrBot?start=start\">LegendasBOT</a> resolve isso pra você de forma simples e eficiente <tg-emoji emoji-id=\"5445284980978621387\">🚀</tg-emoji>","footer":"","reactions":"","buttons":[{"text":"🤖 Legendas BOT","url":"http://t.me/LegendasBrBot?start=start","custom_emoji_id":"5296447931627352804"},{"text":"📺 Central de Novidades","url":"https://t.me/LegendasBOTTopic","custom_emoji_id":"5373330964372004748"}],"step":""}`

func DefaultFixedPostBuilderPayload() string {
	return defaultFixedPostBuilderPayload
}

func validFixedPostBuilderPayload(payload string) bool {
	if strings.TrimSpace(payload) == "" {
		return false
	}
	var raw map[string]any
	return json.Unmarshal([]byte(payload), &raw) == nil
}

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
		&models.ChannelEvent{},
		&models.DefaultCaption{},
		&models.MessagePermission{},
		&models.ButtonsPermission{},
		&models.Button{},
		&models.Separator{},
		&models.CustomCaption{},
		&models.CustomCaptionButton{},
		&models.Vote{},
		&models.UserTelegramSession{},
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
		ID:                      1,
		Maintence:               false,
		ForceJoin:               false,
		GlobalDefaultCaption:    defaultGlobalDefaultCaption,
		FixedPostBuilderEnabled: true,
		FixedPostBuilderKey:     "legendasbot",
		FixedPostBuilderPayload: defaultFixedPostBuilderPayload,
		GlobalNewPackCaption:    defaultGlobalNewPackCaption,
	}

	if err := db.WithContext(context.Background()).FirstOrCreate(&config, models.ServerConfig{ID: 1}).Error; err != nil {
		return err
	}

	changed := false
	if strings.TrimSpace(config.GlobalDefaultCaption) == "" {
		config.GlobalDefaultCaption = defaultGlobalDefaultCaption
		changed = true
	}
	if strings.TrimSpace(config.GlobalNewPackCaption) == "" {
		config.GlobalNewPackCaption = defaultGlobalNewPackCaption
		changed = true
	}
	if config.FixedPostBuilderKey == "" {
		config.FixedPostBuilderKey = "legendasbot"
		changed = true
	}
	if !validFixedPostBuilderPayload(config.FixedPostBuilderPayload) {
		config.FixedPostBuilderPayload = defaultFixedPostBuilderPayload
		config.FixedPostBuilderEnabled = true
		changed = true
	}
	if changed {
		if err := db.WithContext(context.Background()).Save(&config).Error; err != nil {
			return err
		}
	}

	customLogger.DB("✔️ ServerConfig iniciado criadas com sucesso.")
	return nil
}
