package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
)

type AppContainerLocal container.AppContainer

func (app *AppContainerLocal) UpdateDefaultCaptionService(ctx context.Context, channelID int64, captionData types.CaptionDefaultUpdateRequest) (*types.CaptionUpdateResponse, error) {
	if len(captionData.Caption) > 4096 {
		return nil, errors.New("Caption muito longa (máximo 4096 caracteres)")
	}

	now := time.Now()
	result := app.DB.Model(&models.DefaultCaption{}).
		Where("owner_channel_id = ?", channelID).
		Updates(map[string]interface{}{
			"caption":    captionData.Caption,
			"updated_at": now,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar legenda padrão: %w", result.Error)
	}

	fmt.Println("✅ Legenda padrão atualizada com sucesso")

	return &types.CaptionUpdateResponse{
		Success: true,
		Message: "Legenda atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}
