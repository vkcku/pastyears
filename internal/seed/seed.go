// Package seed handles seeding the database for local development and/or
// testing.
package seed

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed seed.sql
var seedQuery string

// Seed seeds all the tables except the prelims questions table.
func Seed(ctx context.Context, tx pgx.Tx) error {
	// This splitting is not robust, but for this limited usecase where the
	// input is known beforehand, it should be fine.
	for query := range strings.SplitSeq(seedQuery, ";") {
		_, err := tx.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("seeding failed for query: %w\n%s", err, query)
		}
	}

	return nil
}
