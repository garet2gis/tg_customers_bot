package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/garet2gis/tg_customers_bot/internal/paid_service"
	"github.com/garet2gis/tg_customers_bot/pkg/client/postgresql"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	repeatable "github.com/garet2gis/tg_customers_bot/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) paid_service.Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

func (r Repository) FindOne(ctx context.Context, id string) (*paid_service.PaidService, error) {
	q := `
		SELECT 
		       service.service_id, 
		       service.name, 
		       service.base_duration
		FROM service
		WHERE service.service_id = $1
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	var s paid_service.PaidService
	err := r.client.QueryRow(ctx, q, id).Scan(&s.ID, &s.Name, &s.BaseDuration)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r Repository) FindOneByIndex(ctx context.Context, id int) (*paid_service.PaidService, error) {
	q := `
		SELECT 
		       service_id, 
		       name, 
		       base_duration
		FROM service
		WHERE service_id IN 
		      (SELECT service_id FROM service
		      ORDER BY created_at
		      LIMIT 1 OFFSET $1)
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	var s paid_service.PaidService
	err := r.client.QueryRow(ctx, q, id-1).Scan(&s.ID, &s.Name, &s.BaseDuration)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r Repository) FindAll(ctx context.Context) ([]paid_service.PaidService, error) {
	q := `
		SELECT 
		       service.service_id, 
		       service.name, 
		       service.base_duration
		FROM service
		ORDER BY service.created_at
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	services := make([]paid_service.PaidService, 0)
	for rows.Next() {
		var s paid_service.PaidService
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

func (r Repository) Create(ctx context.Context, s *paid_service.PaidService) error {
	q := `
		INSERT INTO service (name, base_duration) 
		VALUES ($1, $2)
		RETURNING service_id
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))
	if err := r.client.QueryRow(ctx, q, s.Name, s.BaseDuration).Scan(&s.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("Code: %s, Message: %s, Where: %s, Detail: %s, SQLState: %s", pgErr.Code, pgErr.Message, pgErr.Where, pgErr.Detail, pgErr.SQLState())
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r Repository) Update(ctx context.Context, s paid_service.PaidService) error {
	q := `
		UPDATE service 
		SET name = $1,
			base_duration = $2
		WHERE service_id = $3
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))

	commandTag, err := r.client.Exec(ctx, q, s.Name, s.BaseDuration, s.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("Code: %s, Message: %s, Where: %s, Detail: %s, SQLState: %s", pgErr.Code, pgErr.Message, pgErr.Where, pgErr.Detail, pgErr.SQLState())
			if pgErr.Code == "23503" {
				return paid_service.CanNotDeleteRowForeignKey
			}

			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return paid_service.NoRowsDeleted
	}

	return nil

}

func (r Repository) Delete(ctx context.Context, id int) error {
	q := `
		DELETE 
		FROM  service 
		WHERE service_id IN 
		      (SELECT service_id FROM service
		      ORDER BY created_at
		      LIMIT 1 OFFSET $1)
		`
	r.logger.Trace("SQL query: %s", repeatable.FormatQuery(q))

	commandTag, err := r.client.Exec(ctx, q, id-1)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("Code: %s, Message: %s, Where: %s, Detail: %s, SQLState: %s", pgErr.Code, pgErr.Message, pgErr.Where, pgErr.Detail, pgErr.SQLState())
			if pgErr.Code == "23503" {
				return paid_service.CanNotDeleteRowForeignKey
			}

			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return paid_service.NoRowsDeleted
	}

	return nil
}
