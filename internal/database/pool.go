// Package database has wrappers/helpers to deal with the database.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// Autoimport the sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

// Database is a handle to interact with the database. It is safe for concurrent
// usage.
type Database = sql.DB

// ErrEmptyFilePath is returned when no path to a file is provided when creating
// a new connection.
var ErrEmptyFilePath = errors.New("database: no filepath given")

// New returns a new database connection pool with sensible defaults.
func New(file string) (*Database, error) {
	// TODO: Wrap the sqlite3 driver to add hooks so that queries can be tracked
	// and other lifecycle hooks can be added such as when opening/closing
	// connections. For example, this would be useful to run `PRAGMA optimize`
	// at appropriate times.
	if file == "" {
		return nil, ErrEmptyFilePath
	}

	options := url.Values{
		"_auto_vacuum": []string{"incremental"},
		"_busy_timeout": []string{
			strconv.FormatInt((time.Second * 3).Milliseconds(), 10),
		},
		"_foreign_keys": []string{"true"},
		"_journal_mode": []string{"wal"},
		"_synchronous":  []string{"NORMAL"},
		"_loc":          []string{"UTC"},
		"_txlock":       []string{"immediate"},
	}

	db, err := sql.Open("sqlite3", file+"?"+options.Encode())
	if err != nil {
		return nil, fmt.Errorf("database: failed to open database: %w", err)
	}

	return db, nil
}

// Pool is a connection pool that is safe for concurrent use.
type Pool = pgxpool.Pool

// ErrEmptyConnectionString is returned if the connection string is empty.
var ErrEmptyConnectionString = errors.New("database: empty connection string")

// New2 returns a new database pool.
func New2(ctx context.Context, connstring string) (*Pool, error) {
	// If the connection string is empty `pgx` will use the values from the
	// postgres
	// environment variables like `PGPORT` etc. which I am not a fan of.
	if connstring == "" {
		return nil, ErrEmptyConnectionString
	}

	config, err := pgxpool.ParseConfig(connstring)
	if err != nil {
		return nil, fmt.Errorf(
			"database: failed to parse connection string: %w",
			err,
		)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("database: failed to create pool: %w", err)
	}

	return pool, nil
}
