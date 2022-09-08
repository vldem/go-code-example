package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgres returns DB
func NewPostgres(ctx context.Context, psqlConn string) (*pgxpool.Pool, error) {
	// connect to database
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
