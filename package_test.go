// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"fmt"
	"runtime"
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

// errorLocationValue provides the function name and line number for where this
// function was called from - 1 line. What this means is that the returned value
// will be homed to the file line directly above where this function was called.
// This is a utility for testing error details and that associated error calls
// set the error location correctly.
func errorLocationValue(c *gc.C) string {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(2, rpc[:])
	if n < 1 {
		return ""
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return fmt.Sprintf("%s:%d", frame.Function, frame.Line-1)
}
