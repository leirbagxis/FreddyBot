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

func (app *AppContainerLocal) CreateCustomCaptionService(ctx context.Context, channelID int64, body types.CreateCustomCaptionRequest) (*types.CreateCustomCaptionResponse, error) {
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
	newCaption := models.CustomCaption{
		CaptionID:      uuid.NewString(),
		Code:           body.Code,
		Caption:        body.Caption,
		LinkPreview:    body.LinkPreview,
		OwnerChannelID: channelID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := app.DB.WithContext(ctx).Create(&newCaption).Error; err != nil {
		return nil, fmt.Errorf("erro ao criar custom caption: %w", err)
	}

	return &types.CreateCustomCaptionResponse{
		Success: true,
		Message: "Custom caption criada com sucesso",
		Data: map[string]interface{}{
			"captionId":   newCaption.CaptionID,
			"code":        newCaption.Code,
			"caption":     newCaption.Caption,
			"linkPreview": newCaption.LinkPreview,
			"createdAt":   newCaption.CreatedAt,
		},
	}, nil
}

func (app *AppContainerLocal) DeleteCustomCaptionService(ctx context.Context, channelID int64, captionID string) error {
	if captionID == "" {
		return fmt.Errorf("ID do caption é obrigatório")
	}

	if channelID == 0 {
		return fmt.Errorf("ID do canal é obrigatório")
	}

	result := app.DB.WithContext(ctx).
		Where("caption_id = ? AND owner_channel_id = ?", captionID, channelID).
		Delete(&models.CustomCaption{})

	if result.Error != nil {
		return fmt.Errorf("erro ao deletar customCaption: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("customCaption não encontrado ou não pertence ao canal")
	}

	return nil
}
