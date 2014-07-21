// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/juju/errgo"
)

// NOTE: the SetLocation calls explicitly call into the embedded errgo Err
// structure as gccgo creates an extra entry in the runtime stack for the
// generated method on the outer Err structure that just defers to the
// errgo.Err method.  In order to get the package passing with gccgo this is
// needed.

// New is a drop in replacement for the standard libary errors module that records
// the location that the error is created.
//
// For example:
//    return errors.New("validation failed")
//
func New(message string) error {
	err := &Err{errgo.Err{Message_: message}}
	err.Err.SetLocation(1)
	return err
}

// Errorf creates a new annotated error and records the location that the
// error is created.  This should be a drop in replacement for fmt.Errorf.
//
// For example:
//    return errors.Errorf("validation failed: %s", message)
//
func Errorf(format string, args ...interface{}) error {
	err := &Err{errgo.Err{Message_: fmt.Sprintf(format, args...)}}
	err.Err.SetLocation(1)
	return err
}

// Trace always returns an annotated error.  Trace records the
// location of the Trace call, and adds it to the annotation stack.
// If the error argument is nil, the result of Trace is nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Trace(err)
//   }
//
func Trace(other error) error {
	if other == nil {
		return nil
	}
	err := &Err{errgo.Err{Underlying_: other, Cause_: Cause(other)}}
	err.Err.SetLocation(1)
	return err
}

// Annotate is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
// If the error argument is nil, the result of Annotate is nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Annotate(err, "failed to frombulate")
//   }
//
func Annotate(other error, message string) error {
	if other == nil {
		return nil
	}
	// Underlying is the previous link used for traversing the stack.
	// Cause is the reason for this error.
	err := &Err{
		errgo.Err{
			Underlying_: other,
			Cause_:      Cause(other),
			Message_:    message,
		},
	}
	err.Err.SetLocation(1)
	return err
}

// Annotatef is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
// If the error argument is nil, the result of Annotatef is nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Annotatef(err, "failed to frombulate the %s", arg)
//   }
//
func Annotatef(other error, format string, args ...interface{}) error {
	if other == nil {
		return nil
	}
	// Underlying is the previous link used for traversing the stack.
	// Cause is the reason for this error.
	err := &Err{
		errgo.Err{
			Underlying_: other,
			Cause_:      Cause(other),
			Message_:    fmt.Sprintf(format, args...),
		},
	}
	err.Err.SetLocation(1)
	return err
}

// Wrap changes the error value that is returned with LastError. The location
// of the Wrap call is also stored in the annotation stack.
// If the error argument is nil, the result of Wrap is nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       newErr := &packageError{"more context", private_value}
//       return errors.Wrap(err, newErr)
//   }
//
func Wrap(other, newDescriptive error) error {
	if other == nil {
		return nil
	}
	err := &Err{
		errgo.Err{
			Underlying_: other,
			Cause_:      newDescriptive,
		},
	}
	err.Err.SetLocation(1)
	return err
}

// Check looks at the Cause of the error to see if it matches the checker
// function.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       if errors.Check(err, os.IsNotExist) {
//           return someOtherFunc()
//       }
//   }
//
func Check(err error, checker func(error) bool) bool {
	return checker(errgo.Cause(err))
}

// ErrorStack returns a string representation of the annotated error. If the
// error passed as the parameter is not an annotated error, the result is
// simply the result of the Error() method on that error.
//
// If the error is an annotated error, a multi-line string is returned where
// each line represents one entry in the annotation stack. The full filename
// from the call stack is used in the output.
func ErrorStack(err error) string {
	if err == nil {
		return ""
	}
	// We want the first error first
	var lines []string
	for {
		var buff []byte
		if err, ok := err.(errgo.Locationer); ok {
			loc := err.Location()
			// Strip off the leading GOPATH/src path elements.
			loc.File = trimGoPath(loc.File)
			if loc.IsSet() {
				buff = append(buff, loc.String()...)
				buff = append(buff, ": "...)
			}
		}
		if cerr, ok := err.(errgo.Wrapper); ok {
			message := cerr.Message()
			buff = append(buff, message...)
			// If there is a cause for this error, and it is different to the cause
			// of the underlying error, then output the error string in the stack trace.
			var cause error
			if err1, ok := err.(errgo.Causer); ok {
				cause = err1.Cause()
			}
			err = cerr.Underlying()
			if cause != nil && !sameError(Cause(err), cause) {
				if message != "" {
					buff = append(buff, ": "...)
				}
				buff = append(buff, cause.Error()...)
			}
		} else {
			buff = append(buff, err.Error()...)
			err = nil
		}
		lines = append(lines, string(buff))
		if err == nil {
			break
		}
	}
	// reverse the lines to get the original error, which was at the end of
	// the list, back to the start.
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}
	return strings.Join(result, "\n")
}

// Ideally we'd have a way to check identity, but deep equals will do.
func sameError(e1, e2 error) bool {
	return reflect.DeepEqual(e1, e2)
}
