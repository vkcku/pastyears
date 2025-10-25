package database_test

import (
	"errors"
	"testing"

	"github.com/vkcku/pastyears/internal/database"
	"github.com/vkcku/pastyears/internal/testutils/check"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	t.Run("empty connection string", func(t *testing.T) {
		t.Parallel()

		pool, err := database.NewPool(t.Context(), "", database.NewConfig())

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
			database.NewConfig(),
		)
		if err != nil {
			t.Errorf("got error: %+v", err)
		}

		if pool == nil {
			t.Error("got nil pool")
		}
	})

	t.Run("configurations are set", func(t *testing.T) {
		t.Parallel()

		config := database.NewConfig()

		pool, err := database.NewPool(
			t.Context(),
			"postgresql://username:password@127.0.0.1:5432/foobar",
			database.NewConfig(),
		)
		if err != nil {
			t.Errorf("got error: %+v", err)
		}

		if pool == nil {
			t.Fatal("got nil pool")
		}

		poolConfig := pool.Config()
		actual := database.Config{
			HealthCheckPeriod:     poolConfig.HealthCheckPeriod,
			MaxConnIdleTime:       poolConfig.MaxConnIdleTime,
			MaxConnLifetime:       poolConfig.MaxConnLifetime,
			MaxConnLifetimeJitter: poolConfig.MaxConnLifetimeJitter,
			MaxConns: uint8( //nolint:gosec
				poolConfig.MaxConns,
			),
			MinConns: uint8( //nolint:gosec
				poolConfig.MinConns,
			),
			MinIdleConns: uint8( //nolint:gosec
				poolConfig.MinIdleConns,
			),
		}

		check.Equal(t, config, actual)
	})
}
