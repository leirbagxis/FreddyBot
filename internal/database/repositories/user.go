package repositories

import (
	"context"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(ctx context.Context, user *models.User) error {
	var existing models.User
	err := r.db.WithContext(ctx).First(&existing, "user_id = ?", user.UserId).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	now := time.Now()

	if err == gorm.ErrRecordNotFound {
		user.CreatedAt = now
		user.UpdatedAt = now
		return r.db.WithContext(ctx).Create(user).Error
	}

	return r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
		"first_name": user.FirstName,
		"updated_at": now,
	}).Error
}

func (r *UserRepository) GetUserById(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "user_id = ? ", userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
