package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
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
		Where("owner_id = ? AND id = ?", userId, channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil

}

func (r *ChannelRepository) GetChannelByID(ctx context.Context, channelId int64) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.WithContext(ctx).
		Where("id = ?", channelId).
		First(&channel).Error

	if err != nil {
		return nil, err
	}

	return &channel, nil

}

func (r *ChannelRepository) DeleteChannelByTwoId(ctx context.Context, userId, channelId int64) error {
	result := r.db.WithContext(ctx).
		Where("owner_id = ? AND id = ?", userId, channelId).
		Delete(&models.Channel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("channel not found or you don't have permission to delete it")
	}

	return nil

}

func (r *ChannelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	return r.db.WithContext(ctx).Create(channel).Error
}

func (r *ChannelRepository) CreateChannelWithDefaults(ctx context.Context, channelID int64, title, inviteURL, newPackCaption, caption string, ownerID int64) (*models.Channel, error) {
	channel := &models.Channel{
		ID:             channelID,
		Title:          title,
		NewPackCaption: newPackCaption,
		InviteURL:      inviteURL,
		OwnerID:        ownerID,
		DefaultCaption: &models.DefaultCaption{
			CaptionID:      uuid.New().String(),
			Caption:        caption,
			OwnerChannelID: channelID,
			MessagePermission: &models.MessagePermission{
				MessagePermissionID: uuid.New().String(),
				LinkPreview:         true,
				Message:             true,
				Audio:               true,
				Video:               true,
				Photo:               true,
				Sticker:             true,
				GIF:                 true,
			},
			ButtonsPermission: &models.ButtonsPermission{
				ButtonsPermissionID: uuid.New().String(),
				Message:             true,
				Audio:               true,
				Video:               true,
				Photo:               true,
				Sticker:             true,
				GIF:                 true,
			},
		},
		Buttons: []models.Button{
			{
				ButtonID:       uuid.NewString(),
				NameButton:     title,
				ButtonURL:      inviteURL,
				PositionX:      0,
				PositionY:      0,
				OwnerChannelID: channelID,
			},
		},
	}

	captionID := channel.DefaultCaption.CaptionID
	channel.DefaultCaption.MessagePermission.OwnerCaptionID = captionID
	channel.DefaultCaption.ButtonsPermission.OwnerCaptionID = captionID

	err := r.db.WithContext(ctx).Create(channel).Error
	if err != nil {
		return nil, err
	}

	return channel, nil
}
