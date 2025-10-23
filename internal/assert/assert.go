// Package assert contains helper functions for asserting invariances etc. It is
// a shame there is not builtin assert in golang.
package assert

import (
	"fmt"
)

// Assertf panics if the given condition is not true. It will panic with the
// given format and args after calling [fmt.Sprintf].
func Assertf(condition bool, format string, args ...any) {
	if condition == false {
		msg := fmt.Sprintf(format, args...)
		panic(msg)
	}
}

// NotErrf panics if the given error is not nil. It will panic with the given
// format and args after calling [fmt.Sprintf]. It will also print the result of
// `err.Error()` on a newline after the formatted message.
func NotErrf(err error, format string, args ...any) {
	if err != nil {
		msg := fmt.Sprintf(format, args...) + "\n" + err.Error()
		panic(msg)
	}
}
