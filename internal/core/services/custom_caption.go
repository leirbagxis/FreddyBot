package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

type CustomCaptionService struct {
	customCaptionRepo *repositories.CustomCaptionRepository
	channelRepo       *repositories.ChannelRepository
	cache             *cache.Service
}

func NewCustomCaptionService(customCaptionRepo *repositories.CustomCaptionRepository, channelRepo *repositories.ChannelRepository, cache *cache.Service) *CustomCaptionService {
	return &CustomCaptionService{
		customCaptionRepo: customCaptionRepo,
		channelRepo:       channelRepo,
		cache:             cache,
	}
}

func (s *CustomCaptionService) CreateCustomCaption(ctx context.Context, channelID int64, body types.CreateCustomCaptionRequest) (*models.CustomCaption, error) {
	_, err := s.channelRepo.GetChannelByIDLight(ctx, channelID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	newCaption := &models.CustomCaption{
		CaptionID:      uuid.NewString(),
		Code:           body.Code,
		Caption:        body.Caption,
		LinkPreview:    body.LinkPreview,
		OwnerChannelID: channelID,
	}

	if err := s.customCaptionRepo.CreateCustomCaption(ctx, newCaption); err != nil {
		return nil, errors.Internal(err)
	}

	s.cache.InvalidateChannel(ctx, channelID)
	return newCaption, nil
}

func (s *CustomCaptionService) UpdateCustomCaption(ctx context.Context, channelID int64, captionID string, body types.CreateCustomCaptionRequest) (int64, error) {
	updates := map[string]interface{}{
		"caption":      body.Caption,
		"link_preview": body.LinkPreview,
		"updated_at":   time.Now(),
	}

	rowsAffected, err := s.customCaptionRepo.UpdateCustomCaption(ctx, channelID, captionID, updates)
	if err != nil {
		return 0, errors.Internal(err)
	}

	if rowsAffected == 0 {
		return 0, errors.ErrNotFound
	}

	s.cache.InvalidateChannel(ctx, channelID)
	logger.Bot("✅ Legenda customizada atualizada com sucesso: %s (Canal: %d)", captionID, channelID)

	return rowsAffected, nil
}

func (s *CustomCaptionService) DeleteCustomCaption(ctx context.Context, channelID int64, captionID string) error {
	rowsAffected, err := s.customCaptionRepo.DeleteCustomCaption(ctx, channelID, captionID)
	if err != nil {
		return errors.Internal(err)
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	s.cache.InvalidateChannel(ctx, channelID)
	return nil
}
