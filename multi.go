// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"strings"
)

type multiErrorItem struct {
	err error
	id  string
}

// String returns a string representation of the item.
func (mei multiErrorItem) String() string {
	if mei.id == "" {
		return fmt.Sprintf("%v", mei.err)
	} else {
		return fmt.Sprintf("(%q) %v", mei.id, mei.err)
	}
}

// MultiError represents an ordered set of errors, aggregated into one.
// Each error may be associated with a string ID, which does not need
// to be unique.
type MultiError struct {
	errors []multiErrorItem
	byID   map[string][]error
}

// NewMultiError returns a new MultiError and a function that may be
// used to add errors to the set.

// It also returns a function that sets the error for any of the IDs.
// That function returns false if the provided ID is not recognized and
// true otherwise.
func NewMultiError() (*MultiError, func(error, string)) {
	multi := &MultiError{
		byID: make(map[string][]error),
	}
	return multi, multi.setError
}

func (multi *MultiError) setError(err error, id string) {
	multi.byID[id] = append(multi.byID[id], err)
	multi.errors = append(multi.errors, multiErrorItem{err, id})
}

// Error returns the error string for the error.
func (multi MultiError) Error() string {
	msg := fmt.Sprintf("%d errors", len(multi.errors))
	if len(multi.errors) == 0 {
		return msg
	}
	if len(multi.byID) > 1 {
		// TODO(ericsnow) Don't count the empty string?
		msg += fmt.Sprintf(" (for %d IDs)", len(multi.byID))
	}

	var errors []string
	for _, err := range multi.errors {
		errors = append(errors, err.String())
	}
	msg += ": " + strings.Join(errors, ",")
	return msg
}

// Errors returns a new list containing all the errors, in the order
// they were added.
func (multi MultiError) Errors() ([]error, []string) {
	var errors []error
	var ids []string
	for _, item := range multi.errors {
		errors = append(errors, item.err)
		ids = append(ids, item.id)
	}
	return errors, ids
}

// IsMultiError reports whether err was created with NewMultiError().
func IsMultiError(err error) bool {
	err = Cause(err)
	_, ok := err.(*MultiError)
	return ok
}
