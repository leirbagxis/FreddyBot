package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

func (app *AppContainerLocal) GetCustomCaptionButtons(ctx context.Context, captionID string) ([]models.CustomCaptionButton, error) {
	var buttons []models.CustomCaptionButton

	err := app.DB.WithContext(ctx).Where("owner_caption_id = ?", captionID).Order("position_y ASC").Order("position_x ASC").Find(&buttons).Error
	if err != nil {
		return nil, err
	}

	return buttons, err
}

func (app *AppContainerLocal) CreateCustomCaptionButtonService(ctx context.Context, channelID int64, captionID string, body types.ButtonCreateRequest) (*types.CreateCustomCaptionResponse, error) {
	if err := validateButtonData(body); err != nil {
		return nil, err
	}

	position, err := app.CalculateNextCustomButtonPosition(ctx, captionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular posição: %w", err)
	}

	var channel models.Channel
	if err := app.DB.WithContext(ctx).
		Where("id = ? ", channelID).
		First(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("canal não encontrado ou não pertence ao usuário")
		}
		return nil, fmt.Errorf("erro ao buscar canal: %w", err)
	}

	now := time.Now()
	newCaption := models.CustomCaptionButton{
		ButtonID:       uuid.NewString(),
		NameButton:     body.NameButton,
		ButtonURL:      body.ButtonURL,
		PositionX:      position.X,
		PositionY:      position.Y,
		OwnerCaptionID: captionID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := app.DB.WithContext(ctx).Create(&newCaption).Error; err != nil {
		return nil, fmt.Errorf("erro ao criar custom caption: %w", err)
	}

	fmt.Printf("✅ Botão Customizado criado com sucesso: %s\n", newCaption.ButtonID)
	return &types.CreateCustomCaptionResponse{
		Success: true,
		Message: "Button Custom Caption criada com sucesso",
		Data: map[string]interface{}{
			"buttonId":   newCaption.ButtonID,
			"nameButton": newCaption.NameButton,
			"buttonUrl":  newCaption.ButtonURL,
			"createdAt":  newCaption.CreatedAt,
		},
	}, nil
}

func (app *AppContainerLocal) UpdateCustomCaptionButtonService(ctx context.Context, channelID int64, captionID, buttonID string, body types.ButtonCreateRequest) (*types.CreateCustomCaptionResponse, error) {
	if err := validateButtonData(body); err != nil {
		return nil, err
	}

	now := time.Now()
	result := app.DB.Model(&models.CustomCaptionButton{}).
		Where("button_id = ? AND owner_caption_id = ?", buttonID, captionID).
		Updates(map[string]interface{}{
			"name_button": body.NameButton,
			"button_url":  body.ButtonURL,
			"updated_at":  now,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar um botao customizado: %w", result.Error)
	}
	fmt.Printf("✅ Botão customizado atualizado com sucesso")

	return &types.CreateCustomCaptionResponse{
		Success: true,
		Message: "Botão customizado atualizado com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

func (app *AppContainerLocal) DeleteCustomCaptionButtonService(ctx context.Context, channelID int64, captionID, buttonID string) error {
	if buttonID == "" {
		return fmt.Errorf("ID do botão é obrigatório")
	}

	if captionID == "" {
		return fmt.Errorf("ID do caption é obrigatório")
	}

	if channelID == 0 {
		return fmt.Errorf("ID do canal é obrigatório")
	}

	result := app.DB.WithContext(ctx).
		Debug().
		Where("button_id = ? AND owner_caption_id = ?", buttonID, captionID).
		Delete(&models.CustomCaptionButton{})

	if result.Error != nil {
		return fmt.Errorf("erro ao deletar botão: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("botão não encontrado ou não pertence ao canal")
	}

	return nil
}

func (app *AppContainerLocal) UpdateCustomCaptionLayoutService(
	ctx context.Context,
	channelID int64,
	captionID string,
	layoutData types.UpdateCustomCaptionLayoutRequest,
) (*types.CreateCustomCaptionResponse, error) {
	if len(layoutData.Layout) == 0 {
		return nil, errors.New("layout não pode ser vazio")
	}

	var channel models.Channel
	if err := app.DB.WithContext(ctx).
		Where("id = ? ", channelID).
		First(&channel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("canal não encontrado ou sem permissão")
		}
		return nil, fmt.Errorf("erro ao buscar canal: %w", err)
	}

	var caption models.CustomCaption
	if err := app.DB.WithContext(ctx).
		Where("owner_channel_id = ? AND caption_id = ?", channelID, captionID).
		First(&caption).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("custom caption não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar custom caption: %w", err)
	}

	type ButtonPosition struct {
		ID        string
		PositionX int
		PositionY int
	}

	var updates []ButtonPosition
	for y, row := range layoutData.Layout {
		for x, btn := range row {
			updates = append(updates, ButtonPosition{
				ID:        btn.ID,
				PositionX: x,
				PositionY: y,
			})
		}
	}

	now := time.Now()
	err := app.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, btn := range updates {
			result := tx.Model(&models.CustomCaptionButton{}).
				Where("button_id = ? and owner_caption_id = ?", btn.ID, caption.CaptionID).
				Updates(map[string]interface{}{
					"position_x": btn.PositionX,
					"position_y": btn.PositionY,
					"updated_at": now,
				})
			if result.Error != nil {
				return result.Error
			}

		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar layout dos botões: %w", err)
	}

	return &types.CreateCustomCaptionResponse{
		Success: true,
		Message: "Layout dos botões da custom caption atualizado com sucesso",
		Data: map[string]interface{}{
			"updated_at": now,
			"total":      len(updates),
		},
	}, nil

}

func (app *AppContainerLocal) CalculateNextCustomButtonPosition(ctx context.Context, captionID string) (*types.ButtonPosition, error) {
	buttons, err := app.GetCustomCaptionButtons(ctx, captionID)
	if err != nil {
		return nil, err
	}

	if len(buttons) == 0 {
		return &types.ButtonPosition{X: 0, Y: 0}, nil
	}

	// Encontrar última linha (maxY)
	maxY := buttons[0].PositionY
	for _, button := range buttons {
		if button.PositionY > maxY {
			maxY = button.PositionY
		}
	}

	// Contar botões na última linha
	buttonsInLastRow := 0
	for _, button := range buttons {
		if button.PositionY == maxY {
			buttonsInLastRow++
		}
	}

	if buttonsInLastRow < 3 {
		return &types.ButtonPosition{X: buttonsInLastRow, Y: maxY}, nil
	}

	return &types.ButtonPosition{X: 0, Y: maxY + 1}, nil
}
