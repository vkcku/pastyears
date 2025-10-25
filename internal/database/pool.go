// Package database handles creating connections, running transactions etc.
package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a connection pool to the database.
type Pool = pgxpool.Pool

// ErrEmptyConnectionString is returned if the connection string is empty.
var ErrEmptyConnectionString = errors.New("database: empty connection string")

// Config holds values for configuring the connection pool.
//
// For more details about each value, refer
// [github.com/jackc/pgx/v5/pgxpool.Config].
type Config struct {
	HealthCheckPeriod     time.Duration
	MaxConnIdleTime       time.Duration
	MaxConnLifetime       time.Duration
	MaxConnLifetimeJitter time.Duration
	MaxConns              uint8
	MinConns              uint8
	MinIdleConns          uint8
}

// NewConfig returns a new [Config] object with sane defaults.
func NewConfig() Config {
	return Config{
		HealthCheckPeriod: time.Second * 10,
		MaxConnIdleTime:   time.Minute * 20,
		// One of the responses in the following thread recommended 1 hour.
		//
		// https://www.postgresql.org/message-id/CA%2Bmi_8bnvpxHZtb6EgHSHY-xn29W8VJMzjPU3fiCOv1bfjrNuA%40mail.gmail.com
		MaxConnLifetime:       time.Hour,
		MaxConnLifetimeJitter: time.Minute * 2,
		MaxConns:              10,
		MinConns:              1,
		MinIdleConns:          1,
	}
}

// NewPool returns a new database pool.
func NewPool(
	ctx context.Context,
	connString string,
	config Config,
) (*Pool, error) {
	if connString == "" {
		return nil, ErrEmptyConnectionString
	}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf(
			"database: failed to parse connection string: %w",
			err,
		)
	}

	poolConfig.HealthCheckPeriod = config.HealthCheckPeriod
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnLifetimeJitter = config.MaxConnLifetimeJitter
	poolConfig.MaxConns = int32(config.MaxConns)
	poolConfig.MinConns = int32(config.MinConns)
	poolConfig.MinIdleConns = int32(config.MinIdleConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("database: failed to get pool: %w", err)
	}

	return pool, nil
}
