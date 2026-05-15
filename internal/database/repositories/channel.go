package repositories

import (
	"context"
	"errors"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) CountUserChannels(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Channel{}).
		Where("owner_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *ChannelRepository) GetChannelByTwoID(ctx context.Context, userId, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.owner_id = ? AND channels.id = ?", userId, channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelRepository) GetChannelByUserID(ctx context.Context, userId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.owner_id = ?", userId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelRepository) GetChannelByID(ctx context.Context, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Joins("DefaultCaption").
		Joins("DefaultCaption.MessagePermission").
		Joins("DefaultCaption.ButtonsPermission").
		Joins("Separator").
		Preload("Owner").
		Preload("Buttons").
		Preload("CustomCaptions").
		Preload("CustomCaptions.Buttons").
		Where("channels.id = ?", channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelRepository) GetChannelByIDLight(ctx context.Context, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Where("id = ?", channelId).
		First(&channel).Error
	return &channel, err
}

func (r *ChannelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	return r.db.WithContext(ctx).Create(channel).Error
}

func (r *ChannelRepository) UpdateChannel(ctx context.Context, channel *models.Channel) error {
	return r.db.WithContext(ctx).Save(channel).Error
}

func (r *ChannelRepository) UpdateOwnerChannel(ctx context.Context, channelID, oldOwnerID, newOwnerID int64) error {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", channelID, oldOwnerID).
		First(&channel).Error
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&channel).
		Updates(map[string]any{
			"owner_id":      newOwnerID,
			"token_version": gorm.Expr("token_version + 1"),
		}).Error
}

func (r *ChannelRepository) DeleteChannelWithRelations(ctx context.Context, userId, channelId int64) error {
	result := r.db.WithContext(ctx).
		Where("owner_id = ? AND id = ?", userId, channelId).
		Delete(&models.Channel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("channel not found")
	}
	return nil
}

func (r *ChannelRepository) GetAllChannelsByUserID(ctx context.Context, userID int64) ([]models.Channel, error) {
	var channels []models.Channel
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Order("updated_at ASC").
		Find(&channels).Error
	return channels, err
}

func (r *ChannelRepository) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	var channels []models.Channel
	err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		Find(&channels).Error
	return channels, err
}

func (r *ChannelRepository) GetAllChannelsPaginated(ctx context.Context, limit, offset int) ([]models.Channel, int64, error) {
	var channels []models.Channel
	var total int64
	db := r.db.WithContext(ctx).Model(&models.Channel{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Limit(limit).Offset(offset).Order("updated_at DESC").Find(&channels).Error
	return channels, total, err
}

func (r *ChannelRepository) GetChannelButtons(ctx context.Context, channelId int64) ([]models.Button, error) {
	var buttons []models.Button
	err := r.db.WithContext(ctx).
		Where("owner_channel_id = ?", channelId).
		Order("position_y ASC, position_x ASC").
		Find(&buttons).Error
	return buttons, err
}

func (r *ChannelRepository) UpdateDefaultCaption(ctx context.Context, channelID int64, caption string) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.DefaultCaption{}).
		Where("owner_channel_id = ?", channelID).
		Update("caption", caption)
	return result.RowsAffected, result.Error
}

func (r *ChannelRepository) UpdateNewPackCaption(ctx context.Context, channelID int64, caption string) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.Channel{}).
		Where("id = ?", channelID).
		Update("new_pack_caption", caption)
	return result.RowsAffected, result.Error
}

func (r *ChannelRepository) UpdateReactions(ctx context.Context, channelID int64, reactions string) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.Channel{}).
		Where("id = ?", channelID).
		Update("reactions", reactions)
	return result.RowsAffected, result.Error
}

func (r *ChannelRepository) UpdateReactionPosition(ctx context.Context, channelID int64, position int) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.Channel{}).
		Where("id = ?", channelID).
		Update("reaction_position", position)
	return result.RowsAffected, result.Error
}
