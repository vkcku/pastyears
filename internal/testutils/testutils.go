// Package testutils contains helpers and utility functions for running tests.
package testutils

import (
	"testing"
)

func init() { //nolint:gochecknoinits
	if testing.Testing() == false {
		panic("testutils: cannot import outside of tests")
	}
}

// Main is a common main function that should be called from `TestMain` from
// all the tests.
func Main(m *testing.M) int {
	defer func() {
		closePool()
	}()

	return m.Run()
}
