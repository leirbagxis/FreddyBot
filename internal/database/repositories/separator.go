package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeparatorRepository struct {
	db *gorm.DB
}

func NewSeparatorRepository(db *gorm.DB) *SeparatorRepository {
	return &SeparatorRepository{db: db}
}

func (r *SeparatorRepository) GetSeparatorByOwnerChannelID(ctx context.Context, ownerChannelID int64) (*models.Separator, error) {
	var separator models.Separator

	err := r.db.WithContext(ctx).
		Where("owner_channel_id = ?", ownerChannelID).
		First(&separator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &separator, nil
}

func (r *SeparatorRepository) SaveSeparator(ctx context.Context, separator *models.Separator) error {
	if separator.ID == "" {
		separator.ID = uuid.NewString()
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "owner_channel_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"separator_id", "separator_url", "updated_at"}),
		}).
		Create(separator).Error

	return err
}

func (r *SeparatorRepository) DeleteSeparatorByOwnerChannelId(ctx context.Context, channelID int64) error {
	result := r.db.WithContext(ctx).
		Where("owner_channel_id = ?", channelID).
		Delete(&models.Separator{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("separator not found")
	}

	return nil
}

func (r *SeparatorRepository) GetSeparatorByTwoID(ctx context.Context, channelID int64, separatorID string) (*models.Separator, error) {
	var separator models.Separator

	err := r.db.WithContext(ctx).
		Where("owner_channel_id = ? and separator_id = ?", channelID, separatorID).
		First(&separator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &separator, nil
}

func (r *SeparatorRepository) SavePremiumSeparator(ctx context.Context, channelId int64, separator string, entities []models.TGMessageEntity) (*models.Separator, error) {
	var entsJSON []byte
	if len(entities) > 0 {
		b, err := json.Marshal(entities)
		if err != nil {
			return nil, fmt.Errorf("marshal entities: %w", err)
		}
		entsJSON = b
	}

	var sep models.Separator
	tx := r.db.WithContext(ctx).Where("owner_channel_id = ?", channelId).First(&sep)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			sep = models.Separator{
				OwnerChannelID:        channelId,
				SeparatorText:         separator,
				SeparatorEntitiesJSON: entsJSON,
				// limpa modo arquivo
				// SeparatorID:  "",
				// SeparatorURL: "",
			}
			if err := r.db.WithContext(ctx).Create(&sep).Error; err != nil {
				return nil, err
			}
			return &sep, nil
		}
		return nil, tx.Error
	}

	// atualizar para modo texto
	sep.SeparatorText = separator
	sep.SeparatorEntitiesJSON = entsJSON
	sep.SeparatorID = ""
	sep.SeparatorURL = ""
	if err := r.db.WithContext(ctx).Save(&sep).Error; err != nil {
		return nil, err
	}
	return &sep, nil
}
