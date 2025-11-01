// Package seed handles seeding the database for local development and/or
// testing.
package seed

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/vkcku/pastyears/internal/database"
)

//go:embed seed.sql
var seedQuery string

// Seed seeds all the tables except the prelims questions table.
func Seed(ctx context.Context, db *database.Database) error {
	_, err := db.ExecContext(ctx, seedQuery)
	if err != nil {
		return fmt.Errorf("seed: failed to run seed query: %w", err)
	}

	return nil
}
