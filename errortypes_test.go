// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	stderrors "errors"
	"fmt"

	"github.com/juju/errors"
	gc "gopkg.in/check.v1"
)

// errorInfo holds information about a single error type: its type
// and name, wrapping and variable arguments constructors and message
// suffix.
type errorInfo struct {
	errType         errors.ConstError
	errName         string
	argsConstructor func(string, ...interface{}) error
	wrapConstructor func(error, string) error
	suffix          string
}

// allErrors holds information for all defined errors. When adding new
// errors, add them here as well to include them in tests.
var allErrors = []*errorInfo{
	{errors.Timeout, "Timeout", errors.Timeoutf, errors.NewTimeout, " timeout"},
	{errors.NotFound, "NotFound", errors.NotFoundf, errors.NewNotFound, " not found"},
	{errors.UserNotFound, "UserNotFound", errors.UserNotFoundf, errors.NewUserNotFound, " user not found"},
	{errors.Unauthorized, "Unauthorized", errors.Unauthorizedf, errors.NewUnauthorized, ""},
	{errors.NotImplemented, "NotImplemented", errors.NotImplementedf, errors.NewNotImplemented, " not implemented"},
	{errors.AlreadyExists, "AlreadyExists", errors.AlreadyExistsf, errors.NewAlreadyExists, " already exists"},
	{errors.NotSupported, "NotSupported", errors.NotSupportedf, errors.NewNotSupported, " not supported"},
	{errors.NotValid, "NotValid", errors.NotValidf, errors.NewNotValid, " not valid"},
	{errors.NotProvisioned, "NotProvisioned", errors.NotProvisionedf, errors.NewNotProvisioned, " not provisioned"},
	{errors.NotAssigned, "NotAssigned", errors.NotAssignedf, errors.NewNotAssigned, " not assigned"},
	{errors.MethodNotAllowed, "MethodNotAllowed", errors.MethodNotAllowedf, errors.NewMethodNotAllowed, ""},
	{errors.BadRequest, "BadRequest", errors.BadRequestf, errors.NewBadRequest, ""},
	{errors.Forbidden, "Forbidden", errors.Forbiddenf, errors.NewForbidden, ""},
	{errors.QuotaLimitExceeded, "QuotaLimitExceeded", errors.QuotaLimitExceededf, errors.NewQuotaLimitExceeded, ""},
	{errors.NotYetAvailable, "NotYetAvailable", errors.NotYetAvailablef, errors.NewNotYetAvailable, ""},
}

type errorTypeSuite struct{}

var _ = gc.Suite(&errorTypeSuite{})

func (t *errorInfo) equal(t0 *errorInfo) bool {
	if t0 == nil {
		return false
	}
	return t == t0
}

type errorTest struct {
	err     error
	message string
	errInfo *errorInfo
}

func deferredAnnotatef(err error, format string, args ...interface{}) error {
	errors.DeferredAnnotatef(&err, format, args...)
	return err
}

func mustSatisfy(c *gc.C, err error, errInfo *errorInfo) {
	if errInfo != nil {
		msg := fmt.Sprintf("Is(err, %s) should be TRUE when err := %#v", errInfo.errName, err)
		c.Check(errors.Is(err, errInfo.errType), gc.Equals, true, gc.Commentf(msg))
	}
}

func mustNotSatisfy(c *gc.C, err error, errInfo *errorInfo) {
	if errInfo != nil {
		msg := fmt.Sprintf("Is(err, %s) should be FALSE when err := %#v", errInfo.errName, err)
		c.Check(errors.Is(err, errInfo.errType), gc.Equals, false, gc.Commentf(msg))
	}
}

func checkErrorMatches(c *gc.C, err error, message string, errInfo *errorInfo) {
	if message == "<nil>" {
		c.Check(err, gc.IsNil)
		c.Check(errInfo, gc.IsNil)
	} else {
		c.Check(err, gc.ErrorMatches, message)
	}
}

func runErrorTests(c *gc.C, errorTests []errorTest, checkMustSatisfy bool) {
	for i, t := range errorTests {
		c.Logf("test %d: %T: %v", i, t.err, t.err)
		checkErrorMatches(c, t.err, t.message, t.errInfo)
		if checkMustSatisfy {
			mustSatisfy(c, t.err, t.errInfo)
		}

		// Check all other satisfiers to make sure none match.
		for _, otherErrInfo := range allErrors {
			if checkMustSatisfy && otherErrInfo.equal(t.errInfo) {
				continue
			}
			mustNotSatisfy(c, t.err, otherErrInfo)
		}
	}
}

func (*errorTypeSuite) TestDeferredAnnotatef(c *gc.C) {
	// Ensure DeferredAnnotatef annotates the errors.
	errorTests := []errorTest{}
	for _, errInfo := range allErrors {
		errorTests = append(errorTests, []errorTest{{
			deferredAnnotatef(nil, "comment"),
			"<nil>",
			nil,
		}, {
			deferredAnnotatef(stderrors.New("blast"), "comment"),
			"comment: blast",
			nil,
		}, {
			deferredAnnotatef(errInfo.argsConstructor("foo %d", 42), "comment %d", 69),
			"comment 69: foo 42" + errInfo.suffix,
			errInfo,
		}, {
			deferredAnnotatef(errInfo.argsConstructor(""), "comment"),
			"comment: " + errInfo.suffix,
			errInfo,
		}, {
			deferredAnnotatef(errInfo.wrapConstructor(stderrors.New("pow!"), "woo"), "comment"),
			"comment: woo: pow!",
			errInfo,
		}}...)
	}

	runErrorTests(c, errorTests, true)
}

func (*errorTypeSuite) TestAllErrors(c *gc.C) {
	errorTests := []errorTest{}
	for _, errInfo := range allErrors {
		errorTests = append(errorTests, []errorTest{{
			nil,
			"<nil>",
			nil,
		}, {
			errInfo.argsConstructor("foo %d", 42),
			"foo 42" + errInfo.suffix,
			errInfo,
		}, {
			errInfo.argsConstructor(""),
			errInfo.suffix,
			errInfo,
		}, {
			errInfo.wrapConstructor(stderrors.New("pow!"), "prefix"),
			"prefix: pow!",
			errInfo,
		}, {
			errInfo.wrapConstructor(stderrors.New("pow!"), ""),
			"pow!",
			errInfo,
		}, {
			errInfo.wrapConstructor(nil, "prefix"),
			"prefix",
			errInfo,
		}}...)
	}

	runErrorTests(c, errorTests, true)
}

// TestThatYouAlwaysGetError is a regression test for checking that the wrap
// constructor for our error types always returns a valid error object even if
// don't feed the construct with an instantiated error or a non empty string.
func (*errorTypeSuite) TestThatYouAlwaysGetError(c *gc.C) {
	for _, errType := range allErrors {
		err := errType.wrapConstructor(nil, "")
		c.Assert(err.Error(), gc.Equals, "")
	}
}

func (*errorTypeSuite) TestWithTypeNil(c *gc.C) {
	myErr := errors.ConstError("do you feel lucky?")
	c.Assert(errors.WithType(nil, myErr), gc.IsNil)
}

func (*errorTypeSuite) TestWithType(c *gc.C) {
	myErr := errors.ConstError("do you feel lucky?")
	myErr2 := errors.ConstError("i don't feel lucky")
	err := errors.New("yes")

	err = errors.WithType(err, myErr)
	c.Assert(errors.Is(err, myErr), gc.Equals, true)
	c.Assert(err.Error(), gc.Equals, "yes")
	c.Assert(errors.Is(err, myErr2), gc.Equals, false)
}
