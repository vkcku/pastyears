// Package tassert is the same as the check package except if any of the
// functions fail, then the test fails immediately. This is equivalent to
// testify.require.
package tassert

import (
	"testing"

	"github.com/vkcku/pastyears/internal/testutils/check"
)

// Equal check sif the given values are equal. This fails the test immediately
// if they are not equal.
func Equal[T any](t *testing.T, want T, got T) {
	t.Helper()

	if check.Equal(t, want, got) == false {
		t.FailNow()
	}
}
