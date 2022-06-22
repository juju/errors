// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	gc "gopkg.in/check.v1"

	"github.com/juju/errors"
)

type functionSuite struct {
}

var _ = gc.Suite(&functionSuite{})

func (*functionSuite) TestNew(c *gc.C) {
	err := errors.New("testing")
	loc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "testing")
	c.Assert(errors.Cause(err), gc.Equals, err)
	c.Assert(errors.Details(err), Contains, loc)
}

func (*functionSuite) TestErrorf(c *gc.C) {
	err := errors.Errorf("testing %d", 42)
	loc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "testing 42")
	c.Assert(errors.Cause(err), gc.Equals, err)
	c.Assert(errors.Details(err), Contains, loc)
}

func (*functionSuite) TestTrace(c *gc.C) {
	first := errors.New("first")
	err := errors.Trace(first)
	loc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "first")
	c.Assert(errors.Is(err, first), gc.Equals, true)
	c.Assert(errors.Details(err), Contains, loc)

	c.Assert(errors.Trace(nil), gc.IsNil)
}

func (*functionSuite) TestAnnotate(c *gc.C) {
	first := errors.New("first")
	err := errors.Annotate(first, "annotation")
	loc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "annotation: first")
	c.Assert(errors.Cause(err), gc.Equals, first)
	c.Assert(errors.Details(err), Contains, loc)

	c.Assert(errors.Annotate(nil, "annotate"), gc.IsNil)
}

func (*functionSuite) TestAnnotatef(c *gc.C) {
	first := errors.New("first")
	err := errors.Annotatef(first, "annotation %d", 2) //err annotatefTest
	loc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "annotation 2: first")
	c.Assert(errors.Cause(err), gc.Equals, first)
	c.Assert(errors.Details(err), Contains, loc)

	c.Assert(errors.Annotatef(nil, "annotate"), gc.IsNil)
}

func (*functionSuite) TestDeferredAnnotatef(c *gc.C) {
	// NOTE: this test fails with gccgo
	if runtime.Compiler == "gccgo" {
		c.Skip("gccgo can't determine the location")
	}
	first := errors.New("first")
	test := func() (err error) {
		defer errors.DeferredAnnotatef(&err, "deferred %s", "annotate")
		return first
	}
	err := test()
	c.Assert(err.Error(), gc.Equals, "deferred annotate: first")
	c.Assert(errors.Cause(err), gc.Equals, first)

	err = nil
	errors.DeferredAnnotatef(&err, "deferred %s", "annotate")
	c.Assert(err, gc.IsNil)
}

func (*functionSuite) TestWrap(c *gc.C) {
	first := errors.New("first")
	firstLoc := errorLocationValue(c)

	detailed := errors.New("detailed")
	err := errors.Wrap(first, detailed)
	secondLoc := errorLocationValue(c)

	c.Assert(err.Error(), gc.Equals, "detailed")
	c.Assert(errors.Cause(err), gc.Equals, detailed)
	c.Assert(errors.Details(err), Contains, firstLoc)
	c.Assert(errors.Details(err), Contains, secondLoc)
}

func (*functionSuite) TestWrapOfNil(c *gc.C) {
	detailed := errors.New("detailed")
	err := errors.Wrap(nil, detailed)
	loc := errorLocationValue(c)
	c.Assert(err.Error(), gc.Equals, "detailed")
	c.Assert(errors.Cause(err), gc.Equals, detailed)
	c.Assert(errors.Details(err), Contains, loc)
}

func (*functionSuite) TestWrapf(c *gc.C) {
	first := errors.New("first")
	firstLoc := errorLocationValue(c)
	detailed := errors.New("detailed")
	err := errors.Wrapf(first, detailed, "value %d", 42)
	secondLoc := errorLocationValue(c)
	c.Assert(err.Error(), gc.Equals, "value 42: detailed")
	c.Assert(errors.Cause(err), gc.Equals, detailed)
	c.Assert(errors.Details(err), Contains, firstLoc)
	c.Assert(errors.Details(err), Contains, secondLoc)
}

func (*functionSuite) TestWrapfOfNil(c *gc.C) {
	detailed := errors.New("detailed")
	err := errors.Wrapf(nil, detailed, "value %d", 42)
	loc := errorLocationValue(c)
	c.Assert(err.Error(), gc.Equals, "value 42: detailed")
	c.Assert(errors.Cause(err), gc.Equals, detailed)
	c.Assert(errors.Details(err), Contains, loc)
}

func (*functionSuite) TestMask(c *gc.C) {
	first := errors.New("first")
	err := errors.Mask(first)
	loc := errorLocationValue(c)
	c.Assert(err.Error(), gc.Equals, "first")
	c.Assert(errors.Cause(err), gc.Equals, err)
	c.Assert(errors.Details(err), Contains, loc)

	c.Assert(errors.Mask(nil), gc.IsNil)
}

func (*functionSuite) TestMaskf(c *gc.C) {
	first := errors.New("first")
	err := errors.Maskf(first, "masked %d", 42)
	loc := errorLocationValue(c)
	c.Assert(err.Error(), gc.Equals, "masked 42: first")
	c.Assert(errors.Cause(err), gc.Equals, err)
	c.Assert(errors.Details(err), Contains, loc)

	c.Assert(errors.Maskf(nil, "mask"), gc.IsNil)
}

func (*functionSuite) TestCause(c *gc.C) {
	c.Assert(errors.Cause(nil), gc.IsNil)
	c.Assert(errors.Cause(someErr), gc.Equals, someErr)

	fmtErr := fmt.Errorf("simple")
	c.Assert(errors.Cause(fmtErr), gc.Equals, fmtErr)

	err := errors.Wrap(someErr, fmtErr)
	c.Assert(errors.Cause(err), gc.Equals, fmtErr)

	err = errors.Annotate(err, "annotated")
	c.Assert(errors.Cause(err), gc.Equals, fmtErr)

	err = errors.Maskf(err, "masked")
	c.Assert(errors.Cause(err), gc.Equals, err)

	// Look for a file that we know isn't there.
	dir := c.MkDir()
	_, err = os.Stat(filepath.Join(dir, "not-there"))
	c.Assert(os.IsNotExist(err), gc.Equals, true)

	err = errors.Annotatef(err, "wrap it")
	// Now the error itself isn't a 'IsNotExist'.
	c.Assert(os.IsNotExist(err), gc.Equals, false)
	// However if we use the Check method, it is.
	c.Assert(os.IsNotExist(errors.Cause(err)), gc.Equals, true)
}

type tracer interface {
	StackTrace() []string
}

func (*functionSuite) TestErrorStack(c *gc.C) {
	for i, test := range []struct {
		message   string
		generator func(*gc.C, io.Writer) error
		tracer    bool
	}{{
		message: "nil",
		generator: func(_ *gc.C, _ io.Writer) error {
			return nil
		},
	}, {
		message: "raw error",
		generator: func(c *gc.C, expected io.Writer) error {
			fmt.Fprint(expected, "raw")
			return fmt.Errorf("raw")
		},
	}, {
		message: "single error stack",
		generator: func(c *gc.C, expected io.Writer) error {
			err := errors.New("first error")
			fmt.Fprintf(expected, "%s: first error", errorLocationValue(c))
			return err
		},
		tracer: true,
	}, {
		message: "annotated error",
		generator: func(c *gc.C, expected io.Writer) error {
			err := errors.New("first error")
			fmt.Fprintf(expected, "%s: first error\n", errorLocationValue(c))
			err = errors.Annotate(err, "annotation")
			fmt.Fprintf(expected, "%s: annotation", errorLocationValue(c))
			return err
		},
		tracer: true,
	}, {
		message: "wrapped error",
		generator: func(c *gc.C, expected io.Writer) error {
			err := errors.New("first error")
			fmt.Fprintf(expected, "%s: first error\n", errorLocationValue(c))
			err = errors.Wrap(err, newError("detailed error"))
			fmt.Fprintf(expected, "%s: detailed error", errorLocationValue(c))
			return err
		},
		tracer: true,
	}, {
		message: "annotated wrapped error",
		generator: func(c *gc.C, expected io.Writer) error {
			err := errors.Errorf("first error")
			fmt.Fprintf(expected, "%s: first error\n", errorLocationValue(c))
			err = errors.Wrap(err, fmt.Errorf("detailed error"))
			fmt.Fprintf(expected, "%s: detailed error\n", errorLocationValue(c))
			err = errors.Annotatef(err, "annotated")
			fmt.Fprintf(expected, "%s: annotated", errorLocationValue(c))
			return err
		},
		tracer: true,
	}, {
		message: "traced, and annotated",
		generator: func(c *gc.C, expected io.Writer) error {
			err := errors.New("first error")
			fmt.Fprintf(expected, "%s: first error\n", errorLocationValue(c))
			err = errors.Trace(err)
			fmt.Fprintf(expected, "%s: \n", errorLocationValue(c))
			err = errors.Annotate(err, "some context")
			fmt.Fprintf(expected, "%s: some context\n", errorLocationValue(c))
			err = errors.Trace(err)
			fmt.Fprintf(expected, "%s: \n", errorLocationValue(c))
			err = errors.Annotate(err, "more context")
			fmt.Fprintf(expected, "%s: more context\n", errorLocationValue(c))
			err = errors.Trace(err)
			fmt.Fprintf(expected, "%s: ", errorLocationValue(c))
			return err
		},
		tracer: true,
	}, {
		message: "uncomparable, wrapped with a value error",
		generator: func(c *gc.C, expected io.Writer) error {
			err := newNonComparableError("first error")
			fmt.Fprintln(expected, "first error")
			err = errors.Trace(err)
			fmt.Fprintf(expected, "%s: \n", errorLocationValue(c))
			err = errors.Wrap(err, newError("value error"))
			fmt.Fprintf(expected, "%s: value error\n", errorLocationValue(c))
			err = errors.Maskf(err, "masked")
			fmt.Fprintf(expected, "%s: masked\n", errorLocationValue(c))
			err = errors.Annotate(err, "more context")
			fmt.Fprintf(expected, "%s: more context\n", errorLocationValue(c))
			err = errors.Trace(err)
			fmt.Fprintf(expected, "%s: ", errorLocationValue(c))
			return err
		},
		tracer: true,
	}} {
		c.Logf("%v: %s", i, test.message)
		expected := strings.Builder{}
		err := test.generator(c, &expected)
		stack := errors.ErrorStack(err)
		ok := c.Check(stack, gc.Equals, expected.String())
		if !ok {
			c.Logf("%#v", err)
		}
		tracer, ok := err.(tracer)
		c.Check(ok, gc.Equals, test.tracer)
		if ok {
			stackTrace := tracer.StackTrace()
			c.Check(stackTrace, gc.DeepEquals, strings.Split(stack, "\n"))
		}
	}
}

func (*functionSuite) TestFormat(c *gc.C) {
	formatErrorExpected := &strings.Builder{}
	err := errors.New("TestFormat")
	fmt.Fprintf(formatErrorExpected, "%s: TestFormat\n", errorLocationValue(c))
	err = errors.Mask(err)
	fmt.Fprintf(formatErrorExpected, "%s: ", errorLocationValue(c))

	for i, test := range []struct {
		format string
		expect string
	}{{
		format: "%s",
		expect: "TestFormat",
	}, {
		format: "%v",
		expect: "TestFormat",
	}, {
		format: "%q",
		expect: `"TestFormat"`,
	}, {
		format: "%A",
		expect: `%!A(*errors.Err=TestFormat)`,
	}, {
		format: "%+v",
		expect: formatErrorExpected.String(),
	}} {
		c.Logf("test %d: %q", i, test.format)
		s := fmt.Sprintf(test.format, err)
		c.Check(s, gc.Equals, test.expect)
	}
}

type basicError struct {
	Reason string
}

func (b *basicError) Error() string {
	return b.Reason
}

func (*functionSuite) TestAs(c *gc.C) {
	baseError := &basicError{"I'm an error"}
	testErrors := []error{
		errors.Trace(baseError),
		errors.Annotate(baseError, "annotation"),
		errors.Wrap(baseError, errors.New("wrapper")),
		errors.Mask(baseError),
	}

	for _, err := range testErrors {
		bError := &basicError{}
		val := errors.As(err, &bError)
		c.Check(val, gc.Equals, true)
		c.Check(bError.Reason, gc.Equals, "I'm an error")
	}
}

func (*functionSuite) TestIs(c *gc.C) {
	baseError := &basicError{"I'm an error"}
	testErrors := []error{
		errors.Trace(baseError),
		errors.Annotate(baseError, "annotation"),
		errors.Wrap(baseError, errors.New("wrapper")),
		errors.Mask(baseError),
	}

	for _, err := range testErrors {
		val := errors.Is(err, baseError)
		c.Check(val, gc.Equals, true)
	}
}

func (*functionSuite) TestSetLocationWithNilError(c *gc.C) {
	c.Assert(errors.SetLocation(nil, 1), gc.IsNil)
}

func (*functionSuite) TestSetLocation(c *gc.C) {
	err := errors.New("test")
	err = errors.SetLocation(err, 1)
	stack := fmt.Sprintf("%s: test", errorLocationValue(c))
	_, implements := err.(errors.Locationer)
	c.Assert(implements, gc.Equals, true)

	c.Check(errors.ErrorStack(err), gc.Equals, stack)
}

func (*functionSuite) TestHideErrorStillReturnsErrorString(c *gc.C) {
	err := stderrors.New("This is a simple error")
	err = errors.Hide(err)

	c.Assert(err.Error(), gc.Equals, "This is a simple error")
}

func (*functionSuite) TestQuietWrappedErrorStillSatisfied(c *gc.C) {
	simpleTestError := errors.ConstError("I am a teapot")
	err := fmt.Errorf("fill me up%w", errors.Hide(simpleTestError))
	c.Assert(err.Error(), gc.Equals, "fill me up")
	c.Assert(errors.Is(err, simpleTestError), gc.Equals, true)
}

type ComplexErrorMessage interface {
	error
	ComplexMessage() string
}

type complexError struct {
	Message string
}

func (c *complexError) Error() string {
	return c.Message
}

func (c *complexError) ComplexMessage() string {
	return c.Message
}

type complexErrorOther struct {
	Message string
}

func (c *complexErrorOther) As(e any) bool {
	if ce, ok := e.(**complexError); ok {
		*ce = &complexError{
			Message: c.Message,
		}
		return true
	}
	return false
}

func (c *complexErrorOther) Error() string {
	return c.Message
}

func (c *complexErrorOther) ComplexMessage() string {
	return c.Message
}

func (*functionSuite) TestHasType(c *gc.C) {
	complexErr := &complexError{Message: "complex error message"}
	wrapped1 := fmt.Errorf("wrapping1: %w", complexErr)
	wrapped2 := fmt.Errorf("wrapping2: %w", wrapped1)

	c.Assert(errors.HasType[*complexError](complexErr), gc.Equals, true)
	c.Assert(errors.HasType[*complexError](wrapped1), gc.Equals, true)
	c.Assert(errors.HasType[*complexError](wrapped2), gc.Equals, true)
	c.Assert(errors.HasType[ComplexErrorMessage](wrapped2), gc.Equals, true)
	c.Assert(errors.HasType[*complexErrorOther](wrapped2), gc.Equals, false)
	c.Assert(errors.HasType[*complexErrorOther](nil), gc.Equals, false)

	complexErrOther := &complexErrorOther{Message: "another complex error"}

	c.Assert(errors.HasType[*complexError](complexErrOther), gc.Equals, true)

	wrapped2 = fmt.Errorf("wrapping1: %w", complexErrOther)
	c.Assert(errors.HasType[*complexError](wrapped2), gc.Equals, true)
}

func (*functionSuite) TestAsType(c *gc.C) {
	complexErr := &complexError{Message: "complex error message"}
	wrapped1 := fmt.Errorf("wrapping1: %w", complexErr)
	wrapped2 := fmt.Errorf("wrapping2: %w", wrapped1)

	ce, ok := errors.AsType[*complexError](complexErr)
	c.Assert(ok, gc.Equals, true)
	c.Assert(ce.Message, gc.Equals, complexErr.Message)

	ce, ok = errors.AsType[*complexError](wrapped1)
	c.Assert(ok, gc.Equals, true)
	c.Assert(ce.Message, gc.Equals, complexErr.Message)

	ce, ok = errors.AsType[*complexError](wrapped2)
	c.Assert(ok, gc.Equals, true)
	c.Assert(ce.Message, gc.Equals, complexErr.Message)

	cem, ok := errors.AsType[ComplexErrorMessage](wrapped2)
	c.Assert(ok, gc.Equals, true)
	c.Assert(cem.ComplexMessage(), gc.Equals, complexErr.Message)

	ceo, ok := errors.AsType[*complexErrorOther](wrapped2)
	c.Assert(ok, gc.Equals, false)
	c.Assert(ceo, gc.Equals, (*complexErrorOther)(nil))

	ceo, ok = errors.AsType[*complexErrorOther](nil)
	c.Assert(ok, gc.Equals, false)
	c.Assert(ceo, gc.Equals, (*complexErrorOther)(nil))

	complexErrOther := &complexErrorOther{Message: "another complex error"}
	ce, ok = errors.AsType[*complexError](complexErrOther)
	c.Assert(ok, gc.Equals, true)
	c.Assert(ce.Message, gc.Equals, complexErrOther.Message)

	wrapped2 = fmt.Errorf("wrapping1: %w", complexErrOther)
	ce, ok = errors.AsType[*complexError](wrapped2)
	c.Assert(ok, gc.Equals, true)
	c.Assert(ce.Message, gc.Equals, complexErrOther.Message)
}

func ExampleHide() {
	myConstError := errors.ConstError("I don't want to be fmt printed")
	err := fmt.Errorf("don't show this error%w", errors.Hide(myConstError))

	fmt.Println(err)
	fmt.Println(stderrors.Is(err, myConstError))

	// Output:
	// don't show this error
	// true
}

type MyError struct {
	Message string
}

func (m *MyError) Error() string {
	return m.Message
}

func ExampleHasType() {
	myErr := &MyError{Message: "these are not the droids you're looking for"}
	err := fmt.Errorf("wrapped: %w", myErr)
	is := errors.HasType[*MyError](err)
	fmt.Println(is)

	// Output:
	// true
}

func ExampleAsType() {
	myErr := &MyError{Message: "these are not the droids you're looking for"}
	err := fmt.Errorf("wrapped: %w", myErr)
	myErr, is := errors.AsType[*MyError](err)
	fmt.Println(is)
	fmt.Println(myErr.Message)

	// Output:
	// true
	// these are not the droids you're looking for
}
