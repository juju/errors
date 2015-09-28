// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type multiSuite struct{}

var _ = gc.Suite(&multiSuite{})

func (*multiSuite) TestNewMultiError(c *gc.C) {
	err, _ := errors.NewMultiError()
	errs, ids := errors.ExposeMultiError(err)

	c.Check(errs, gc.HasLen, 0)
	c.Check(ids, gc.HasLen, 0)
}

func (*multiSuite) TestSetErrorOkay(c *gc.C) {
	err, setError := errors.NewMultiError()
	expectedA := errors.Errorf("<failure>")
	expectedB := errors.Errorf("<failure>")
	setError(expectedB, "b")
	setError(expectedA, "a")
	errs, ids := errors.ExposeMultiError(err)

	c.Check(errs, jc.DeepEquals, []error{expectedB, expectedA})
	c.Check(ids, jc.DeepEquals, []string{"b", "a"})
}

func (*multiSuite) TestSetErrorNoID(c *gc.C) {
	err, setError := errors.NewMultiError()
	expected := errors.Errorf("<failure>")
	setError(expected, "")
	errs, ids := errors.ExposeMultiError(err)

	c.Check(errs, jc.DeepEquals, []error{expected})
	c.Check(ids, jc.DeepEquals, []string{""})
}

func (*multiSuite) TestSetErrorAlreadySet(c *gc.C) {
	err, setError := errors.NewMultiError()
	first := errors.Errorf("<first>")
	setError(first, "b")
	second := errors.Errorf("<failure>")
	setError(second, "b")
	errs, ids := errors.ExposeMultiError(err)

	c.Check(errs, jc.DeepEquals, []error{first, second})
	c.Check(ids, jc.DeepEquals, []string{"b", "b"})
}

func (*multiSuite) TestErrorStringSameIDs(c *gc.C) {
	err, setError := errors.NewMultiError()
	for _, id := range []string{"b", "b", "b"} {
		setError(errors.Errorf(id), id)
	}

	c.Check(err, gc.ErrorMatches, `3 errors: .*`)
}

func (*multiSuite) TestErrorStringMixedIDs(c *gc.C) {
	err, setError := errors.NewMultiError()
	for _, id := range []string{"a", "b", "", "b"} {
		setError(errors.Errorf(id), id)
	}

	c.Check(err, gc.ErrorMatches, `4 errors \(for 3 IDs\): .*`)
}

func (*multiSuite) TestErrorStringUniqueIDs(c *gc.C) {
	err, setError := errors.NewMultiError()
	for _, id := range []string{"a", "b", "c"} {
		setError(errors.Errorf(id), id)
	}

	c.Check(err, gc.ErrorMatches, `3 errors \(for 3 IDs\): .*`)
}

func (*multiSuite) TestErrorStringNoIDs(c *gc.C) {
	err, setError := errors.NewMultiError()
	for _, id := range []string{"", "", ""} {
		setError(errors.Errorf(id), id)
	}

	c.Check(err, gc.ErrorMatches, `3 errors: .*`)
}

func (*multiSuite) TestErrorStringEmpty(c *gc.C) {
	err, _ := errors.NewMultiError()

	c.Check(err, gc.ErrorMatches, `0 errors`)
}

func (*multiSuite) TestErrorsOkay(c *gc.C) {
	err, setError := errors.NewMultiError()
	expectedA := errors.Errorf("<failure>")
	expectedB := errors.Errorf("<failure>")
	setError(expectedB, "b")
	setError(expectedA, "a")
	errs, ids := err.Errors()

	c.Check(errs, jc.DeepEquals, []error{expectedB, expectedA})
	c.Check(ids, jc.DeepEquals, []string{"b", "a"})
}

func (*multiSuite) TestErrorsMixed(c *gc.C) {
	expectedIDs := []string{"b", "a", "", "b", "c"}
	var expectedErrs []error
	err, setError := errors.NewMultiError()
	for _, id := range expectedIDs {
		expectedErr := errors.Errorf(id)
		expectedErrs = append(expectedErrs, expectedErr)
		setError(expectedErr, id)
	}
	errs, ids := err.Errors()

	c.Check(errs, jc.DeepEquals, expectedErrs)
	c.Check(ids, jc.DeepEquals, expectedIDs)
}

func (*multiSuite) TestErrorsEmpty(c *gc.C) {
	err, _ := errors.NewMultiError()
	errs, ids := err.Errors()

	c.Check(errs, gc.HasLen, 0)
	c.Check(ids, gc.HasLen, 0)
}

func (*multiSuite) TestIsMultiErrorDirect(c *gc.C) {
	err, _ := errors.NewMultiError()
	ok := errors.IsMultiError(err)

	c.Check(ok, jc.IsTrue)
}

func (*multiSuite) TestIsMultiErrorIndirect(c *gc.C) {
	err, _ := errors.NewMultiError()
	wrapped := errors.Trace(err)
	ok := errors.IsMultiError(wrapped)

	c.Check(ok, jc.IsTrue)
}

func (*multiSuite) TestIsMultiErrorFalse(c *gc.C) {
	err := errors.Errorf("not multi")
	ok := errors.IsMultiError(err)

	c.Check(ok, jc.IsFalse)
}
