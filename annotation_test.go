// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the GPLv3, see LICENCE file for details.

package errors_test

import (
	"fmt"
	"os"
	"path/filepath"

	gc "launchpad.net/gocheck"

	"github.com/juju/errors"
)

type annotationSuite struct{}

var _ = gc.Suite(&annotationSuite{})

func echo(value interface{}) interface{} {
	return value
}

func (*annotationSuite) TestErrorString(c *gc.C) {
	for i, test := range []struct {
		message   string
		generator func() error
		expected  string
	}{
		{
			message: "uncomparable errors",
			generator: func() error {
				err := errors.Annotatef(newNonComparableError("uncomparable"), "annotation")
				return errors.Annotatef(err, "another")
			},
			expected: "another: annotation: uncomparable",
		}, {
			message: "Errorf",
			generator: func() error {
				return errors.Errorf("first error")
			},
			expected: "first error",
		}, {
			message: "annotating nil",
			generator: func() error {
				return errors.Annotatef(nil, "annotation")
			},
			expected: "annotation",
		}, {
			message: "annotated error",
			generator: func() error {
				err := errors.Errorf("first error")
				return errors.Annotatef(err, "annotation")
			},
			expected: "annotation: first error",
		}, {
			message: "test annotation format",
			generator: func() error {
				err := errors.Errorf("first %s", "error")
				return errors.Annotatef(err, "%s", "annotation")
			},
			expected: "annotation: first error",
		}, {
			message: "wrapped error",
			generator: func() error {
				//first := Errorf("first error")
				err := newError("first error")
				return errors.Wrap(err, newError("detailed error"))
			},
			expected: "detailed error",
		}, {
			message: "wrapped annotated error",
			generator: func() error {
				err := errors.Errorf("first error")
				err = errors.Annotatef(err, "annotated")
				return errors.Wrap(err, fmt.Errorf("detailed error"))
			},
			expected: "detailed error",
		}, {
			message: "annotated wrapped error",
			generator: func() error {
				err := errors.Errorf("first error")
				err = errors.Wrap(err, fmt.Errorf("detailed error"))
				return errors.Annotatef(err, "annotated")
			},
			expected: "annotated: detailed error",
		},
	} {
		c.Logf("%v: %s", i, test.message)
		err := test.generator()
		ok := c.Check(err.Error(), gc.Equals, test.expected)
		if !ok {
			c.Logf("%#v", test.generator())
		}
	}
}

func (*annotationSuite) TestAnnotatedErrorCheck(c *gc.C) {
	// Look for a file that we know isn't there.
	dir := c.MkDir()
	_, err := os.Stat(filepath.Join(dir, "not-there"))
	c.Assert(os.IsNotExist(err), gc.Equals, true)
	c.Assert(errors.Check(err, os.IsNotExist), gc.Equals, true)

	err = errors.Annotatef(err, "wrap it")
	// Now the error itself isn't a 'IsNotExist'.
	c.Assert(os.IsNotExist(err), gc.Equals, false)
	// However if we use the Check method, it is.
	c.Assert(errors.Check(err, os.IsNotExist), gc.Equals, true)
}

// This is an uncomparable error type, as it is a struct that supports the
// error interface (as opposed to a pointer type).
type error_ struct {
	info  string
	slice []string
}

// Create a non-comparable error
func newNonComparableError(message string) error {
	return error_{info: message}
}

func (e error_) Error() string {
	return e.info
}

func newError(message string) error {
	return testError{message}
}

// The testError is a value type error for ease of seeing results
// when the test fails.
type testError struct {
	message string
}

func (e testError) Error() string {
	return e.message
}
