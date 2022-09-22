package paid_service

import (
	"context"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) FindOne(ctx context.Context, id string) (*PaidService, error) {
	one, err := s.repository.FindOne(ctx, id)
	if err != err {
		return nil, err
	}
	return one, nil
}

func (s *Service) FindAll(ctx context.Context) ([]PaidService, error) {
	all, err := s.repository.FindAll(ctx)
	if err != err {
		return nil, err
	}
	return all, nil
}

func (s *Service) Create(ctx context.Context, ps *PaidService) (string, error) {
	err := s.repository.Create(ctx, ps)
	if err != err {
		return "", err
	}
	return ps.ID, nil
}
