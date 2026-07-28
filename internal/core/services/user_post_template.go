package services

import (
	"context"
	"fmt"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
)

type UserPostTemplateService struct {
	repo *repositories.UserPostTemplateRepository
}

func NewUserPostTemplateService(repo *repositories.UserPostTemplateRepository) *UserPostTemplateService {
	return &UserPostTemplateService{repo: repo}
}

func (s *UserPostTemplateService) SaveTemplate(ctx context.Context, ownerID int64, name, templateData string) (*models.UserPostTemplate, error) {
	template := &models.UserPostTemplate{
		ID:           fmt.Sprintf("tpl_%d_%d", time.Now().UnixNano(), ownerID),
		OwnerID:      ownerID,
		Name:         name,
		TemplateData: templateData,
	}

	if err := s.repo.Create(ctx, template); err != nil {
		return nil, errors.Internal(err)
	}
	return template, nil
}

func (s *UserPostTemplateService) ListTemplates(ctx context.Context, ownerID int64) ([]models.UserPostTemplate, error) {
	return s.repo.GetByOwner(ctx, ownerID)
}

func (s *UserPostTemplateService) DeleteTemplate(ctx context.Context, id string, ownerID int64) error {
	return s.repo.Delete(ctx, id, ownerID)
}

func (s *UserPostTemplateService) GetTemplateByID(ctx context.Context, id string) (*models.UserPostTemplate, error) {
	tpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return tpl, nil
}
