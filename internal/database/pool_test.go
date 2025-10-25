package database_test

import (
	"errors"
	"testing"

	"github.com/vkcku/pastyears/internal/database"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	t.Run("empty connection string", func(t *testing.T) {
		t.Parallel()

		pool, err := database.NewPool(t.Context(), "")

		if errors.Is(err, database.ErrEmptyConnectionString) == false {
			t.Errorf(
				"wanted '%+v', got '%+v'",
				database.ErrEmptyConnectionString,
				err,
			)
		}

		if pool != nil {
			t.Errorf("got non nil pool")
		}
	})

	t.Run("valid connection string", func(t *testing.T) {
		t.Parallel()

		pool, err := database.NewPool(
			t.Context(),
			"postgresql://username:password@127.0.0.1:5432/foobar",
		)
		if err != nil {
			t.Errorf("got error: %+v", err)
		}

		if pool == nil {
			t.Error("got nil pool")
		}
	})
}
