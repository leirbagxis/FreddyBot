package service

import (
	"context"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
)

func (app *AppContainerLocal) UpdateMessagePermissionService(ctx context.Context, channelID int64, messagePermissionsData types.UpdateMessagePermissionRequest) (*types.UpdatePermissionsResponse, error) {
	var caption models.DefaultCaption
	err := app.DB.Where("owner_channel_id = ?", channelID).
		Preload("MessagePermission").
		First(&caption).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao criar botão: %w", err)
	}

	now := time.Now()
	result := app.DB.Model(&caption.MessagePermission).
		Updates(messagePermissionsData)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar permissao message: %w", result.Error)
	}

	fmt.Println("✅ MessagePermission atualizada com sucesso")

	return &types.UpdatePermissionsResponse{
		Success: true,
		Message: "MessagePermission atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil

}

func (app *AppContainerLocal) UpdateButtonsPermissionService(ctx context.Context, channelID int64, buttonsPermissionsData types.UpdateButtonsPermissionRequest) (*types.UpdatePermissionsResponse, error) {
	var caption models.DefaultCaption
	err := app.DB.Where("owner_channel_id = ?", channelID).
		Preload("ButtonsPermission").
		First(&caption).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao criar botão: %w", err)
	}

	now := time.Now()
	result := app.DB.Model(&caption.ButtonsPermission).
		Updates(buttonsPermissionsData)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar ButtonsPermission: %w", result.Error)
	}

	fmt.Println("✅ ButtonsPermission atualizada com sucesso")

	return &types.UpdatePermissionsResponse{
		Success: true,
		Message: "ButtonsPermission atualizada com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil

}
