// Copyright 2023 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"fmt"

	"github.com/juju/errors"
)

func ExampleTrace() {
	err := fmt.Errorf("Too many gophers to count")
	tracedErr := errors.Trace(err)

	fmt.Println(tracedErr)
	fmt.Println(errors.Is(tracedErr, err))
	fmt.Println(errors.ErrorStack(tracedErr))
	fmt.Println()

	tracedErr = errors.Trace(tracedErr)
	fmt.Println(errors.ErrorStack(tracedErr))
	fmt.Println()

	tracedErr = errors.Trace(fmt.Errorf("foobar: %w", tracedErr))
	fmt.Println(errors.ErrorStack(tracedErr))

	// Output: Too many gophers to count
	// true
	// Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:14: Too many gophers to count
	//
	// Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:14: Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:21: Too many gophers to count
	//
	// Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:14: Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:21: Too many gophers to count
	// foobar: Too many gophers to count
	// github.com/juju/errors_test.ExampleTrace:25: foobar: Too many gophers to count
}
