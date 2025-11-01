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
