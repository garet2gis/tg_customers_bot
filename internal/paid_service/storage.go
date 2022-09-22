package paid_service

import (
	"context"
)

type Repository interface {
	FindOne(ctx context.Context, id string) (*PaidService, error)
	FindAll(ctx context.Context) ([]PaidService, error)
	Create(ctx context.Context, s *PaidService) error
	Update(ctx context.Context, s PaidService) error
	Delete(ctx context.Context, id string) error
}

type ServiceTemporaryRepository interface {
	Update(userID int64, service *CreatePaidServiceDTO, bucket string) error
	Get(userID int64, bucket string) (*CreatePaidServiceDTO, error)
	Delete(userID int64, bucket string) error
}
