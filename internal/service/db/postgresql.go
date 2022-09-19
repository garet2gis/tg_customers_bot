package db

import (
	"context"
	"fmt"
	"github.com/garet2gis/tg_customers_bot/internal/service"
	"github.com/garet2gis/tg_customers_bot/pkg/client/postgresql"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	repeatable "github.com/garet2gis/tg_customers_bot/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) service.Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

func (r Repository) FindOne(ctx context.Context, id string) (*service.PaidService, error) {
	q := `
		SELECT 
		       service.service_id, 
		       service.name, 
		       service.base_duration
		FROM service
		WHERE service.service_id = $1
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	var s service.PaidService
	err := r.client.QueryRow(ctx, q, id).Scan(&s.ID, &s.Name, &s.BaseDuration)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r Repository) FindAll(ctx context.Context) ([]service.PaidService, error) {
	q := `
		SELECT 
		       service.service_id, 
		       service.name, 
		       service.base_duration
		FROM service 
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	services := make([]service.PaidService, 0)
	for rows.Next() {
		var s service.PaidService
		if err = rows.Scan(&s.ID, &s.Name, &s.BaseDuration); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (r Repository) Create(ctx context.Context, s *service.PaidService) error {
	q := `
		INSERT INTO service (name, base_duration) 
		VALUES ($1, $2)
		RETURNING service_id
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	if err := r.client.QueryRow(ctx, q, s.Name, s.BaseDuration).Scan(&s.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("Code: %s, Message: %s, Where: %s, Detail: %s, SQLState: %s", pgErr.Code, pgErr.Message, pgErr.Where, pgErr.Detail, pgErr.SQLState())
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r Repository) Update(ctx context.Context, s service.PaidService) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
