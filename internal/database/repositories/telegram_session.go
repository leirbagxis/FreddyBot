package repositories

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type TelegramSessionRepository struct {
	db *gorm.DB
}

func NewTelegramSessionRepository(db *gorm.DB) *TelegramSessionRepository {
	return &TelegramSessionRepository{db: db}
}

func (r *TelegramSessionRepository) Upsert(ctx context.Context, session *models.UserTelegramSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *TelegramSessionRepository) GetByUserID(ctx context.Context, userID int64) (*models.UserTelegramSession, error) {
	var session models.UserTelegramSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&session).Error
	return &session, err
}

func (r *TelegramSessionRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserTelegramSession{}).Error
}

func (r *TelegramSessionRepository) SetActive(ctx context.Context, userID int64, active bool) error {
	return r.db.WithContext(ctx).Model(&models.UserTelegramSession{}).
		Where("user_id = ?", userID).
		Update("is_active", active).Error
}

func (r *TelegramSessionRepository) GetActiveUserIDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Model(&models.UserTelegramSession{}).
		Where("is_active = ?", true).
		Pluck("user_id", &ids).Error
	return ids, err
}
