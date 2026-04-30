package service

import (
	"context"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
)

func (app *AppContainerLocal) UpdateMessagePermissionService(ctx context.Context, channelID int64, messagePermissionsData types.UpdateMessagePermissionRequest) (*types.UpdatePermissionsResponse, error) {
	// Subquery rápida para encontrar o ID da permissão vinculada ao canal sem carregar tudo
	result := app.DB.Model(&models.MessagePermission{}).
		Where("owner_caption_id = (SELECT caption_id FROM default_captions WHERE owner_channel_id = ?)", channelID).
		Updates(messagePermissionsData)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar permissao message: %w", result.Error)
	}

	return &types.UpdatePermissionsResponse{
		Success: true,
		Message: "MessagePermission atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    time.Now(),
		},
	}, nil
}

func (app *AppContainerLocal) UpdateButtonsPermissionService(ctx context.Context, channelID int64, buttonsPermissionsData types.UpdateButtonsPermissionRequest) (*types.UpdatePermissionsResponse, error) {
	// Atualização direta via subquery para performance máxima
	result := app.DB.Model(&models.ButtonsPermission{}).
		Where("owner_caption_id = (SELECT caption_id FROM default_captions WHERE owner_channel_id = ?)", channelID).
		Updates(buttonsPermissionsData)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar ButtonsPermission: %w", result.Error)
	}

	return &types.UpdatePermissionsResponse{
		Success: true,
		Message: "ButtonsPermission atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    time.Now(),
		},
	}, nil
}

func (app *AppContainerLocal) UpdateReactionsActiveService(ctx context.Context, channelID int64, active bool) (*types.UpdatePermissionsResponse, error) {
	result := app.DB.Model(&models.MessagePermission{}).
		Where("owner_caption_id = (SELECT caption_id FROM default_captions WHERE owner_channel_id = ?)", channelID).
		Update("reactions", active)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar reacoes: %w", result.Error)
	}

	return &types.UpdatePermissionsResponse{
		Success: true,
		Message: "Status das reações atualizado com sucesso",
		Data: map[string]interface{}{
			"active":     active,
			"updated_at": time.Now(),
		},
	}, nil
}
