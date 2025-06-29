package repositories

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) CountUserChannels(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Model(&models.Channel{}).Where("owner_id = ?", userID).Error
	return count, err
}
