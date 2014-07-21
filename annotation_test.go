// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	gc "launchpad.net/gocheck"

	"github.com/juju/errgo"
	"github.com/juju/errors"
)

type annotationSuite struct{}

var _ = gc.Suite(&annotationSuite{})

func echo(value interface{}) interface{} {
	return value
}

func (*annotationSuite) TestNilArgs(c *gc.C) {
	c.Assert(errors.Trace(nil), gc.IsNil)
	c.Assert(errors.Annotate(nil, "foo"), gc.IsNil)
	c.Assert(errors.Annotatef(nil, "foo %d", 2), gc.IsNil)
	c.Assert(errors.Wrap(nil, errors.New("omg")), gc.IsNil)
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
		}, {
			message: "traced, and annotated",
			generator: func() error {
				err := errors.New("first error")
				err = errors.Trace(err)
				err = errors.Annotate(err, "some context")
				err = errors.Trace(err)
				err = errors.Annotate(err, "more context")
				return errors.Trace(err)
			},
			expected: "more context: some context: first error",
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

func (*annotationSuite) TestErrorStack(c *gc.C) {
	for i, test := range []struct {
		message   string
		generator func() error
		expected  string
	}{
		{
			message: "raw error",
			generator: func() error {
				return fmt.Errorf("raw")
			},
			expected: "raw",
		}, {
			message: "single error stack",
			generator: func() error {
				return errors.New("first error") //err single
			},
			expected: "$single$: first error",
		}, {
			message: "annotated error",
			generator: func() error {
				err := errors.New("first error")          //err annotated-0
				return errors.Annotate(err, "annotation") //err annotated-1
			},
			expected: "" +
				"$annotated-0$: first error\n" +
				"$annotated-1$: annotation",
		}, {
			message: "wrapped error",
			generator: func() error {
				err := errors.New("first error")                    //err wrapped-0
				return errors.Wrap(err, newError("detailed error")) //err wrapped-1
			},
			expected: "" +
				"$wrapped-0$: first error\n" +
				"$wrapped-1$: detailed error",
		}, {
			message: "annotated wrapped error",
			generator: func() error {
				err := errors.Errorf("first error")                  //err ann-wrap-0
				err = errors.Wrap(err, fmt.Errorf("detailed error")) //err ann-wrap-1
				return errors.Annotatef(err, "annotated")            //err ann-wrap-2
			},
			expected: "" +
				"$ann-wrap-0$: first error\n" +
				"$ann-wrap-1$: detailed error\n" +
				"$ann-wrap-2$: annotated",
		}, {
			message: "traced, and annotated",
			generator: func() error {
				err := errors.New("first error")           //err stack-0
				err = errors.Trace(err)                    //err stack-1
				err = errors.Annotate(err, "some context") //err stack-2
				err = errors.Trace(err)                    //err stack-3
				err = errors.Annotate(err, "more context") //err stack-4
				return errors.Trace(err)                   //err stack-5
			},
			expected: "" +
				"$stack-0$: first error\n" +
				"$stack-1$: \n" +
				"$stack-2$: some context\n" +
				"$stack-3$: \n" +
				"$stack-4$: more context\n" +
				"$stack-5$: ",
		}, {
			message: "uncomparable, mixed errgo, value error",
			generator: func() error {
				err := newNonComparableError("first error")                        //err mixed-0
				err = errors.Trace(err)                                            //err mixed-1
				err = errgo.WithCausef(err, newError("value error"), "annotation") //err mixed-2
				err = errors.Trace(err)                                            //err mixed-3
				err = errors.Annotate(err, "more context")                         //err mixed-4
				return errors.Trace(err)                                           //err mixed-5
			},
			expected: "" +
				"first error\n" +
				"$mixed-1$: \n" +
				"$mixed-2$: annotation: value error\n" +
				"$mixed-3$: \n" +
				"$mixed-4$: more context\n" +
				"$mixed-5$: ",
		},
	} {
		c.Logf("%v: %s", i, test.message)
		err := test.generator()
		expected := replaceLocations(test.expected)
		ok := c.Check(errors.ErrorStack(err), gc.Equals, expected)
		if !ok {
			c.Logf("%#v", err)
		}
	}
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

// Copied from the errgo/errors_test.go.

func replaceLocations(s string) string {
	t := ""
	for {
		i := strings.Index(s, "$")
		if i == -1 {
			break
		}
		t += s[0:i]
		s = s[i+1:]
		i = strings.Index(s, "$")
		if i == -1 {
			panic("no second $")
		}
		t += location(s[0:i]).String()
		s = s[i+1:]
	}
	t += s
	return t
}

func location(tag string) errgo.Location {
	line, ok := tagToLine[tag]
	if !ok {
		panic(fmt.Errorf("tag %q not found", tag))
	}
	return errgo.Location{
		File: "github.com/juju/errors/annotation_test.go",
		Line: line,
	}
}

var tagToLine = make(map[string]int)

func init() {
	data, err := ioutil.ReadFile("annotation_test.go")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if j := strings.Index(line, "//err "); j >= 0 {
			tagToLine[line[j+len("//err "):]] = i + 1
		}
	}
}
