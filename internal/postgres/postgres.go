// Package postgres manages dev/test instances of Postgres. This requires that
// postgres binaries `pgctl`, `initdb` etc. to be in PATH.
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
)

// Status specifies the status of the postgres instance.
type Status uint8

const (
	// Running means the Postgres instance is running.
	Running Status = iota

	// NotRunning means a Postgres instance exists, but it is not running.
	NotRunning

	// NoInstance means there is no Postgres instance.
	NoInstance
)

var (
	errEmptyPgDir    = errors.New("postgres: pgDir was empty")
	errNoPort        = errors.New("postgres: port was 0")
	errEmptyDatabase = errors.New("postgres: database was empty")
	errNoInstance    = errors.New("postgres: no instance found")
)

// New will create a new Postgres instance in the given directory that listens
// on the given port. The superuser will be the current user as per
// `user.Current()`. This will also create the given database. This will fail if
// the database is already running or the given directory exists.
func New(
	ctx context.Context,
	pgDir string,
	port uint16,
	database string,
) error {
	if pgDir == "" {
		return errEmptyPgDir
	}

	if port == 0 {
		return errNoPort
	}

	if database == "" {
		return errEmptyDatabase
	}

	var (
		dataDir   = getDataDir(pgDir)
		logFile   = filepath.Join(pgDir, "logs.log")
		socketDir = pgDir
	)

	u, err := user.Current() //nolint:varnamelen
	if err != nil {
		return fmt.Errorf("postgres: failed to get current user: %w", err)
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
		return wrapExitError("initdb", err)
	}

	file, err := os.OpenFile(
		filepath.Join(dataDir, "postgresql.conf"),
		os.O_RDWR|os.O_APPEND,
		0777,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to open postgresql.conf: %w", err)
	}
	defer file.Close() //nolint:errcheck

	_, err = fmt.Fprintf(file, "port = %d\n", port)
	if err != nil {
		return fmt.Errorf(
			"postgres: failed to write to postgresql.conf: %w",
			err,
		)
	}

	cmd = exec.CommandContext(
		ctx,
		"pg_ctl",
		"--pgdata",
		dataDir,
		"--options",
		fmt.Sprintf("'--unix_socket_directories=%s'", socketDir),
		"--log",
		logFile,
		"start",
	)
	if _, err = cmd.Output(); err != nil {
		return wrapExitError("pg_ctl", err)
	}

	cmd = exec.CommandContext(
		ctx,
		"createdb",
		"--host",
		socketDir,
		"--owner",
		u.Username,
		database,
	)
	if _, err = cmd.Output(); err != nil {
		return wrapExitError("createdb", err)
	}

	return nil
}

// Stop the postgres instance if it is running.
func Stop(ctx context.Context, pgDir string) error {
	status, err := GetStatus(ctx, pgDir)
	if err != nil {
		return err
	}

	if status != Running {
		return nil
	}

	cmd := exec.CommandContext(
		ctx,
		"pg_ctl",
		"--pgdata",
		getDataDir(pgDir),
		"stop",
	)
	if _, err = cmd.Output(); err != nil {
		return wrapExitError("pg_ctl", err)
	}

	return nil
}

// Remove stops and removes the postgres instance.
func Remove(ctx context.Context, pgDir string) error {
	err := Stop(ctx, pgDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(pgDir)
	if err != nil {
		return fmt.Errorf("postgres: failed to remove pgDir: %w", err)
	}

	return nil
}

// Start starts the postgres instance if it is not running. This will not create
// the instance if one does not exist.
func Start(ctx context.Context, pgDir string) error {
	status, err := GetStatus(ctx, pgDir)
	if err != nil {
		return err
	}

	switch status {
	case Running:
		return nil
	case NoInstance:
		return errNoInstance
	case NotRunning:
		// no-op
	}

	var (
		dataDir   = getDataDir(pgDir)
		socketDir = pgDir
		logFile   = filepath.Join(pgDir, "logs.log")
	)

	cmd := exec.CommandContext(
		ctx,
		"pg_ctl",
		"--pgdata",
		dataDir,
		"--options",
		fmt.Sprintf("'--unix_socket_directories=%s'", socketDir),
		"--log",
		logFile,
		"start",
	)
	if _, err = cmd.Output(); err != nil {
		return wrapExitError("pg_ctl", err)
	}

	return nil
}

// GetStatus returns the status of the postgres instance.
func GetStatus(ctx context.Context, pgDir string) (Status, error) {
	var (
		dataDir = getDataDir(pgDir)
		cmd     = exec.CommandContext(
			ctx,
			"pg_ctl",
			"--pgdata",
			dataDir,
			"status",
		)
	)

	_, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			switch exitErr.ExitCode() {
			case 3:
				return NotRunning, nil
			case 4:
				return NoInstance, nil
			}
		}

		return NotRunning, wrapExitError("pg_ctl", err)
	}

	return Running, nil
}

func wrapExitError(command string, err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf(
			"postgres: failed to run '%s': %w\n%s",
			command,
			err,
			exitErr.Stderr,
		)
	}

	return fmt.Errorf("postgres: failed to run '%s': %w", command, err)
}

func getDataDir(pgDir string) string {
	return filepath.Join(pgDir, "data")
}
