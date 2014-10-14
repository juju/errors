// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"reflect"
	"runtime"
)

// Location describes a source code location.
type Location struct {
	File string
	Line int
}

// String returns a location in filename.go:99 format.
func (loc Location) String() string {
	return fmt.Sprintf("%s:%d", loc.File, loc.Line)
}

// isSet reports whether the location has been set.
func (loc Location) isSet() bool {
	return loc.File != ""
}

// Err holds a description of an error along with information about
// where the error was created.
//
// It may be embedded in custom error types to add extra information that
// this errors package can understand.
type Err struct {
	// message holds an annotation of the error.
	message string

	// cause holds the cause of the error as returned
	// by the Cause method.
	cause error

	// Previous holds the Previous error in the error stack, if any.
	previous error

	// Location holds the source code location where the error was
	// created.
	location Location
}

// NewErr is used to return a Err for the purpose of embedding in other
// structures.  The location is not specified, and needs to be set with a call
// to SetLocation.
//
// For example:
// 		type FooError struct {
//			errors.Err
//			code int
//		}
//
//      func NewFooError(code int) error {
//			err := &FooError{errors.NewErr("foo"), code}
//			err.SetLocation(1)
//			return err
// 		}
func NewErr(format string, args ...interface{}) Err {
	return Err{
		message: fmt.Sprintf(format, args...),
	}
}

// Location is the file and line of where the error was most recently
// created or annotated.
func (e *Err) Location() Location {
	return e.location
}

// Previous returns the previous error in the error stack, if any.
func (e *Err) Previous() error {
	return e.previous
}

// The Cause of an error is the most recent error in the error stack that
// meets one of these criteria: the original error that was raised; the new
// error that was passed into the Wrap function; the most recently masked
// error; or nil if the error itself is considered the Cause.  Normally this
// method is not invoked directly, but instead through the Cause stand alone
// function.
func (e *Err) Cause() error {
	return e.cause
}

// Message returns the message stored with the most recent location. This is
// the empty string if the most recent call was Trace, or the message stored
// with Annotate or Mask.
func (e *Err) Message() string {
	return e.message
}

// Error implements error.Error.
func (e *Err) Error() string {
	// We want to walk up the stack of errors showing the annotations
	// as long as the cause is the same.
	err := e.previous
	if !sameError(Cause(err), e.cause) && e.cause != nil {
		err = e.cause
	}
	switch {
	case err == nil:
		return e.message
	case e.message == "":
		return err.Error()
	}
	return fmt.Sprintf("%s: %v", e.message, err)
}

// SetLocation records the source location of the error by at callDepth stack
// frames above the call.
func (e *Err) SetLocation(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.location = Location{trimGoPath(file), line}
}

// Ideally we'd have a way to check identity, but deep equals will do.
func sameError(e1, e2 error) bool {
	return reflect.DeepEqual(e1, e2)
}
