package repositories

import (
	"context"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type ScheduledPostRepository struct {
	db *gorm.DB
}

func NewScheduledPostRepository(db *gorm.DB) *ScheduledPostRepository {
	return &ScheduledPostRepository{db: db}
}

func (r *ScheduledPostRepository) Create(ctx context.Context, post *models.ScheduledPost) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *ScheduledPostRepository) GetByID(ctx context.Context, id string) (*models.ScheduledPost, error) {
	var post models.ScheduledPost
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *ScheduledPostRepository) GetByOwner(ctx context.Context, ownerID int64) ([]models.ScheduledPost, error) {
	var posts []models.ScheduledPost
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND status IN ('pending', 'paused')", ownerID).
		Order("next_run_at ASC").
		Find(&posts).Error
	return posts, err
}

func (r *ScheduledPostRepository) GetDuePosts(ctx context.Context, now time.Time) ([]models.ScheduledPost, error) {
	var posts []models.ScheduledPost
	err := r.db.WithContext(ctx).
		Where("status = 'pending' AND next_run_at <= ?", now).
		Order("next_run_at ASC").
		Limit(10).
		Find(&posts).Error
	return posts, err
}

func (r *ScheduledPostRepository) GetQueueGroup(ctx context.Context, queueGroupID string) ([]models.ScheduledPost, error) {
	var posts []models.ScheduledPost
	err := r.db.WithContext(ctx).
		Where("queue_group_id = ? AND status IN ('pending', 'sent')", queueGroupID).
		Order("queue_position ASC").
		Find(&posts).Error
	return posts, err
}

func (r *ScheduledPostRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ScheduledPostRepository) UpdateNextRunAt(ctx context.Context, id string, nextRunAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Update("next_run_at", nextRunAt).Error
}

func (r *ScheduledPostRepository) MarkSent(ctx context.Context, id string, sentAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "sent",
			"sent_at":    sentAt,
			"sent_count": gorm.Expr("sent_count + 1"),
		}).Error
}

func (r *ScheduledPostRepository) UpdateError(ctx context.Context, id string, lastError string, retryCount int) error {
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "failed",
			"last_error":  lastError,
			"retry_count": retryCount,
		}).Error
}

func (r *ScheduledPostRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.ScheduledPost{}).Error
}

func (r *ScheduledPostRepository) UpdateScheduleTime(ctx context.Context, id string, nextRunAt time.Time, scheduleTime string) error {
	updates := map[string]interface{}{
		"next_run_at": nextRunAt,
	}
	if scheduleTime != "" {
		updates["schedule_time"] = scheduleTime
	}
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *ScheduledPostRepository) UpdatePinMessage(ctx context.Context, id string, pinMessage bool) error {
	return r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("id = ?", id).
		Update("pin_message", pinMessage).Error
}

func (r *ScheduledPostRepository) GetByOwnerAndChannel(ctx context.Context, ownerID, channelID int64) ([]models.ScheduledPost, error) {
	var posts []models.ScheduledPost
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND channel_id = ? AND status IN ('pending', 'paused')", ownerID, channelID).
		Order("next_run_at ASC").
		Find(&posts).Error
	return posts, err
}

func (r *ScheduledPostRepository) CountByOwner(ctx context.Context, ownerID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ScheduledPost{}).
		Where("owner_id = ? AND status = 'pending'", ownerID).
		Count(&count).Error
	return count, err
}
