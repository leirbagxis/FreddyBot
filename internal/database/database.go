package database

import (
	//"github.com/glebarez/sqlite"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

		&models.SubscriptionPlan{},
		&models.Subscription{},
		&models.Coupon{},
		&models.GiftCard{},
	)
	if err != nil {
		panic(err)
	}

	SeedPlans(db)

	return db
}

func SeedPlans(db *gorm.DB) {
	plans := []models.SubscriptionPlan{
		{
			ID:         0001,
			Name:       "Gratuito",
			PricePix:   0,
			PriceStars: 0,
			Duration:   0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         0002,
			Name:       "Gold",
			PricePix:   10.99, // R$15,00
			PriceStars: 100,   // 150 Stars
			Duration:   30,    // 30 dias
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, plan := range plans {
		db.Clauses(
			clause.OnConflict{UpdateAll: true},
		).Create(&plan)
	}
}
