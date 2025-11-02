// Format SQL code.
//
// This exists because `sql-formatter` does not accept multiple files as
// input...
package main

import (
	"bytes"
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

	appendErr := func(err error) {
		mu.Lock()

		errs = append(errs, err)

		mu.Unlock()
	}

	for _, file := range files {
		wg.Go(func() {
			// `sql-formatter` will always write to the file even if there are
			// no formatting changes required which means `treefmt` will
			// interpret that as a failure in CI (also it just won't work
			// because nix makes the files read-only during the checks). So, do
			// a little bit of manual checking before writing the formatted file
			// back to disk.
			cmd := exec.CommandContext( //nolint:gosec
				ctx,
				"sql-formatter",
				"-l",
				"postgresql",
				file,
			)

			stdout, err := cmd.Output()
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

				appendErr(err)

				return
			}

			actual, err := os.ReadFile(file) //nolint:gosec
			if err != nil {
				appendErr(fmt.Errorf("failed to read '%s': %w", file, err))

				return
			}

			if bytes.Equal(actual, stdout) {
				return
			}

			err = os.WriteFile(file, stdout, 0640) //nolint:gosec
			if err != nil {
				appendErr(fmt.Errorf("failed to write '%s': %w", file, err))

				return
			}
		})
	}

	wg.Wait()

	return errors.Join(errs...)
}
