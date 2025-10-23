// Package check contains test helpers to check for equality etc. These helpers
// will "soft" fail meaning the test will keep running. This is equivalent to
// `testify.require`.
package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal checks if the given values are equal. If they are not equal, then the
// diff is added to the error with `[testing.T.Error]` and this returns false.
// If they are equal then this returns true.
func Equal[T any](t *testing.T, want T, got T) bool {
	t.Helper()

	diff := cmp.Diff(want, got)
	equal := diff == ""

	if equal == false {
		t.Error(diff)
	}

	return equal
}
