package services

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
)

type ServerService struct {
	serverRepo *repositories.ServerConfigRepository
}

func NewServerService(serverRepo *repositories.ServerConfigRepository) *ServerService {
	return &ServerService{serverRepo: serverRepo}
}

func (s *ServerService) GetConfig(ctx context.Context) (*models.ServerConfig, error) {
	config, err := s.serverRepo.GetServerConfig(ctx)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return config, nil
}

func (s *ServerService) UpdateConfig(ctx context.Context, maintenance, forceJoin bool, globalDefaultCaption, globalNewPackCaption string) (*models.ServerConfig, error) {
	config, err := s.serverRepo.GetServerConfig(ctx)
	if err != nil {
		return nil, errors.Internal(err)
	}

	config.Maintence = maintenance
	config.ForceJoin = forceJoin
	config.GlobalDefaultCaption = globalDefaultCaption
	config.GlobalNewPackCaption = globalNewPackCaption

	if err := s.serverRepo.UpdateServerConfig(ctx, config); err != nil {
		return nil, errors.Internal(err)
	}
	return config, nil
}

func (s *ServerService) GetMaintenance(ctx context.Context) (bool, error) {
	config, err := s.serverRepo.GetServerConfig(ctx)
	if err != nil {
		return false, errors.Internal(err)
	}
	return config.Maintence, nil
}

func (s *ServerService) ToggleMaintenance(ctx context.Context) (bool, error) {
	config, err := s.serverRepo.GetServerConfig(ctx)
	if err != nil {
		return false, errors.Internal(err)
	}

	config.Maintence = !config.Maintence
	if err := s.serverRepo.UpdateServerConfig(ctx, config); err != nil {
		return false, errors.Internal(err)
	}
	return config.Maintence, nil
}
