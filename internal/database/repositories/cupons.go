package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (r *CouponRepository) ValidateCoupon(ctx context.Context, code string, userID int64) (*models.Coupon, error) {

	var coupon models.Coupon
	if err := r.db.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, errors.New("Este cupom é inválido")
	}

	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Este cupom está expirado")
	}

	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return nil, errors.New("Este cupom está esgotado")
	}

	if coupon.OnePerUser {
		key := fmt.Sprintf("%s:%d", coupon.ID, userID)
		var count int64
		r.db.Model(&models.CouponUsage{}).
			Where("unique_key = ?", key).
			Count(&count)

		if count > 0 {
			return nil, errors.New("Você já utilizou esse cupom")
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
	if err := r.db.WithContext(ctx).Where("is_auto = true").First(&coupon).Error; err != nil {
		return nil, nil
	}

	return &coupon, nil
}

func (r *CouponRepository) ApplyCouponToPayment(
	ctx context.Context,
	coupon *models.Coupon,
	payment *models.Payment,
	userID int64,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 🔒 Lock no payment (anti race)
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&payment, "id = ?", payment.ID).Error; err != nil {

			return err
		}

		// 🚫 Não permitir trocar cupom
		fmt.Println(payment.CouponID)
		if payment.CouponID != nil {
			if strings.TrimSpace(*payment.CouponID) != "" {
				return errors.New("este pagamento já possui um cupom aplicado")
			}
		}

		// 🧮 Calcular valor final
		finalAmount := payment.Amount

		switch coupon.DiscountType {
		case "percent":
			finalAmount = payment.Amount - (payment.Amount*coupon.DiscountValue)/100
		case "fixed":
			finalAmount = payment.Amount - coupon.DiscountValue
		}

		if finalAmount < 1 {
			finalAmount = 1
		}

		// 🧾 Registrar uso do cupom
		usage := models.CouponUsage{
			ID:       uuid.NewString(),
			CouponID: coupon.ID,
			UserID:   userID,
			UsedAt:   time.Now(),
		}

		if coupon.OnePerUser {
			usage.UniqueKey = fmt.Sprintf("%s:%d", coupon.ID, userID)
		} else {
			usage.UniqueKey = uuid.NewString()
		}

		if err := tx.Create(&usage).Error; err != nil {
			return err
		}

		// ➕ Incrementar contador global
		if err := tx.Model(&models.Coupon{}).
			Where("id = ?", coupon.ID).
			Update("used_count", gorm.Expr("used_count + 1")).Error; err != nil {

			return err
		}

		// 💳 Atualizar payment
		if err := tx.Model(&models.Payment{}).
			Where("id = ?", payment.ID).
			Updates(map[string]interface{}{
				"amount":    finalAmount,
				"coupon_id": coupon.ID,
			}).Error; err != nil {

			return err
		}

		return nil
	})
}
