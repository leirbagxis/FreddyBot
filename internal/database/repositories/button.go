package repositories

import (
	"context"
	"fmt"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
	"gorm.io/gorm"
)

type ButtonRepository struct {
	db *gorm.DB
}

func NewButtonRepository(db *gorm.DB) *ButtonRepository {
	return &ButtonRepository{
		db: db,
	}
}

func (r *ButtonRepository) GetUserChannelsAsButtons(ctx context.Context, userID int64) ([][]parser.Button, error) {
	var channels []models.Channel
	err := r.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&channels).Error
	if err != nil || len(channels) == 0 {
		return nil, err
	}

	var buttons [][]parser.Button
	for _, channel := range channels {
		row := []parser.Button{
			{
				Text:         channel.Title,
				CallbackData: fmt.Sprintf("config:%d", channel.ID),
			},
		}
		buttons = append(buttons, row)
	}

	return buttons, nil

}
