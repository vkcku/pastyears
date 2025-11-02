package testutils

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"

	migrations "github.com/vkcku/pastyears/db"
	"github.com/vkcku/pastyears/internal/database"
)

var (
	connstring = flag.String(
		"pastyears.database",
		os.Getenv("PASTYEARS_TEST_DB_URL"),
		"The connection string for the test database instance. If empty, those tests are skipped. The default value is taken from the $PASTYEARS_TEST_DB_URL.", //nolint:lll
	)

	getPool = sync.OnceValue(func() *database.Pool {
		ctx := context.Background()
		pool, err := database.New2(ctx, *connstring)
		if err != nil {
			panic(fmt.Errorf("testutils: %w", err))
		}

		if err := pool.Ping(ctx); err != nil {
			panic(fmt.Errorf("testutils: ping failed: %w", err))
		}

		// Multiple test packages may try to run the migrations simultaneously.
		_, err = pool.Exec(ctx, "SELECT pg_advisory_lock(1)")
		if err != nil {
			panic(
				fmt.Errorf(
					"testutils: failed to get lock for running migrations: %w",
					err,
				),
			)
		}
		defer func() {
			_, err := pool.Exec(ctx, "SELECT pg_advisory_unlock(1)")
			if err != nil {
				panic(
					fmt.Errorf(
						"testutils: failed to release lock for running migrations: %w",
						err,
					),
				)
			}
		}()

		err = migrations.Run(*connstring, "")
		if err != nil {
			panic(fmt.Errorf("testutils: %w", err))
		}

		return pool
	})
)

// TestTx returns a transaction that is automatically rolled back at the end of
// the test. It is an error for this to be committed or rolled back by the test.
func TestTx(t *testing.T) pgx.Tx { //nolint:ireturn
	t.Helper()

	if *connstring == "" {
		t.Skip("connection string not set via -database")
	}

	tx, err := getPool().Begin(t.Context())
	if err != nil {
		t.Fatalf("testutils: failed to start transaction: %+v", err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil {
			panic(
				fmt.Errorf(
					"testutils: failed to rollback transaction: %w",
					err,
				),
			)
		}
	})

	return tx
}
