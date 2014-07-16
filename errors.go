// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"

	"github.com/juju/errgo"
	"github.com/juju/loggo"
)

// NOTE: the SetLocation calls explicitly call into the embedded errgo Err
// structure as gccgo creates an extra entry in the runtime stack for the
// generated method on the outer Err structure that just defers to the
// errgo.Err method.  In order to get the package passing with gccgo this is
// needed.

// A juju Err is an errgo.Err but the error string generated walks up the
// stack of errors adding the messages but stops if the cause of the
// underlying error is different.
type Err struct {
	errgo.Err
}

// Error implements error.Error.
func (e *Err) Error() string {
	// We want to walk up the stack of errors showing the annotations
	// as long as the cause is the same.
	err := e.Underlying_
	if !sameError(Cause(e.Underlying_), e.Cause_) && e.Cause_ != nil {
		err = e.Cause_
	}
	switch {
	case err == nil:
		return e.Message_
	case e.Message_ == "":
		return err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Message_, err)
}

// Cause returns the cause of the given error.  If err does not implement
// errgo.Causer or its Cause method returns nil, it returns err itself.
func Cause(err error) error {
	return errgo.Cause(err)
}

// newer is implemented by error types that can add a context message
// while preserving their type.
type newer interface {
	new(msg string) error
}

// wrap is a helper to construct an *wrapper.
func wrap(err error, format, suffix string, args ...interface{}) Err {
	newErr := Err{
		errgo.Err{
			Message_:    fmt.Sprintf(format+suffix, args...),
			Underlying_: err,
		},
	}
	newErr.Err.SetLocation(2)
	return newErr
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
	newErr := &Err{
		errgo.Err{
			Message_:    msg,
			Underlying_: *err,
			Cause_:      Cause(*err),
		},
	}
	newErr.Err.SetLocation(1)
	*err = newErr
}

// Maskf masks the given error (when it is not nil) with the given
// format string and arguments (like fmt.Sprintf), returning a new
// error. If *err is nil, Maskf does nothing.
func Maskf(err *error, format string, args ...interface{}) {
	if *err == nil {
		return
	}
	newErr := &Err{
		errgo.Err{
			Message_:    fmt.Sprintf(format, args...),
			Underlying_: *err,
		},
	}
	newErr.Err.SetLocation(1)
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
	return &notFound{wrap(nil, format, " not found", args...)}
}

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	return &notFound{wrap(err, msg, "")}
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
	return &unauthorized{wrap(nil, format, "", args...)}
}

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	return &unauthorized{wrap(err, msg, "")}
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
	return &notImplemented{wrap(nil, format, " not implemented", args...)}
}

// NewNotImplemented returns an error which wraps err and satisfies
// IsNotImplemented().
func NewNotImplemented(err error, msg string) error {
	return &notImplemented{wrap(err, msg, "")}
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
	return &alreadyExists{wrap(nil, format, " already exists", args...)}
}

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string) error {
	return &alreadyExists{wrap(err, msg, "")}
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
	return &notSupported{wrap(nil, format, " not supported", args...)}
}

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	return &notSupported{wrap(err, msg, "")}
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
