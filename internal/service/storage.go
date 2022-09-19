package service

import "context"

type Repository interface {
	FindOne(ctx context.Context, id string) (*PaidService, error)
	FindAll(ctx context.Context) ([]PaidService, error)
	Create(ctx context.Context, s *PaidService) error
	Update(ctx context.Context, s PaidService) error
	Delete(ctx context.Context, id string) error
}
