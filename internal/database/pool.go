// Package database handles creating connections, running transactions etc.
package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a connection pool to the database.
type Pool = pgxpool.Pool

// ErrEmptyConnectionString is returned if the connection string is empty.
var ErrEmptyConnectionString = errors.New("database: empty connection string")

// NewPool returns a new database pool.
func NewPool(ctx context.Context, connString string) (*Pool, error) {
	if connString == "" {
		return nil, ErrEmptyConnectionString
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf(
			"database: failed to parse connection string: %w",
			err,
		)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("database: failed to get pool: %w", err)
	}

	return pool, nil
}
