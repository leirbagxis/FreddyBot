package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetPlanByID(ctx context.Context, planID string) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := r.db.WithContext(ctx).First(&plan, "id = ?", planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscription plan %q not found", planID)
		}
		return nil, err
	}
	return &plan, nil
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

func (r *SubscriptionRepository) AssignSubscriptionToChannel(ctx context.Context, channelId int64, plan models.SubscriptionPlan, paymentMethod string, appliedCouponID *string, redeemedGiftID *string) (*models.Subscription, error) {

	var result *models.Subscription

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&models.Subscription{}).
			Where("channel_id = ? AND status = ?", channelId, "active").
			Updates(map[string]any{
				"status":     "expired",
				"end_date":   gorm.Expr("COALESCE(end_date, ?)", time.Now().UTC()),
				"updated_at": time.Now().UTC(),
			}).Error; err != nil {
			return err
		}

		now := time.Now().UTC()
		start := now
		var end time.Time
		if plan.Duration > 0 {
			end = start.AddDate(0, 0, plan.Duration)
		} else {
			end = time.Time{}
		}

		sub := models.Subscription{
			ChannelID:       channelId,
			PlanID:          toPlanID(plan.ID),
			Status:          "active",
			StartDate:       start,
			EndDate:         end,
			PaymentMethod:   paymentMethod,
			AppliedCouponID: appliedCouponID,
			RedeemedGiftID:  redeemedGiftID,
		}

		if err := tx.Create(&sub).Error; err != nil {
			return err
		}

		if err := tx.Preload("Plan").First(&sub, "id = ?", sub.ID).Error; err != nil {
			return err
		}

		result = &sub
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func toPlanID(id any) string {
	switch v := id.(type) {
	case string:
		return v
	case uint:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
