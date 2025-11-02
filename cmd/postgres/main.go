//nolint:wrapcheck
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/vkcku/pastyears/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "postgres: %+v\n", err)

		os.Exit(1)
	}
}

func getHost() string {
	if host := os.Getenv("PGHOST"); host != "" {
		return host
	}

	return filepath.Join(os.TempDir(), "pastyears-pg")
}

func getPort() uint {
	s := os.Getenv("PGPORT")
	if s == "" {
		return 5432
	}

	port, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		panic(fmt.Errorf("failed to parse port: %s", s)) //nolint:err113
	}

	return uint(port)
}

func getDatabase() string {
	if s := os.Getenv("PGDATABASE"); s != "" {
		return s
	}

	return "pastyears"
}

func run() error {
	var (
		ctx    = context.Background()
		remove = flag.Bool(
			"remove",
			false,
			"If enabled, remove the instance if it exists.",
		)
		host = flag.String(
			"host",
			getHost(),
			"The directory that contains the postgres instance.",
		)
		port = flag.Uint(
			"port",
			getPort(),
			"The port to start postgres on.",
		)
		database = flag.String(
			"database",
			getDatabase(),
			"The database to create.",
		)
	)

	flag.Parse()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if *remove {
		return postgres.Remove(ctx, *host)
	}

	status, err := postgres.GetStatus(ctx, *host)
	if err != nil {
		return err
	}

	switch status {
	case postgres.Running:
		return nil
	case postgres.NotRunning:
		return postgres.Start(ctx, *host)
	case postgres.NotFound:
		_, err := postgres.New(
			ctx,
			*database,
			*host,
			uint16(*port), //nolint:gosec
		)

		return err
	}

	panic("unreachable")
}
