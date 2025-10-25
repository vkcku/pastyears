// A CLI for managing Pastyears because shell scripts suck.
//
//nolint:exhaustruct
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := &cli.Command{
		Usage: "A CLI for managing Pastyears because shell scripts suck.",
	}

	err := cmd.Run(ctx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
