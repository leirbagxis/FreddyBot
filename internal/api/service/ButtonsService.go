package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

func (app *AppContainerLocal) CreateButtonService(ctx context.Context, channelID int64, buttonData types.ButtonCreateRequest) (*types.ButtonCreateResponse, error) {
	if err := validateButtonData(buttonData); err != nil {
		return nil, err
	}

	position, err := app.calculateNextButtonPosition(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular posição: %w", err)
	}

	newButton := &models.Button{
		ButtonID:       uuid.NewString(),
		OwnerChannelID: channelID,
		NameButton:     buttonData.NameButton,
		ButtonURL:      buttonData.ButtonURL,
		PositionX:      position.X,
		PositionY:      position.Y,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = app.ButtonRepo.CreateButton(ctx, newButton)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar botão: %w", err)
	}

	fmt.Printf("✅ Botão padrão criado com sucesso: %s\n", newButton.ButtonID)

	return &types.ButtonCreateResponse{
		Success: true,
		Message: "Botão criado com sucesso",
		Data:    newButton,
	}, nil
}

func (app *AppContainerLocal) UpdateButtonService(ctx context.Context, channelID int64, buttonID string, buttonData types.ButtonCreateRequest) (*types.ButtonResponse, error) {
	if err := validateButtonData(buttonData); err != nil {
		return nil, err
	}

	now := time.Now()
	result := app.DB.Model(&models.Button{}).
		Where("button_id = ? AND owner_channel_id = ?", buttonID, channelID).
		Updates(map[string]interface{}{
			"name_button": buttonData.NameButton,
			"button_url":  buttonData.ButtonURL,
			"updated_at":  now,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar um botao padrão: %w", result.Error)
	}
	fmt.Printf("✅ Botão padrão atualizado com sucesso")

	return &types.ButtonResponse{
		Success: true,
		Message: "Botão padrao atualizado com sucesso",
		Data: map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"updated_at":    now,
		},
	}, nil
}

func (app *AppContainerLocal) DeleteDefaulfButtonService(ctx context.Context, channelID int64, buttonID string) error {
	if buttonID == "" {
		return fmt.Errorf("ID do botão é obrigatório")
	}

	if channelID == 0 {
		return fmt.Errorf("ID do canal é obrigatório")
	}

	result := app.DB.WithContext(ctx).
		Debug().
		Where("button_id = ? AND owner_channel_id = ?", buttonID, channelID).
		Delete(&models.Button{})

	if result.Error != nil {
		return fmt.Errorf("erro ao deletar botão: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("botão não encontrado ou não pertence ao canal")
	}

	return nil
}

func (app *AppContainerLocal) UpdateButtonsLayoutService(ctx context.Context, channelID int64, layoutData types.UpdateLayoutRequest) (*types.UpdateLayoutResponse, error) {
	if len(layoutData.Layout) == 0 {
		return nil, errors.New("layout não pode ser vazio")
	}

	var channel models.Channel
	err := app.DB.WithContext(ctx).
		Where("id = ? ", channelID).
		First(&channel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("canal não encontrado ou sem permissão")
		}
		return nil, fmt.Errorf("erro ao buscar canal: %w", err)
	}

	var buttonsToUpdate []struct {
		ID        string
		PositionX int
		PositionY int
	}
	for y, row := range layoutData.Layout {
		for x, button := range row {
			buttonsToUpdate = append(buttonsToUpdate, struct {
				ID        string
				PositionX int
				PositionY int
			}{
				ID:        button.ID,
				PositionX: x,
				PositionY: y,
			})
		}
	}

	now := time.Now()
	err = app.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, btn := range buttonsToUpdate {
			result := tx.Model(&models.Button{}).
				Where("button_id = ? AND owner_channel_id = ?", btn.ID, channelID).
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

	fmt.Println("✅ Layout dos botões atualizado com sucesso")

	return &types.UpdateLayoutResponse{
		Success: true,
		Message: "Layout dos botões atualizado com sucesso",
		Data: map[string]interface{}{
			"updated_at": now,
			"total":      len(buttonsToUpdate),
		},
	}, nil
}

func (app *AppContainerLocal) calculateNextButtonPosition(ctx context.Context, channelId int64) (*types.ButtonPosition, error) {
	buttons, err := app.ChannelRepo.GetChannelButtons(ctx, channelId)
	if err != nil {
		return nil, err
	}

	if len(buttons) == 0 {
		return &types.ButtonPosition{X: 0, Y: 0}, nil
	}

	// Encontrar última linha (maxY)
	maxY := 0
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

	// Se última linha tem menos de 3 botões, adicionar na mesma linha
	if buttonsInLastRow < 3 {
		return &types.ButtonPosition{X: buttonsInLastRow, Y: maxY}, nil
	}

	// Criar nova linha
	return &types.ButtonPosition{X: 0, Y: maxY + 1}, nil
}

func validateButtonData(data types.ButtonCreateRequest) error {
	// Validar nome do botão
	if strings.TrimSpace(data.NameButton) == "" {
		return fmt.Errorf("texto do botão é obrigatório")
	}

	if len(data.NameButton) > 64 {
		return fmt.Errorf("texto do botão muito longo (máximo 64 caracteres)")
	}

	// Validar URL se fornecida
	if data.ButtonURL != "" && !isValidURL(data.ButtonURL) {
		return fmt.Errorf("URL inválida")
	}

	return nil
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}
