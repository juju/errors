// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
)

// a ConstError is a prototype for a certain type of error
type ConstError string

// ConstError implements error
func (e ConstError) Error() string {
	return ""
}

// Different types of errors
const (
	// Timeout represents an error on timeout.
	Timeout = ConstError("timeout")
	// NotFound represents an error when something has not been found.
	NotFound = ConstError("not found")
	// UserNotFound represents an error when a non-existent user is looked up.
	UserNotFound = ConstError("user not found")
	// Unauthorized represents an error when an operation is unauthorized.
	Unauthorized = ConstError("unauthorized")
	// NotImplemented represents an error when something is not
	// implemented.
	NotImplemented = ConstError("not implemented")
	// AlreadyExists represents and error when something already exists.
	AlreadyExists = ConstError("already exists")
	// NotSupported represents an error when something is not supported.
	NotSupported = ConstError("not supported")
	// NotValid represents an error when something is not valid.
	NotValid = ConstError("not valid")
	// NotProvisioned represents an error when something is not yet provisioned.
	NotProvisioned = ConstError("not provisioned")
	// NotAssigned represents an error when something is not yet assigned to
	// something else.
	NotAssigned = ConstError("not assigned")
	// BadRequest represents an error when a request has bad parameters.
	BadRequest = ConstError("bad request")
	// MethodNotAllowed represents an error when an HTTP request
	// is made with an inappropriate method.
	MethodNotAllowed = ConstError("method not allowed")
	// Forbidden represents an error when a request cannot be completed because of
	// missing privileges.
	Forbidden = ConstError("forbidden")
	// QuotaLimitExceeded is emitted when an action failed due to a quota limit check.
	QuotaLimitExceeded = ConstError("quota limit exceeded")
	// NotYetAvailable is the error returned when a resource is not yet available
	// but it might be in the future.
	NotYetAvailable = ConstError("not yet available")
)

// errWithType is an Err bundled with its error type (a ConstError)
type errWithType struct {
	err     Err
	errType ConstError
}

// errWithType implements error
func (e *errWithType) Error() string {
	return e.err.Error()
}

// e.Is compares `target` with e's error type
func (e *errWithType) Is(target error) bool {
	if &e.errType == nil {
		return false
	} else {
		return target == e.errType
	}
}

// Unwrapping an errWithType gives the underlying Err
func (e *errWithType) Unwrap() error {
	return &e.err
}

// wrap is a helper to construct an *wrapper.
func wrap(err error, format, suffix string, args ...interface{}) Err {
	newErr := Err{
		message:  fmt.Sprintf(format+suffix, args...),
		previous: err,
	}
	newErr.SetLocation(2)
	return newErr
}

// Timeoutf returns an error which satisfies IsTimeout().
func Timeoutf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " timeout", args...),
		Timeout,
	}
}

// NewTimeout returns an error which wraps err that satisfies
// IsTimeout().
func NewTimeout(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		Timeout,
	}
}

// IsTimeout reports whether err was created with Timeoutf() or
// NewTimeout().
func IsTimeout(err error) bool {
	return Is(err, Timeout)
}

// NotFoundf returns an error which satisfies IsNotFound().
func NotFoundf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not found", args...),
		NotFound,
	}
}

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotFound,
	}
}

// IsNotFound reports whether err was created with NotFoundf() or
// NewNotFound().
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}

// UserNotFoundf returns an error which satisfies IsUserNotFound().
func UserNotFoundf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " user not found", args...),
		UserNotFound,
	}
}

// NewUserNotFound returns an error which wraps err and satisfies
// IsUserNotFound().
func NewUserNotFound(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		UserNotFound,
	}
}

// IsUserNotFound reports whether err was created with UserNotFoundf() or
// NewUserNotFound().
func IsUserNotFound(err error) bool {
	return Is(err, UserNotFound)
}

// Unauthorizedf returns an error which satisfies IsUnauthorized().
func Unauthorizedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		Unauthorized,
	}
}

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		Unauthorized,
	}
}

// IsUnauthorized reports whether err was created with Unauthorizedf() or
// NewUnauthorized().
func IsUnauthorized(err error) bool {
	return Is(err, Unauthorized)
}

// NotImplementedf returns an error which satisfies IsNotImplemented().
func NotImplementedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not implemented", args...),
		NotImplemented,
	}
}

// NewNotImplemented returns an error which wraps err and satisfies
// IsNotImplemented().
func NewNotImplemented(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotImplemented,
	}
}

// IsNotImplemented reports whether err was created with
// NotImplementedf() or NewNotImplemented().
func IsNotImplemented(err error) bool {
	return Is(err, NotImplemented)
}

// AlreadyExistsf returns an error which satisfies IsAlreadyExists().
func AlreadyExistsf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " already exists", args...),
		AlreadyExists,
	}
}

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		AlreadyExists,
	}
}

// IsAlreadyExists reports whether the error was created with
// AlreadyExistsf() or NewAlreadyExists().
func IsAlreadyExists(err error) bool {
	return Is(err, AlreadyExists)
}

// NotSupportedf returns an error which satisfies IsNotSupported().
func NotSupportedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not supported", args...),
		NotSupported,
	}
}

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotSupported,
	}
}

// IsNotSupported reports whether the error was created with
// NotSupportedf() or NewNotSupported().
func IsNotSupported(err error) bool {
	return Is(err, NotSupported)
}

// NotValidf returns an error which satisfies IsNotValid().
func NotValidf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not valid", args...),
		NotValid,
	}
}

// NewNotValid returns an error which wraps err and satisfies IsNotValid().
func NewNotValid(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotValid,
	}
}

// IsNotValid reports whether the error was created with NotValidf() or
// NewNotValid().
func IsNotValid(err error) bool {
	return Is(err, NotValid)
}

// NotProvisionedf returns an error which satisfies IsNotProvisioned().
func NotProvisionedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not provisioned", args...),
		NotProvisioned,
	}
}

// NewNotProvisioned returns an error which wraps err that satisfies
// IsNotProvisioned().
func NewNotProvisioned(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotProvisioned,
	}
}

// IsNotProvisioned reports whether err was created with NotProvisionedf() or
// NewNotProvisioned().
func IsNotProvisioned(err error) bool {
	return Is(err, NotProvisioned)
}

// NotAssignedf returns an error which satisfies IsNotAssigned().
func NotAssignedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, " not assigned", args...),
		NotAssigned,
	}
}

// NewNotAssigned returns an error which wraps err that satisfies
// IsNotAssigned().
func NewNotAssigned(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotAssigned,
	}
}

// IsNotAssigned reports whether err was created with NotAssignedf() or
// NewNotAssigned().
func IsNotAssigned(err error) bool {
	return Is(err, NotAssigned)
}

// BadRequestf returns an error which satisfies IsBadRequest().
func BadRequestf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		BadRequest,
	}
}

// NewBadRequest returns an error which wraps err that satisfies
// IsBadRequest().
func NewBadRequest(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		BadRequest,
	}
}

// IsBadRequest reports whether err was created with BadRequestf() or
// NewBadRequest().
func IsBadRequest(err error) bool {
	return Is(err, BadRequest)
}

// MethodNotAllowedf returns an error which satisfies IsMethodNotAllowed().
func MethodNotAllowedf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		MethodNotAllowed,
	}
}

// NewMethodNotAllowed returns an error which wraps err that satisfies
// IsMethodNotAllowed().
func NewMethodNotAllowed(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		MethodNotAllowed,
	}
}

// IsMethodNotAllowed reports whether err was created with MethodNotAllowedf() or
// NewMethodNotAllowed().
func IsMethodNotAllowed(err error) bool {
	return Is(err, MethodNotAllowed)
}

// Forbiddenf returns an error which satistifes IsForbidden()
func Forbiddenf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		Forbidden,
	}
}

// NewForbidden returns an error which wraps err that satisfies
// IsForbidden().
func NewForbidden(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		Forbidden,
	}
}

// IsForbidden reports whether err was created with Forbiddenf() or
// NewForbidden().
func IsForbidden(err error) bool {
	return Is(err, Forbidden)
}

// QuotaLimitExceededf returns an error which satisfies IsQuotaLimitExceeded.
func QuotaLimitExceededf(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		QuotaLimitExceeded,
	}
}

// NewQuotaLimitExceeded returns an error which wraps err and satisfies
// IsQuotaLimitExceeded.
func NewQuotaLimitExceeded(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		QuotaLimitExceeded,
	}
}

// IsQuotaLimitExceeded returns true if the given error represents a
// QuotaLimitExceeded error.
func IsQuotaLimitExceeded(err error) bool {
	return Is(err, QuotaLimitExceeded)
}

// NotYetAvailablef returns an error which satisfies IsNotYetAvailable.
func NotYetAvailablef(format string, args ...interface{}) error {
	return &errWithType{
		wrap(nil, format, "", args...),
		NotYetAvailable,
	}
}

// NewNotYetAvailable returns an error which wraps err and satisfies
// IsNotYetAvailable.
func NewNotYetAvailable(err error, msg string) error {
	return &errWithType{
		wrap(err, msg, ""),
		NotYetAvailable,
	}
}

// IsNotYetAvailable reports err was created with NotYetAvailablef or
// NewNotYetAvailable.
func IsNotYetAvailable(err error) bool {
	return Is(err, NotYetAvailable)
}
