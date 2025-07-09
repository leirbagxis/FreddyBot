package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
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
