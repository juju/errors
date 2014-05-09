package errors

import (
	"fmt"

	"github.com/juju/errgo"
)

// Errorf creates a new annotated error and records the location that the
// error is created.  This should be a drop in replacement for fmt.Errorf.
//
// For example:
//    return errors.Errorf("validation failed")
//
func Errorf(format string, args ...interface{}) error {
	err := &errgo.Err{Message_: fmt.Sprintf(format, args...)}
	err.SetLocation(1)
	return err
}

// Trace always returns an annotated error.  Trace records the
// location of the Trace call, and adds it to the annotation stack.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Trace(err)
//   }
//
func Trace(other error) error {
	err := &errgo.Err{Underlying_: other}
	err.SetLocation(1)
	return err
}

// Annotate is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Annotate(err, "failed to frombulate")
//   }
//
func Annotate(other error, message string) error {
	// Underlying is the previous link used for traversing the stack.
	// Cause is the reason for this error.
	err := &errgo.Err{
		Underlying_: other,
		Cause_:      other,
		Message_:    message,
	}
	err.SetLocation(1)
	return err
}

// Annotatef is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Annotatef(err, "failed to frombulate the %s", arg)
//   }
//
func Annotatef(other error, format string, args ...interface{}) error {
	// Underlying is the previous link used for traversing the stack.
	// Cause is the reason for this error.
	err := &errgo.Err{
		Underlying_: other,
		Cause_:      other,
		Message_:    fmt.Sprintf(format, args...),
	}
	err.SetLocation(1)
	return err
}

// Wrap changes the error value that is returned with LastError. The location
// of the Wrap call is also stored in the annotation stack.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       newErr := &packageError{"more context", private_value}
//       return errors.Wrap(err, newErr)
//   }
//
func Wrap(other, newDescriptive error) error {
	err := &errgo.Err{
		Underlying_: other,
		Cause_:      newDescriptive,
	}
	err.SetLocation(1)
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
	var buff []byte
	for {
		if err, ok := err.(errgo.Locationer); ok {
			loc := err.Location()
			// trimGoPath(&loc)
			if loc.IsSet() {
				buff = append(buff, loc.String()...)
				buff = append(buff, ": "...)
			}
		}
		if cerr, ok := err.(errgo.Wrapper); ok {
			buff = append(buff, cerr.Message()...)
			err = cerr.Underlying()
		} else {
			buff = append(buff, err.Error()...)
			err = nil
		}
		if err == nil {
			break
		}
		buff = append(buff, '\n')
	}
	return string(buff)
}
