package repositories

import (
	"context"
	"errors"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetChannelSubscription(ctx context.Context, channelId int64) (*models.Subscription, error) {
	var sub models.Subscription

	err := r.db.WithContext(ctx).Where("channel_id = ? AND status = ?", channelId, "active").Order("end_date DESC NULLS LAST, updated_at DESC").
		Preload("Plan").
		First(&sub).Error

	if err == nil {
		return &sub, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2) Caso não exista ativa, pegar a mais recente (qualquer status)
	err = r.db.WithContext(ctx).
		Where("channel_id = ?", channelId).
		Order("end_date DESC NULLS LAST, updated_at DESC").
		Preload("Plan").
		First(&sub).Error

	if err != nil {
		return nil, err
	}
	return &sub, nil

}
