package seed_test

import (
	"database/sql"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/vkcku/pastyears/internal/seed"
	"github.com/vkcku/pastyears/internal/testutils"
)

func TestSeed(t *testing.T) {
	t.Parallel()

	var (
		tables []string

		ctx        = t.Context()
		tx         = testutils.TestTx(t)
		exceptions = []string{
			"questions.rich_text",
			"questions.prelims_questions",
			"questions.prelims_questions_topics",
		}
	)

	// Fetch the table names.
	{
		rows, err := tx.Query(
			ctx,
			"select schemaname || '.' || tablename from pg_catalog.pg_tables where schemaname not in ('pg_catalog', 'information_schema')", //nolint:lll
		)
		if err != nil {
			t.Fatal(err)
		}

		tables, err = pgx.CollectRows(
			rows,
			func(row pgx.CollectableRow) (string, error) {
				var table string

				err := row.Scan(&table)

				return table, err //nolint:wrapcheck
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := seed.Seed(ctx, tx); err != nil {
		t.Fatal(err)
	}

	for _, table := range tables {
		var (
			value int
			query = "SELECT 1 FROM " + table + " LIMIT 1"
			row   = tx.QueryRow(ctx, query)
		)

		err := row.Scan(&value)

		if slices.Contains(exceptions, table) {
			if errors.Is(err, pgx.ErrNoRows) == false {
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
		tx    = testutils.TestTx(t)
	)

	err := seed.Seed(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	row := tx.QueryRow(
		ctx,
		"SELECT 1 FROM questions.question_papers WHERE year = $1",
		time.Now().Year(),
	)
	if err := row.Scan(&value); err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatalf("wanted 1, got %d", value)
	}
}
