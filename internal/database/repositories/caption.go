package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type CaptionRepository struct {
	db *gorm.DB
}

func NewCaptionRepository(db *gorm.DB) *CaptionRepository {
	return &CaptionRepository{db: db}
}

func (r *CaptionRepository) SavePremiumCaption(
	ctx context.Context,
	channelId int64,
	text string,
	entities []models.TGMessageEntity,
) (*models.DefaultCaption, error) {

	// Serializa entities para JSON
	var entitiesJSON []byte
	var err error
	if len(entities) > 0 {
		entitiesJSON, err = json.Marshal(entities)
		if err != nil {
			return nil, fmt.Errorf("marshal entities: %w", err)
		}
	}

	var dc models.DefaultCaption
	// Busca por OwnerChannelID único
	tx := r.db.WithContext(ctx).Where("owner_channel_id = ?", channelId).First(&dc)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			// Criar
			dc = models.DefaultCaption{
				Caption:             text,
				CaptionEntitiesJSON: entitiesJSON,
				OwnerChannelID:      channelId,
			}
			if err := r.db.WithContext(ctx).Create(&dc).Error; err != nil {
				return nil, err
			}
			return &dc, nil
		}
		return nil, tx.Error
	}

	// Atualizar
	dc.Caption = text
	dc.CaptionEntitiesJSON = entitiesJSON
	if err := r.db.WithContext(ctx).Save(&dc).Error; err != nil {
		return nil, err
	}
	return &dc, nil
}
