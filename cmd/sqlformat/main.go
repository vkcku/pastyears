// Format SQL code.
//
// This exists because `sql-formatter` does not accept multiple files as
// input...
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprint(
			os.Stderr,
			"sqlformat: failed to format files\n",
			err.Error(),
		)

		os.Exit(1)
	}
}

func run(files []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var (
		wg   sync.WaitGroup
		errs []error
		mu   sync.Mutex
	)

	for _, file := range files {
		wg.Go(func() {
			cmd := exec.CommandContext( //nolint:gosec
				ctx,
				"sql-formatter",
				"-l",
				"sqlite",
				"--fix",
				file,
			)

			_, err := cmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					err = fmt.Errorf(
						"formatting '%s' failed: %s: %w",
						file,
						exitErr.Stderr,
						err,
					)
				}

				mu.Lock()

				errs = append(errs, err)

				mu.Unlock()
			}
		})
	}

	wg.Wait()

	return errors.Join(errs...)
}
