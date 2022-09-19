package postgresql

import (
	"context"
	"fmt"
	"github.com/garet2gis/tg_customers_bot/internal/config"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, log logging.Logger, dc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dc.Username, dc.Password, dc.Host, dc.Port, dc.DatabaseName)
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, dc.MaxAttempts, 2*time.Second)

	if err != nil {
		log.Fatal("failed to connect to postgresql")
		return nil, err
	}

	return pool, nil
}
