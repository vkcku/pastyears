package database_test

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/vkcku/pastyears/internal/database"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		db, err := database.New(":memory:")

		t.Cleanup(func() {
			if db != nil {
				if err := db.Close(); err != nil {
					panic(err)
				}
			}
		})

		if err != nil {
			t.Errorf("got error: %+v", err)
		}

		err = db.PingContext(t.Context())
		if err != nil {
			t.Errorf("got error when pinging: %+v", err)
		}
	})

	t.Run("empty file path", func(t *testing.T) {
		t.Parallel()

		db, err := database.New("")

		if db != nil {
			t.Errorf("wanted nil db, got +%v", db)
		}

		if errors.Is(err, database.ErrEmptyFilePath) == false {
			t.Errorf("wanted +%v, got +%v", database.ErrEmptyFilePath, err)
		}
	})
}

func TestDefaults(t *testing.T) {
	t.Parallel()

	getDB := func(t *testing.T) *sql.DB {
		t.Helper()

		// Some of the defaults, such as `journal_mode` do not give the correct
		// result for in memory databases.
		file := filepath.Join(t.TempDir(), "data.db")

		db, err := database.New(file)
		if err != nil {
			t.Fatalf("failed to get database: %+v", err)
		}

		t.Cleanup(func() {
			if err := db.Close(); err != nil {
				panic(fmt.Errorf("failed to close database: %w", err))
			}
		})

		return db
	}

	type testCase struct {
		pragma   string
		expected string
	}

	cases := []testCase{
		{pragma: "foreign_keys", expected: "1"},
		{pragma: "journal_mode", expected: "wal"},
		// 1 => NORMAL
		{pragma: "synchronous", expected: "1"},
	}

	for _, tc := range cases {
		t.Run(tc.pragma, func(t *testing.T) {
			t.Parallel()

			var (
				actual string
				db     = getDB(t)
			)

			row := db.QueryRowContext(t.Context(), "PRAGMA "+tc.pragma)

			err := row.Scan(&actual)
			if err != nil {
				t.Fatalf("query error: %+v", err)
			}

			if actual != tc.expected {
				t.Fatalf("wanted '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}

func TestNew2(t *testing.T) {
	t.Parallel()

	t.Run("empty connection string", func(t *testing.T) {
		t.Parallel()

		pool, err := database.New2(t.Context(), "")

		if errors.Is(err, database.ErrEmptyConnectionString) == false {
			t.Errorf(
				"wanted '%+v', got '+%v'",
				database.ErrEmptyConnectionString,
				err,
			)
		}

		if pool != nil {
			t.Error("pool must be nil")
		}
	})
}
