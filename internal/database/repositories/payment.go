package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) GetPricePlan(ctx context.Context, key string) (*models.PriceConfig, error) {
	var cfg models.PriceConfig

	err := r.db.WithContext(ctx).
		Where("key = ? AND active = true", key).
		First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("preço não encontrado")
		}
		return nil, err
	}

	return &cfg, nil
}

func (r *PaymentRepository) CreateNewPayment(
	ctx context.Context,
	payload models.Payment,
) (*models.Payment, error) {

	p := models.Payment{
		ID:        uuid.NewString(),
		UserID:    payload.UserID,
		Amount:    payload.Amount,
		Payload:   payload.Payload,
		CreatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}
