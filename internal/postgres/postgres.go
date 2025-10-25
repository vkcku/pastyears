// Package postgres manages dev/test instances of Postgres. This requires that
// postgres binaries `pgctl`, `initdb` etc. to be in PATH.
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

var (
	errEmptyPgDir    = errors.New("postgres: pgDir was empty")
	errNoPort        = errors.New("postgres: port was 0")
	errEmptyDatabase = errors.New("postgres: database was empty")
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

	u, err := user.Current()
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
	defer file.Close()

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
