// Run the tests after creating a test postgres instance.
//
//nolint:wrapcheck
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/vkcku/pastyears/internal/postgres"
)

func main() {
	err := run()
	if err == nil {
		return
	}

	var (
		rc      = 1
		exitErr *exec.ExitError
	)

	if errors.As(err, &exitErr) {
		rc = exitErr.ExitCode()
	} else {
		fmt.Fprintf(os.Stderr, "test: %s\n", err)
	}

	os.Exit(rc)
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tmpdir, err := os.MkdirTemp(os.TempDir(), "pastyears-pg-test")
	if err != nil {
		return fmt.Errorf("failed to create temp dir for postgres: %w", err)
	}

	defer func() {
		if err := postgres.Remove(ctx, tmpdir); err != nil {
			panic(
				fmt.Errorf("failed to remove test postgres instance: %w", err),
			)
		}
	}()

	connstring, err := postgres.New(ctx, "pastyears_test", tmpdir, 6543)
	if err != nil {
		return err
	}

	args := append([]string{"test"}, os.Args[1:]...)
	cmd := exec.CommandContext(ctx, "go", args...) //nolint:gosec
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "PASTYEARS_TEST_DB_URL="+connstring)

	return cmd.Run()
}
