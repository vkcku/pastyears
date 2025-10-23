// Package check contains test helpers to check for equality etc. These helpers
// will "soft" fail meaning the test will keep running. This is equivalent to
// `testify.require`.
package check

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var useColor = flag.Bool( //nolint:gochecknoglobals
	"pastyears.color",
	true,
	"If enabled, display colored diffs in test outputs.",
)

// Equal checks if the given values are equal. If they are not equal, then the
// diff is added to the error with `[testing.T.Error]` and this returns false.
// If they are equal then this returns true.
func Equal[T any](t *testing.T, want T, got T) bool {
	t.Helper()

	diff := cmp.Diff(want, got)
	equal := diff == ""

	if equal == false {
		t.Error(ansiDiff(diff))
	}

	return equal
}

func ansiDiff(diff string) string {
	// REFERENCE:
	// https://github.com/google/go-cmp/issues/230#issuecomment-665750648
	if diff == "" {
		return ""
	}

	if *useColor == false {
		return diff
	}

	lines := strings.Split(diff, "\n")
	for i, s := range lines {
		switch {
		case strings.HasPrefix(s, "-"):
			lines[i] = escapeCode(31) + s + escapeCode(0)
		case strings.HasPrefix(s, "+"):
			lines[i] = escapeCode(32) + s + escapeCode(0)
		}
	}

	return strings.Join(lines, "\n")
}

func escapeCode(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}
