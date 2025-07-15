package repositories

import (
	"context"
	"errors"
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

func (r *ButtonRepository) CreateButton(ctx context.Context, button *models.Button) error {
	if button == nil {
		return fmt.Errorf("botão não pode ser nil")
	}

	err := r.db.WithContext(ctx).Create(button).Error
	if err != nil {
		return fmt.Errorf("erro ao criar botão: %w", err)
	}

	return nil
}

// func (r *ButtonRepository) GetFirstButton(ctx context.Context, channelID int64) (*models.Button, error) {
// 	var button models.Button
// 	err := r.db.WithContext(ctx).Where("owner_channel_id = ?", channelID).FirstOrCreate(&button).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &button, nil
// }

// Função para buscar o primeiro botão (mais antigo) de um canal
func (r *ButtonRepository) GetFirstButton(ctx context.Context, channelID int64) (*models.Button, error) {
	var button models.Button

	err := r.db.WithContext(ctx).
		Where("owner_channel_id = ?", channelID).
		Order("created_at ASC"). // Ordenar por data de criação (mais antigo primeiro)
		First(&button).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("nenhum botão encontrado para o canal %d", channelID)
		}
		return nil, fmt.Errorf("erro ao buscar primeiro botão: %w", err)
	}

	return &button, nil
}
