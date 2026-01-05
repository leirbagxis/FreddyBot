package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{
		db: db,
	}
}

// Validação do cupom (função única)
func (r *CouponRepository) ValidateCoupon(
	ctx context.Context,
	code string,
	userID int64,
) (*models.Coupon, error) {

	var coupon models.Coupon
	if err := r.db.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, errors.New("cupom inválido")
	}

	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("cupom expirado")
	}

	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return nil, errors.New("cupom esgotado")
	}

	if coupon.OnePerUser {
		key := fmt.Sprintf("%s:%d", coupon.ID, userID)
		var count int64
		r.db.Model(&models.CouponUsage{}).
			Where("unique_key = ?", key).
			Count(&count)

		if count > 0 {
			return nil, errors.New("cupom já utilizado")
		}
	}

	return &coupon, nil
}

// Cupom automático (primeira compra)
func (r *CouponRepository) GetAutoCoupon(ctx context.Context, userID int64) (*models.Coupon, error) {
	var count int64
	r.db.Model(&models.Payment{}).
		Where("user_id = ?", userID).
		Count(&count)

	if count > 0 {
		return nil, nil
	}

	var coupon models.Coupon
	if err := r.db.Where("is_auto = true").First(&coupon).Error; err != nil {
		return nil, nil
	}

	return &coupon, nil
}
