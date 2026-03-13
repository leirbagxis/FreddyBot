package repositories

import (
	"context"
	"errors"
	"fmt"
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

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "username = ? ", username).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUSers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Order("updated_at ASC").Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUserAdmin(ctx context.Context, userID int64) (bool, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("usuário não encontrado")
		}
		return false, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	newValue := !user.IsAdmin

	err = r.db.WithContext(ctx).
		Model(&user).
		Updates(map[string]any{
			"is_admin": newValue,
		}).Error
	if err != nil {
		return false, fmt.Errorf("erro ao atualizar status de admin do usuário: %w", err)
	}

	return newValue, nil
}
