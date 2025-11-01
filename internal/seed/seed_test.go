package seed_test

import (
	"database/sql"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/vkcku/pastyears/internal/seed"
	"github.com/vkcku/pastyears/internal/testutils"
)

func TestSeed(t *testing.T) {
	t.Parallel()

	var (
		ctx        = t.Context()
		db         = testutils.TestDB(t)
		tables     = make([]string, 0, 10)
		exceptions = []string{
			"rich_text",
			"prelims_questions",
			"prelims_questions_topics",
		}
	)

	// Fetch the table names.
	{
		rows, err := db.QueryContext(
			ctx,
			"select name from sqlite_schema where type = 'table'",
		)
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			if err := rows.Err(); err != nil {
				t.Fatal(err)
			}

			if err := rows.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		for rows.Next() {
			var table string

			err := rows.Scan(&table)
			if err != nil {
				t.Fatal(err)
			}

			tables = append(tables, table)
		}
	}

	if err := seed.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}

	for _, table := range tables {
		var (
			value int
			query = "SELECT 1 FROM " + table + " LIMIT 1"
			row   = db.QueryRowContext(ctx, query)
		)

		err := row.Scan(&value)

		if slices.Contains(exceptions, table) {
			if errors.Is(err, sql.ErrNoRows) == false {
				t.Errorf(
					"wanted '%s', got '%s' for table '%s'",
					sql.ErrNoRows,
					err,
					table,
				)
			}

			continue
		}

		if err != nil {
			t.Errorf("got error for table '%s': %+v", table, err)
		}

		if value != 1 {
			t.Errorf("wanted 1, got %+v for table '%s'", value, table)
		}
	}
}

// TestQuestionPaperLatestYear ensures that the seeding includes the current
// year for the question papers table. This is to ensure that the `CHECK`
// constraint is updated when the year changes.
func TestQuestionPaperLatestYear(t *testing.T) {
	t.Parallel()

	var (
		value int
		ctx   = t.Context()
		db    = testutils.TestDB(t)
	)

	if err := seed.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}

	row := db.QueryRowContext(
		ctx,
		"SELECT 1 FROM question_papers WHERE year = ?",
		time.Now().Year(),
	)
	if err := row.Scan(&value); err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("wanted 1, got %d", value)
	}
}
