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
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := &cli.Command{
		Usage:    "A CLI for managing Pastyears because shell scripts suck.",
		Commands: []*cli.Command{lintCommand(), preCommitCommand()},
	}

	err := cmd.Run(ctx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
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
					cmd := newCommand(
						ctx,
						"golangci-lint",
						"run",
						"--fix",
						"./...",
					)

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

// newCommand returns a command with stdout and stderr set to `os.Stdout` and
// `os.Stderr`.
func newCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
