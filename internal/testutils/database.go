//nolint:gochecknoglobals
package testutils

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/vkcku/pastyears/internal/database"
)

var (
	migrations     string
	migrationsLock sync.RWMutex
)

// TestDB returns a clean in-memory database for the given test. This will be
// automatically closed at the end of the tests.
func TestDB(t *testing.T) *database.Database {
	t.Helper()

	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("testutils: failed to get database: %+v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	})

	_, err = db.ExecContext(t.Context(), getSchema(t))
	if err != nil {
		t.Fatalf("testutils: failed to run migrations: %+v", err)
	}

	return db
}

// getSchema returns the schema of the database.
//
// The API `dbmate` provides cannot be used because there is no way to inject an
// instance of `*sql.DB` in a concurrent safe manner into the sqlite driver.
// This means we cannot use an in-memory database.
func getSchema(t *testing.T) string {
	t.Helper()

	migrationsLock.RLock()

	if migrations != "" {
		m := migrations

		migrationsLock.RUnlock()

		return m
	}

	migrationsLock.RUnlock()

	migrationsLock.Lock()
	defer migrationsLock.Unlock()

	_, file, _, ok := runtime.Caller(0)
	if ok == false {
		t.Fatalf("testutils: failed to get root directory via runtime.Caller")
	}

	// ../../db/schema.sql
	schemaFile := filepath.Join(
		filepath.Dir(file),
		"..",
		"..",
		"db",
		"schema.sql",
	)

	migrationsRaw, err := os.ReadFile(schemaFile) //nolint:gosec
	if err != nil {
		t.Fatalf(
			"testutils: failed to read migrations file '%s': %+v",
			schemaFile,
			err,
		)
	}

	migrations = string(migrationsRaw)

	return migrations
}
