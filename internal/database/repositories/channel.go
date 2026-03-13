package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) CountUserChannels(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Channel{}).
		Where("owner_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *ChannelRepository) GetChannelByTwoID(ctx context.Context, userId, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		// Usar Joins para relações 1:1 (melhor performance)
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		// Usar Preload para relações 1:N
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.owner_id = ? AND channels.id = ?", userId, channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil

}

func (r *ChannelRepository) GetChannelByUserID(ctx context.Context, userId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		// Usar Joins para relações 1:1 (melhor performance)
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		// Usar Preload para relações 1:N
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.owner_id = ?", userId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil

}

func (r *ChannelRepository) GetChannelByID(ctx context.Context, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		// Usar Joins para relações 1:1 (melhor performance)
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		// Usar Preload para relações 1:N
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.id = ?", channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil

}

func (r *ChannelRepository) DeleteChannelByTwoId(ctx context.Context, userId, channelId int64) error {
	result := r.db.WithContext(ctx).
		Where("owner_id = ? AND id = ?", userId, channelId).
		Delete(&models.Channel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("channel not found or you don't have permission to delete it")
	}

	return nil

}

func (r *ChannelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	return r.db.WithContext(ctx).Create(channel).Error
}

func (r *ChannelRepository) CreateChannelWithDefaults(ctx context.Context, channelID int64, title, inviteURL, newPackCaption, caption string, ownerID int64) (*models.Channel, error) {
	channel := &models.Channel{
		ID:             channelID,
		Title:          title,
		NewPackCaption: newPackCaption,
		InviteURL:      inviteURL,
		OwnerID:        ownerID,
		DefaultCaption: &models.DefaultCaption{
			CaptionID:      uuid.New().String(),
			Caption:        caption,
			OwnerChannelID: channelID,
			MessagePermission: &models.MessagePermission{
				MessagePermissionID: uuid.New().String(),
				LinkPreview:         true,
				Message:             true,
				Audio:               true,
				Video:               true,
				Photo:               true,
				Sticker:             true,
				GIF:                 true,
			},
			ButtonsPermission: &models.ButtonsPermission{
				ButtonsPermissionID: uuid.New().String(),
				Message:             true,
				Audio:               true,
				Video:               true,
				Photo:               true,
				Sticker:             true,
				GIF:                 true,
			},
		},
		Buttons: []models.Button{
			{
				ButtonID:       uuid.NewString(),
				NameButton:     title,
				ButtonURL:      inviteURL,
				PositionX:      0,
				PositionY:      0,
				OwnerChannelID: channelID,
			},
		},
	}

	captionID := channel.DefaultCaption.CaptionID
	channel.DefaultCaption.MessagePermission.OwnerCaptionID = captionID
	channel.DefaultCaption.ButtonsPermission.OwnerCaptionID = captionID

	err := r.db.WithContext(ctx).Create(channel).Error
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) DeleteChannelWithRelations(ctx context.Context, userId, channelId int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Buscar o canal com todas as relações
		var channel models.Channel
		err := tx.Preload("DefaultCaption").
			Preload("DefaultCaption.MessagePermission").
			Preload("DefaultCaption.ButtonsPermission").
			Preload("Buttons").
			Preload("Separator").
			Preload("CustomCaptions").
			Preload("CustomCaptions.Buttons").
			Where("owner_id = ? AND id = ?", userId, channelId).
			First(&channel).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("channel not found or you don't have permission to delete it")
			}
			return fmt.Errorf("failed to find channel: %w", err)
		}

		// 2. Deletar CustomCaptionButtons
		for _, customCaption := range channel.CustomCaptions {
			for _, button := range customCaption.Buttons {
				if err := tx.Delete(&button).Error; err != nil {
					return fmt.Errorf("failed to delete custom caption button: %w", err)
				}
			}
		}

		// 3. Deletar CustomCaptions
		for _, customCaption := range channel.CustomCaptions {
			if err := tx.Delete(&customCaption).Error; err != nil {
				return fmt.Errorf("failed to delete custom caption: %w", err)
			}
		}

		// 4. Deletar Buttons
		for _, button := range channel.Buttons {
			if err := tx.Delete(&button).Error; err != nil {
				return fmt.Errorf("failed to delete button: %w", err)
			}
		}

		// 5. Deletar Separator
		if channel.Separator != nil {
			if err := tx.Delete(channel.Separator).Error; err != nil {
				return fmt.Errorf("failed to delete separator: %w", err)
			}
		}

		// 6. Deletar MessagePermission e ButtonsPermission
		if channel.DefaultCaption != nil {
			if channel.DefaultCaption.MessagePermission != nil {
				if err := tx.Delete(channel.DefaultCaption.MessagePermission).Error; err != nil {
					return fmt.Errorf("failed to delete message permission: %w", err)
				}
			}

			if channel.DefaultCaption.ButtonsPermission != nil {
				if err := tx.Delete(channel.DefaultCaption.ButtonsPermission).Error; err != nil {
					return fmt.Errorf("failed to delete buttons permission: %w", err)
				}
			}

			// 7. Deletar DefaultCaption
			if err := tx.Delete(channel.DefaultCaption).Error; err != nil {
				return fmt.Errorf("failed to delete default caption: %w", err)
			}
		}

		// 8. Finalmente, deletar o Channel
		if err := tx.Delete(&channel).Error; err != nil {
			return fmt.Errorf("failed to delete channel: %w", err)
		}

		return nil
	})
}

func (r *ChannelRepository) GetChannelWithRelations(ctx context.Context, channelId int64) (*models.Channel, error) {
	var channel models.Channel

	err := r.db.WithContext(ctx).
		// Usar Joins para relações 1:1 (melhor performance)
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		// Usar Preload para relações 1:N
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.id = ?", channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *ChannelRepository) UpdateOwnerChannel(ctx context.Context, channelID, oldOwnerID, newOwnerID int64) error {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", channelID, oldOwnerID).
		First(&channel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("canal não encontrado ou você não tem permissão para modificá-lo")
		}
		return fmt.Errorf("Ërro ao buscar canal %w", err)
	}

	err = r.db.WithContext(ctx).
		Model(&channel).
		Updates(map[string]any{
			"owner_id":      newOwnerID,
			"token_version": gorm.Expr("token_version + 1"),
		}).Error
	//Updates("owner_id", newOwnerID).Error

	if err != nil {
		return fmt.Errorf("Erro ao atualizar proprietario do canal: %w", err)
	}

	return nil

}

func (r *ChannelRepository) GetAllChannelsByUserID(ctx context.Context, userID int64) ([]models.Channel, error) {
	var channel []models.Channel
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Find(&channel).
		Order("updated_at ASC").Error

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	var channel []models.Channel
	err := r.db.WithContext(ctx).
		Find(&channel).Error

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) GetChannelButtons(ctx context.Context, channelId int64) ([]models.Button, error) {
	var buttons []models.Button

	err := r.db.WithContext(ctx).
		Where("owner_channel_id = ?", channelId).
		Order("position_y ASC, position_x ASC").
		Find(&buttons).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar botões: %w", err)
	}

	return buttons, nil
}

func (r *ChannelRepository) UpdateChannelBasicInfo(ctx context.Context, channelID int64, title, inviteURL string) error {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Where("id = ?", channelID).
		First(&channel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("canal não encontrado ou você não tem permissão para modificá-lo")
		}
		return fmt.Errorf("Ërro ao buscar canal %w", err)
	}

	now := time.Now()
	err = r.db.WithContext(ctx).Model(&channel).Updates(map[string]interface{}{
		"title":      utils.RemoveHTMLTags(title),
		"invite_url": inviteURL,
		"updated_at": now,
	}).Error

	if err != nil {
		return fmt.Errorf("Erro ao atualizar basic info do canal: %w", err)
	}

	return nil
}

// Função integrada para atualizar informações básicas do canal E o primeiro botão
func (r *ChannelRepository) UpdateChannelBasicInfoAndFirstButton(ctx context.Context, channel *models.Channel) error {
	// Usar transação para garantir atomicidade
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Atualizar informações básicas do canal
	result := tx.Model(&models.Channel{}).
		Where("id = ?", channel.ID).
		Updates(map[string]interface{}{
			"title":      channel.Title,
			"invite_url": channel.InviteURL,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar informações básicas do canal: %w", result.Error)
	}

	// 2. Atualizar o primeiro botão se existir
	if len(channel.Buttons) > 0 {
		firstButton := channel.Buttons[0]

		result = tx.Model(&models.Button{}).
			Where("button_id = ?", firstButton.ButtonID).
			Updates(map[string]interface{}{
				"name_button": firstButton.NameButton,
				"button_url":  firstButton.ButtonURL,
				"updated_at":  time.Now(),
			})

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao atualizar primeiro botão: %w", result.Error)
		}

		if result.RowsAffected > 0 {
			log.Printf("🔘 Primeiro botão do canal %d atualizado no banco", channel.ID)
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("erro ao fazer commit da transação: %w", err)
	}

	log.Printf("✅ Canal %d: informações básicas e primeiro botão atualizados no banco", channel.ID)
	return nil
}
