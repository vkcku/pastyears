// A CLI for managing Pastyears because shell scripts suck.
//
//nolint:exhaustruct,wrapcheck
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/vkcku/pastyears/internal/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := &cli.Command{
		Usage: "A CLI for managing Pastyears because shell scripts suck.",
		Commands: []*cli.Command{
			lintCommand(),
			preCommitCommand(),
			postgresCommand(),
		},
	}

	err := cmd.Run(ctx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

type errMissingValue struct {
	field string
}

func (e errMissingValue) Error() string {
	return e.field + " not set"
}

func postgresCommand() *cli.Command {
	var (
		createInstance bool
		pgDir          string
		port           uint16 = 5432
		database              = "pastyears"

		pgDirFlag = &cli.StringFlag{
			Name:        "pgdir",
			Usage:       "The directory to store all the postgres related files.",
			Sources:     cli.EnvVars("PGHOST"),
			Destination: &pgDir,
			Validator: func(s string) error {
				if s == "" {
					return errMissingValue{"pgdir"}
				}

				return nil
			},
		}
	)

	return &cli.Command{
		Name:    "postgres",
		Usage:   "Manage the dev postgres instance.",
		Aliases: []string{"pg"},
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "Create a new Postgres instance.",
				Flags: []cli.Flag{
					pgDirFlag,
					&cli.StringFlag{
						Name:        "database",
						Usage:       "The database to create.",
						Value:       "pastyears",
						Sources:     cli.EnvVars("PGDATABASE"),
						Destination: &database,
						Validator: func(s string) error {
							if s == "" {
								return errMissingValue{"database"}
							}

							return nil
						},
					},
					&cli.Uint16Flag{
						Name:        "port",
						Usage:       "The port to start postgres on listen on.",
						Value:       5432,
						Sources:     cli.EnvVars("PGPORT"),
						Destination: &port,
						Validator: func(u uint16) error {
							if u == 0 {
								return errMissingValue{"port"}
							}

							return nil
						},
					},
				},
				Action: func(ctx context.Context, _ *cli.Command) error {
					return postgres.New(ctx, pgDir, port, database)
				},
			},
			{
				Name:  "start",
				Usage: "Start the postgres instance if it is not running.",
				Flags: []cli.Flag{pgDirFlag, &cli.BoolFlag{
					Name:        "create",
					Usage:       "Create the instance if one does not already exist. This will use the defaults from `postgres new`.",
					Value:       true,
					Destination: &createInstance,
				}},
				Action: func(ctx context.Context, _ *cli.Command) error {
					if createInstance {
						status, err := postgres.Status(ctx, pgDir)
						if err != nil {
							return err
						}

						if status == postgres.NoInstance {
							return postgres.New(ctx, pgDir, port, database)
						}
					}

					return postgres.Start(ctx, pgDir)
				},
			},
			{
				Name:  "stop",
				Usage: "Stop the postgres instance if it is running.",
				Flags: []cli.Flag{pgDirFlag},
				Action: func(ctx context.Context, _ *cli.Command) error {
					return postgres.Stop(ctx, pgDir)
				},
			},
		},
	}
}

func lintCommand() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Run all formatters and linters while autofixing issues where possible.",
		Commands: []*cli.Command{
			{
				Name:  "go",
				Usage: "Lint all go files.",
				Action: func(ctx context.Context, _ *cli.Command) error {
					args := []string{"run"}
					if isInCI() == false {
						// Running `--fix` in CI is not great, because the error
						// comes up as permission denied since the files are
						// moved
						// to the nix store before running the linting.
						args = append(args, "--fix", "--fast-only")
					}

					args = append(args, "./...")

					cmd := newCommand(ctx, "golangci-lint", args...)

					return cmd.Run()
				},
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			cmd := newCommand(ctx, "treefmt")

			return cmd.Run()
		},
	}
}

func preCommitCommand() *cli.Command {
	return &cli.Command{
		Name:  "pre-commit",
		Usage: "Run all pre-commit hooks.",
		Action: func(ctx context.Context, _ *cli.Command) (err error) {
			stashed := false

			defer func() {
				if stashed == false {
					return
				}

				cmd := newCommand(ctx, "git", "stash", "pop", "--quiet")
				err = errors.Join(err, cmd.Run())
			}()

			cmd := newCommand(
				ctx,
				"git",
				"stash",
				"--quiet",
				"--keep-index",
				"--include-untracked",
			)
			if err = cmd.Run(); err != nil {
				return err
			}

			stashed = true

			changedFilesRaw := strings.Builder{}
			cmd = newCommand(
				ctx,
				"git",
				"diff",
				"--diff-filter",
				"d",
				"--name-only",
				"--cached",
			)

			cmd.Stdout = &changedFilesRaw
			if err = cmd.Run(); err != nil {
				return err
			}

			args := []string{"--ci"}
			args = append(
				args,
				strings.Split(changedFilesRaw.String(), "\n")...)

			cmd = newCommand(ctx, "treefmt", args...)
			if err = cmd.Run(); err != nil {
				return err
			}

			return err
		},
	}
}

func isInCI() bool {
	_, ok := os.LookupEnv("CI")

	return ok
}

// newCommand returns a command with stdout and stderr set to `os.Stdout` and
// `os.Stderr`.
func newCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
