package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"glovee-worker/internal/config"
)

type DB struct {
	*sql.DB
}

func NewDB(ctx context.Context, config *config.Config) (*DB, error) {
	const operation = "repository.postgres.NewDB"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Database,
		config.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
