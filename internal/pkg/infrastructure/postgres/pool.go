package postgres

import (
	"context"
	"laboratory-internet-ai-test/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, config config.Postgres) (*pgxpool.Pool, error) {

	cfg, err := pgxpool.ParseConfig(config.Dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
