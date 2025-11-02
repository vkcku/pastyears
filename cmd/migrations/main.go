//nolint:wrapcheck
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	migrations "github.com/vkcku/pastyears/db"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "migrations: %+v\n", err)

		os.Exit(1)
	}
}

var errEmptyConnString = errors.New("no connection string")

func run() error {
	var (
		connstring = flag.String(
			"connstring",
			os.Getenv("PASTYEARS_DB_URL"),
			"The connection string. This must be provided.",
		)
		run  = flag.Bool("run", false, "Run all pending migrations.")
		drop = flag.Bool(
			"drop",
			false,
			"If enabled, drop and recreate the database and run all migrations.",
		)
		down = flag.Bool(
			"down",
			false,
			"If enabled, rollback the latest migration.",
		)
		newMigration = flag.String(
			"new",
			"",
			"Create a migration with the given file.",
		)
		schemaFile = flag.String(
			"schema",
			"./db/schema.sql",
			"The dumped schema file.",
		)
	)

	flag.Parse()

	if *connstring == "" {
		return errEmptyConnString
	}

	switch {
	case *newMigration != "":
		return migrations.New(*connstring, *newMigration)
	case *down:
		return migrations.Down(*connstring, *schemaFile)
	case *run:
		if *drop {
			if err := migrations.Drop(*connstring); err != nil {
				return err
			}
		}

		return migrations.Run(*connstring, *schemaFile)
	}

	return nil
}
