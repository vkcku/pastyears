// A CLI for managing Pastyears because shell scripts suck.
//
//nolint:exhaustruct
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := &cli.Command{
		Usage:    "A CLI for managing Pastyears because shell scripts suck.",
		Commands: []*cli.Command{lintCommand()},
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

// newCommand returns a command with stdout and stderr set to `os.Stdout` and
// `os.Stderr`.
func newCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
