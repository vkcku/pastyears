// Package migrations deals with running database migrations.
package migrations

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	// Register the driver.
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

var (
	//go:embed migrations/*.sql
	migrations embed.FS
)

func newDb(connstring string) (*dbmate.DB, error) {
	u, err := url.Parse(connstring)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the connection string: %w", err)
	}

	db := dbmate.New(u)

	db.AutoDumpSchema = false
	db.Log = io.Discard
	db.FS = migrations
	db.MigrationsDir = []string{"migrations"}

	return db, nil
}

// New creates a new migration file with the given name.
func New(connstring string, name string) error {
	db, err := newDb(connstring)
	if err != nil {
		return err
	}

	if err := db.NewMigration(name); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	return nil
}

// Drop the database if it exists.
func Drop(connstring string) error {
	db, err := newDb(connstring)
	if err != nil {
		return err
	}

	if err := db.Drop(); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}

// Run runs all the pending migrations. If schemaFile is an non-empty string,
// then the schema is dumped to that file.
func Run(connstring string, schemaFile string) error {
	db, err := newDb(connstring)
	if err != nil {
		return err
	}

	if err := db.CreateAndMigrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return dumpSchemaToFile(db, schemaFile)
}

// Down rollsback the latest migration.
func Down(connstring string, schemaFile string) error {
	db, err := newDb(connstring)
	if err != nil {
		return err
	}

	if err := db.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	return dumpSchemaToFile(db, schemaFile)
}

func dumpSchemaToFile(db *dbmate.DB, file string) error {
	if file == "" {
		return nil
	}

	schema, err := os.Create(file) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", file, err)
	}

	defer func() {
		if err := schema.Close(); err != nil {
			panic(fmt.Errorf("failed to close '%s': %w", file, err))
		}
	}()

	return dumpSchema(db, schema)
}

// dumpSchema writes the schema to the given writer.
func dumpSchema(db *dbmate.DB, writer io.Writer) error {
	// Unfortunately, dbmate does not provide an API to dump the schema into a
	// writer. So write the schema into a temp file and read from that.
	tmpdir, err := os.MkdirTemp(os.TempDir(), "pastyears-pg-schema")
	if err != nil {
		return fmt.Errorf("failed to make temp dir: %w", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			panic(fmt.Errorf("failed to remove tmp dir: %w", err))
		}
	}()

	schemaFile := filepath.Join(tmpdir, "schema.sql")

	db.SchemaFile = schemaFile
	if err := db.DumpSchema(); err != nil {
		return fmt.Errorf("failed to dump schema: %w", err)
	}

	schema, err := os.Open(schemaFile) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open schema file: %w", err)
	}

	defer func() {
		if err := schema.Close(); err != nil {
			panic(fmt.Errorf("failed to close schema file: %w", err))
		}
	}()

	scanner := bufio.NewScanner(schema)
	for scanner.Scan() {
		b := scanner.Bytes()

		// REFERENCE: https://github.com/amacneil/dbmate/issues/678
		if bytes.HasPrefix(b, []byte("\\restrict ")) ||
			bytes.HasPrefix(b, []byte("\\unrestrict ")) {
			continue
		}

		_, err := writer.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write processed schema: %w", err)
		}

		_, err = writer.Write([]byte{'\n'})
		if err != nil {
			return fmt.Errorf("failed to write processed schema: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read dumped schema file: %w", err)
	}

	return nil
}
