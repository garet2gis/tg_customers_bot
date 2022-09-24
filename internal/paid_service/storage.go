package paid_service

import (
	"context"
	"errors"
)

var NoRowsDeleted = errors.New("No row found to delete")
var CanNotDeleteRowForeignKey = errors.New("Can not delete row with foreign key")

type Repository interface {
	FindOne(ctx context.Context, id string) (*PaidService, error)
	FindAll(ctx context.Context) ([]PaidService, error)
	Create(ctx context.Context, s *PaidService) error
	Update(ctx context.Context, s PaidService) error
	Delete(ctx context.Context, id int) error
}

type ServiceTemporaryRepository interface {
	Update(userID int64, service *CreatePaidServiceDTO, bucket string) error
	Get(userID int64, bucket string) (*CreatePaidServiceDTO, error)
	Delete(userID int64, bucket string) error
}
