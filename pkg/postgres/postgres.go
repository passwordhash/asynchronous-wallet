package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*pgxpool.Config)

// WithMaxConns sets the maximum number of connections in the pool.
func WithMaxConns(maxConns int32) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MaxConns = maxConns
	}
}

// NewPool creates a new PostgreSQL connection pool with the provided DSN and options.
func NewPool(ctx context.Context, dsn string, opts ...Option) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(config)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
