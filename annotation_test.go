package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gc "launchpad.net/gocheck"
)

func Test(t *testing.T) { gc.TestingT(t) }

type arrarSuite struct{}

var _ = gc.Suite(&arrarSuite{})

func echo(value interface{}) interface{} {
	return value
}

func (*arrarSuite) TestErrorString(c *gc.C) {
	for i, test := range []struct {
		message   string
		generator func() error
		expected  string
	}{
		{
			message: "uncomparable errors",
			generator: func() error {
				err := Annotatef(newNonComparableError("uncomparable"), "annotation")
				return Annotatef(err, "another")
			},
			expected: "another, annotation: uncomparable",
		}, {
			message: "Errorf",
			generator: func() error {
				return Errorf("first error")
			},
			expected: "first error",
		}, {
			message: "annotating nil",
			generator: func() error {
				return Annotatef(nil, "annotation")
			},
			expected: "annotation",
		}, {
			message: "annotated error",
			generator: func() error {
				err := Errorf("first error")
				return Annotatef(err, "annotation")
			},
			expected: "annotation: first error",
		}, {
			message: "test annotation format",
			generator: func() error {
				err := Errorf("first %s", "error")
				return Annotatef(err, "%s", "annotation")
			},
			expected: "annotation: first error",
		}, {
			message: "wrapped error",
			generator: func() error {
				//first := Errorf("first error")
				err := newError("first error")
				return Wrap(err, newError("more detail"))
			},
			expected: "more detail (first error)",
		}, {
			message: "wrapped annotated error",
			generator: func() error {
				err := Errorf("first error")
				err = Annotatef(err, "annotated")
				return Wrap(err, fmt.Errorf("more detail"))
			},
			expected: "more detail (annotated: first error)",
		}, {
			message: "annotated wrapped error",
			generator: func() error {
				err := Errorf("first error")
				err = Wrap(err, fmt.Errorf("more detail"))
				return Annotatef(err, "annotated")
			},
			expected: "annotated: more detail (first error)",
		}, {
			message: "annotated wrapped annotated error",
			generator: func() error {
				err := Errorf("first error")
				err = Annotatef(err, "annotated")
				err = Wrap(err, fmt.Errorf("more detail"))
				return Annotatef(err, "context")
			},
			expected: "context: more detail (annotated: first error)",
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

func (*arrarSuite) TestAnnotatedErrorCheck(c *gc.C) {
	// Look for a file that we know isn't there.
	dir := c.MkDir()
	_, err := os.Stat(filepath.Join(dir, "not-there"))
	c.Assert(os.IsNotExist(err), gc.Equals, true)
	c.Assert(Check(err, os.IsNotExist), gc.Equals, true)

	err = Annotatef(err, "wrap it")
	// Now the error itself isn't a 'IsNotExist'.
	c.Assert(os.IsNotExist(err), gc.Equals, false)
	// However if we use the Check method, it is.
	c.Assert(Check(err, os.IsNotExist), gc.Equals, true)
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
