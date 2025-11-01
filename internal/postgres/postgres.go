// Package postgres manages the lifecycle of a Postgres instance for
// development/testing.
//
//nolint:gosec
package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
)

// Status indicates the status of the Postgres instance.
type Status uint

const (
	// Running indicates that the Postgres instance is running.
	Running Status = iota
	// NotRunning indicates that the Postgres instance exists, but is not
	// running.
	NotRunning
	// NotFound indicates that a Postgres instance was not found.
	NotFound
)

// New creates a new postgres instance in the given directory listening to the
// given port. This returns a connection string to the instance.
func New(
	ctx context.Context,
	database string,
	dir string,
	port uint16,
) (string, error) {
	var (
		dataDir = filepath.Join(dir, "data")
		logFile = filepath.Join(dir, "logs.log")
	)

	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		"initdb",
		"--username",
		u.Username,
		"--auth=trust",
		"--no-instructions",
		dataDir,
	)
	if _, err = cmd.Output(); err != nil {
		return "", wrapError("initdb", err)
	}

	file, err := os.OpenFile(
		filepath.Join(dataDir, "postgresql.conf"),
		os.O_APPEND|os.O_WRONLY,
		0640,
	)
	if err != nil {
		return "", fmt.Errorf("failed to open postgresql.conf: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("failed to close postgresql.conf: %w", err))
		}
	}()

	_, err = fmt.Fprintf(file, "port = %d\n", port)
	if err != nil {
		return "", fmt.Errorf("failed to set port: %w", err)
	}

	if err := start(ctx, dir, dataDir, logFile); err != nil {
		return "", err
	}

	cmd = exec.CommandContext(
		ctx,
		"createdb",
		"--owner",
		u.Username,
		"--host",
		dir,
		"--port",
		strconv.FormatUint(uint64(port), 10),
		database,
	)
	if _, err := cmd.Output(); err != nil {
		return "", wrapError("createdb", err)
	}

	return "", nil
}

// Start starts the Postgres instance. The instance must have been initialized.
func Start(ctx context.Context, dir string) error {
	return start(
		ctx,
		dir,
		filepath.Join(dir, "data"),
		filepath.Join(dir, "logs.log"),
	)
}

func start(
	ctx context.Context,
	dir string,
	dataDir string,
	logfile string,
) error {
	cmd := exec.CommandContext(
		ctx,
		"pg_ctl",
		"--pgdata",
		dataDir,
		"--log",
		logfile,
		fmt.Sprintf("--options='--unix_socket_directories=%s'", dir),
		"start",
	)
	if _, err := cmd.Output(); err != nil {
		return wrapError("pg_ctl start", err)
	}

	return nil
}

// GetStatus returns the status of the running Postgres instance.
func GetStatus(ctx context.Context, dir string) (Status, error) {
	dataDir := filepath.Join(dir, "data")

	cmd := exec.CommandContext(ctx, "pg_ctl", "--pgdata", dataDir, "status")

	_, err := cmd.Output()
	if err == nil {
		return Running, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) == false {
		return NotFound, fmt.Errorf("pg_ctl failed: %w", err)
	}

	switch exitErr.ExitCode() {
	case 3:
		return NotRunning, nil
	case 4:
		return NotFound, nil
	}

	return NotFound, fmt.Errorf( //nolint:err113
		"pg_ctl error: %s",
		exitErr.Stderr,
	)
}

// Remove stops and removes the instance at the given directory if it exists.
func Remove(ctx context.Context, dir string) error {
	status, err := GetStatus(ctx, dir)
	if err != nil {
		return err
	}

	dataDir := filepath.Join(dir, "data")

	if status == Running {
		cmd := exec.CommandContext(ctx, "pg_ctl", "--pgdata", dataDir, "stop")
		if _, err := cmd.Output(); err != nil {
			return wrapError("pg_ctl stop", err)
		}
	}

	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to remove '%s': %w", dir, err)
	}

	return nil
}

func wrapError(command string, err error) error {
	if err == nil {
		panic("expected non-nil error")
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf( //nolint:err113
			"'%s' error: %s",
			command,
			exitErr.Stderr,
		)
	}

	return fmt.Errorf("'%s' failed: %w", command, err)
}
