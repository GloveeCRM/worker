package postgres

import (
	"context"
	"fmt"
	"glovee-worker/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, config *config.Config) (*DB, error) {
	const operation = "repository.postgres.NewDB"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
		config.Postgres.SSLMode,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
