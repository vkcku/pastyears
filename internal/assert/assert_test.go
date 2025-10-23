package assert_test

import (
	"errors"
	"testing"

	"github.com/vkcku/pastyears/internal/assert"
	"github.com/vkcku/pastyears/internal/testutils/check"
)

func TestAssertf(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		assert.Assertf(true, "")
	})

	t.Run("panics", func(t *testing.T) {
		t.Parallel()

		var (
			err     any
			panicFn = func() {
				defer func() {
					err = recover()
				}()

				assert.Assertf(false, "something went wrong: key=%s", "value")
			}
		)

		panicFn()
		check.Equal(t, "something went wrong: key=value", err)
	})
}

func TestNotErrorf(t *testing.T) {
	t.Parallel()

	t.Run("not an error", func(t *testing.T) {
		t.Parallel()

		assert.NotErrf(nil, "")
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		var (
			panicErr any
			err      = errors.New("some_error") //nolint:err113
			panicFn  = func() {
				defer func() {
					panicErr = recover()
				}()

				assert.NotErrf(err, "something went wrong: key=%s", "value")
			}
		)

		panicFn()
		check.Equal(t, "something went wrong: key=value\nsome_error", panicErr)
	})
}
