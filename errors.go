// Copyright 2014 Canonical Ltd.
// Licensed under the GPLv3, see LICENCE file for details.

package errors

import (
	"fmt"

	"github.com/juju/errgo"
	"github.com/juju/loggo"
)

// A juju Err is an errgo.Err, but formatted with the Cause.
type Err struct {
	errgo.Err
}

// Error implements error.Error.
func (e *Err) Error() string {
	switch {
	case e.Message_ == "" && e.Cause_ == nil:
		return "<no error>"
	case e.Message_ == "":
		return e.Cause_.Error()
	case e.Cause_ == nil:
		return e.Message_
	}
	return fmt.Sprintf("%s: %v", e.Message_, e.Cause_)
}

// newer is implemented by error types that can add a context message
// while preserving their type.
type newer interface {
	new(msg string) error
}

// wrap is a helper to construct an *wrapper.
func wrap(err error, format, suffix string, args ...interface{}) Err {
	return Err{
		errgo.Err{
			Message_:    fmt.Sprintf(format+suffix, args...),
			Underlying_: err,
		},
	}
}

// Contextf prefixes any error stored in err with text formatted
// according to the format specifier. If err does not contain an
// error, Contextf does nothing. All errors created with functions
// from this package are preserved when wrapping.
func Contextf(err *error, format string, args ...interface{}) {
	if *err == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	errNewer, ok := (*err).(newer)
	if ok {
		*err = errNewer.new(msg)
		return
	}
	*err = fmt.Errorf("%s: %v", msg, *err)
}

// Maskf masks the given error (when it is not nil) with the given
// format string and arguments (like fmt.Sprintf), returning a new
// error. If *err is nil, Maskf does nothing.
func Maskf(err *error, format string, args ...interface{}) {
	if *err == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	newErr := &Err{
		errgo.Err{
			Message_:    fmt.Sprintf("%s: %v", msg, *err),
			Underlying_: *err,
		},
	}
	newErr.SetLocation(1)
	*err = newErr
}

// notFound represents an error when something has not been found.
type notFound struct {
	Err
}

func (e *notFound) new(msg string) error {
	return NewNotFound(e, msg)
}

// NotFoundf returns an error which satisfies IsNotFound().
func NotFoundf(format string, args ...interface{}) error {
	err := &notFound{wrap(nil, format, " not found", args...)}
	err.SetLocation(1)
	return err
}

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	newErr := &notFound{wrap(err, msg, "")}
	newErr.SetLocation(1)
	return newErr
}

// IsNotFound reports whether err was created with NotFoundf() or
// NewNotFound().
func IsNotFound(err error) bool {
	err = errgo.Cause(err)
	_, ok := err.(*notFound)
	return ok
}

// unauthorized represents an error when an operation is unauthorized.
type unauthorized struct {
	Err
}

func (e *unauthorized) new(msg string) error {
	return NewUnauthorized(e, msg)
}

// Unauthorizedf returns an error which satisfies IsUnauthorized().
func Unauthorizedf(format string, args ...interface{}) error {
	err := &unauthorized{wrap(nil, format, "", args...)}
	err.SetLocation(1)
	return err
}

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	newErr := &unauthorized{wrap(err, msg, "")}
	newErr.SetLocation(1)
	return newErr
}

// IsUnauthorized reports whether err was created with Unauthorizedf() or
// NewUnauthorized().
func IsUnauthorized(err error) bool {
	err = errgo.Cause(err)
	_, ok := err.(*unauthorized)
	return ok
}

// notImplemented represents an error when something is not
// implemented.
type notImplemented struct {
	Err
}

func (e *notImplemented) new(msg string) error {
	return NewNotImplemented(e, msg)
}

// NotImplementedf returns an error which satisfies IsNotImplemented().
func NotImplementedf(format string, args ...interface{}) error {
	err := &notImplemented{wrap(nil, format, " not implemented", args...)}
	err.SetLocation(1)
	return err
}

// NewNotImplemented returns an error which wraps err and satisfies
// IsNotImplemented().
func NewNotImplemented(err error, msg string) error {
	newErr := &notImplemented{wrap(err, msg, "")}
	newErr.SetLocation(1)
	return newErr
}

// IsNotImplemented reports whether err was created with
// NotImplementedf() or NewNotImplemented().
func IsNotImplemented(err error) bool {
	err = errgo.Cause(err)
	_, ok := err.(*notImplemented)
	return ok
}

// alreadyExists represents and error when something already exists.
type alreadyExists struct {
	Err
}

func (e *alreadyExists) new(msg string) error {
	return NewAlreadyExists(e, msg)
}

// AlreadyExistsf returns an error which satisfies IsAlreadyExists().
func AlreadyExistsf(format string, args ...interface{}) error {
	err := &alreadyExists{wrap(nil, format, " already exists", args...)}
	err.SetLocation(1)
	return err
}

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string) error {
	newErr := &alreadyExists{wrap(err, msg, "")}
	newErr.SetLocation(1)
	return newErr
}

// IsAlreadyExists reports whether the error was created with
// AlreadyExistsf() or NewAlreadyExists().
func IsAlreadyExists(err error) bool {
	err = errgo.Cause(err)
	_, ok := err.(*alreadyExists)
	return ok
}

// notSupported represents an error when something is not supported.
type notSupported struct {
	Err
}

func (e *notSupported) new(msg string) error {
	return NewNotSupported(e, msg)
}

// NotSupportedf returns an error which satisfies IsNotSupported().
func NotSupportedf(format string, args ...interface{}) error {
	err := &notSupported{wrap(nil, format, " not supported", args...)}
	err.SetLocation(1)
	return err
}

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	newErr := &notSupported{wrap(err, msg, "")}
	newErr.SetLocation(1)
	return newErr
}

// IsNotSupported reports whether the error was created with
// NotSupportedf() or NewNotSupported().
func IsNotSupported(err error) bool {
	err = errgo.Cause(err)
	_, ok := err.(*notSupported)
	return ok
}

// LoggedErrorf logs the error and return an error with the same text.
func LoggedErrorf(logger loggo.Logger, format string, a ...interface{}) error {
	logger.Logf(loggo.ERROR, format, a...)
	return fmt.Errorf(format, a...)
}
