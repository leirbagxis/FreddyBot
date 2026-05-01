package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

type AppContainerLocal container.AppContainer

func (app *AppContainerLocal) isEmoji(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

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
		return nil, result.Error
	}

	logger.Bot("✅ Legenda padrão atualizada com sucesso (Canal: %d)", channelID)

	return &types.CaptionUpdateResponse{
		Success: true,
		Message: "Legenda atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

func (app *AppContainerLocal) UpdateNewPackCaptionService(ctx context.Context, channelID int64, captionData types.CaptionDefaultUpdateRequest) (*types.CaptionUpdateResponse, error) {
	if len(captionData.Caption) > 4096 {
		return nil, errors.New("Caption muito longa (máximo 4096 caracteres)")
	}

	now := time.Now()
	result := app.DB.Model(&models.Channel{}).
		Where("id = ?", channelID).
		Updates(map[string]interface{}{
			"NewPackCaption": captionData.Caption,
			"updated_at":     now,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	logger.Bot("✅ NewPackCaption atualizada com sucesso (Canal: %d)", channelID)

	return &types.CaptionUpdateResponse{
		Success: true,
		Message: "NewPackCaption atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

func (app *AppContainerLocal) UpdateReactionsService(ctx context.Context, channelID int64, reactionsData types.ReactionsUpdateRequest) (*types.CaptionUpdateResponse, error) {
	// Validação de emojis
	if reactionsData.Reactions != "" {
		parts := strings.Split(reactionsData.Reactions, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if !app.isEmoji(p) {
				return nil, errors.New("apenas emojis são permitidos como reações")
			}
		}
	}

	now := time.Now()
	result := app.DB.Model(&models.Channel{}).
		Where("id = ?", channelID).
		Updates(map[string]interface{}{
			"reactions":  reactionsData.Reactions,
			"updated_at": now,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar reações: %w", result.Error)
	}

	return &types.CaptionUpdateResponse{
		Success: true,
		Message: "Reações atualizadas com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

func (app *AppContainerLocal) UpdateReactionPositionService(ctx context.Context, channelID int64, posData types.ReactionPositionUpdateRequest) (*types.CaptionUpdateResponse, error) {
	// Verify if any button is already in this row
	var count int64
	app.DB.Model(&models.Button{}).Where("owner_channel_id = ? AND position_y = ?", channelID, posData.ReactionPosition).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("esta linha já possui botões e não pode ser usada para reações")
	}

	now := time.Now()
	result := app.DB.Model(&models.Channel{}).
		Where("id = ?", channelID).
		Updates(map[string]interface{}{
			"reaction_position": posData.ReactionPosition,
			"updated_at":        now,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar posição das reações: %w", result.Error)
	}

	return &types.CaptionUpdateResponse{
		Success: true,
		Message: "Posição das reações atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

