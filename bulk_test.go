// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type bulkSuite struct{}

var _ = gc.Suite(&bulkSuite{})

func (*bulkSuite) TestNewBulkErrorOkay(c *gc.C) {
	err, _ := errors.NewBulkError("a", "b", "c")
	ids := err.IDs()

	c.Check(ids, jc.DeepEquals, []string{"a", "b", "c"})
}

func (*bulkSuite) TestNewBulkErrorEmpty(c *gc.C) {
	err, _ := errors.NewBulkError()
	ids := err.IDs()

	c.Check(ids, gc.HasLen, 0)
}

func (*bulkSuite) TestSetErrorOkay(c *gc.C) {
	err, setError := errors.NewBulkError("a", "b", "c")
	expected := errors.Errorf("<failure>")
	ok := setError("b", expected)
	errors := err.Enumerate()

	c.Check(ok, jc.IsTrue)
	c.Check(errors, jc.DeepEquals, []error{nil, expected, nil})
}

func (*bulkSuite) TestSetErrorAlreadySet(c *gc.C) {
	err, setError := errors.NewBulkError("a", "b", "c")
	ok := setError("b", errors.Errorf("<first>"))
	c.Assert(ok, jc.IsTrue)
	expected := errors.Errorf("<failure>")
	ok = setError("b", expected)
	errors := err.Enumerate()

	c.Check(ok, jc.IsTrue)
	c.Check(errors, jc.DeepEquals, []error{nil, expected, nil})
}

func (*bulkSuite) TestSetErrorUnrecognized(c *gc.C) {
	err, setError := errors.NewBulkError("a", "b", "c")
	ok := setError("d", errors.Errorf("<failure>"))
	errors := err.Enumerate()

	c.Check(ok, jc.IsFalse)
	c.Check(errors, jc.DeepEquals, []error{nil, nil, nil})
}

func (*bulkSuite) TestErrorStringFull(c *gc.C) {
	ids := []string{"a", "b", "c"}
	expected := []error{
		errors.Errorf("a"),
		errors.Errorf("b"),
		errors.Errorf("c"),
	}
	err, setError := errors.NewBulkError(ids...)
	for i, id := range ids {
		setError(id, expected[i])
	}

	c.Check(err, gc.ErrorMatches, `3/3 items failed a bulk request: .*`)
}

func (*bulkSuite) TestErrorStringPartial(c *gc.C) {
	ids := []string{"a", "b", "c"}
	expected := []error{
		errors.Errorf("a"),
		nil,
		errors.Errorf("c"),
	}
	err, setError := errors.NewBulkError(ids...)
	for i, id := range ids {
		setError(id, expected[i])
	}

	c.Check(err, gc.ErrorMatches, `2/3 items failed a bulk request: .*`)
}

func (*bulkSuite) TestErrorStringNone(c *gc.C) {
	err, _ := errors.NewBulkError("a", "b", "c")

	c.Check(err, gc.ErrorMatches, `0/3 items failed a bulk request`)
}

func (*bulkSuite) TestErrorStringEmpty(c *gc.C) {
	err, _ := errors.NewBulkError()

	c.Check(err, gc.ErrorMatches, `0/0 items failed a bulk request`)
}

func (*bulkSuite) TestIDsOkay(c *gc.C) {
	expected := []string{"a", "b", "c"}
	err, _ := errors.NewBulkError(expected...)
	ids := err.IDs()

	c.Check(ids, jc.DeepEquals, expected)
}

func (*bulkSuite) TestIDsEmpty(c *gc.C) {
	err, _ := errors.NewBulkError()
	ids := err.IDs()

	c.Check(ids, gc.HasLen, 0)
}

func (*bulkSuite) TestEnumerateFull(c *gc.C) {
	ids := []string{"a", "b", "c"}
	expected := []error{
		errors.Errorf("a"),
		errors.Errorf("b"),
		errors.Errorf("c"),
	}
	err, setError := errors.NewBulkError(ids...)
	for i, id := range ids {
		setError(id, expected[i])
	}
	errors := err.Enumerate()

	c.Check(errors, jc.DeepEquals, expected)
}

func (*bulkSuite) TestEnumeratePartial(c *gc.C) {
	ids := []string{"a", "b", "c"}
	expected := []error{
		errors.Errorf("a"),
		nil,
		errors.Errorf("c"),
	}
	err, setError := errors.NewBulkError(ids...)
	for i, id := range ids {
		setError(id, expected[i])
	}
	errors := err.Enumerate()

	c.Check(errors, jc.DeepEquals, expected)
}

func (*bulkSuite) TestEnumerateNone(c *gc.C) {
	err, _ := errors.NewBulkError("a", "b", "c")
	errors := err.Enumerate()

	c.Check(errors, jc.DeepEquals, []error{nil, nil, nil})
}

func (*bulkSuite) TestEnumerateEmpty(c *gc.C) {
	err, _ := errors.NewBulkError()
	errors := err.Enumerate()

	c.Check(errors, gc.HasLen, 0)
}

func (*bulkSuite) TestIsBulkErrorDirect(c *gc.C) {
	err, _ := errors.NewBulkError("a", "b", "c")
	ok := errors.IsBulkError(err)

	c.Check(ok, jc.IsTrue)
}

func (*bulkSuite) TestIsBulkErrorIndirect(c *gc.C) {
	err, _ := errors.NewBulkError("a", "b", "c")
	wrapped := errors.Trace(err)
	ok := errors.IsBulkError(wrapped)

	c.Check(ok, jc.IsTrue)
}

func (*bulkSuite) TestIsBulkErrorFalse(c *gc.C) {
	err := errors.Errorf("not bulk")
	ok := errors.IsBulkError(err)

	c.Check(ok, jc.IsFalse)
}
