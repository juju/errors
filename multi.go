// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"strings"
)

// TODO(ericsnow) Make thread-safe?

// MultiError represents an ordered set of errors, aggregated into one.
// Each error may be associated with a string ID, which does not need
// to be unique.
type MultiError struct {
	errors []error
	ids    []string
}

// NewMultiError returns a new MultiError and a function that may be
// used to add errors to the set.

// It also returns a function that sets the error for any of the IDs.
// That function returns false if the provided ID is not recognized and
// true otherwise.
func NewMultiError() (*MultiError, func(error, string)) {
	multi := &MultiError{}
	return multi, multi.setError
}

func (multi *MultiError) setError(err error, id string) {
	multi.errors = append(multi.errors, err)
	multi.ids = append(multi.ids, id)
}

// Error returns the error string for the error.
func (multi MultiError) Error() string {
	byID, errors, ids := multi.collate()

	msg := fmt.Sprintf("%d errors", len(errors))
	if len(errors) == 0 {
		return msg
	}
	if len(byID) > 1 {
		// TODO(ericsnow) Don't count the empty string?
		msg += fmt.Sprintf(" (for %d IDs)", len(byID))
	}

	var errStrs []string
	for i, err := range errors {
		id := ids[i]
		var msg string
		if id != "" {
			msg = fmt.Sprintf("%v", err)
		} else {
			msg = fmt.Sprintf("(%q) %v", id, err)
		}
		errStrs = append(errStrs, msg)
	}
	msg += ": " + strings.Join(errStrs, ",")
	return msg
}

// Errors returns a new list containing all the errors, in the order
// they were added.
func (multi MultiError) Errors() ([]error, []string) {
	var errors []error
	var ids []string
	for i, err := range multi.errors {
		id := multi.ids[i]
		errors = append(errors, err)
		ids = append(ids, id)
	}
	return errors, ids
}

func (multi MultiError) collate() (map[string][]error, []error, []string) {
	collated := make(map[string][]error)
	errors, ids := multi.Errors()
	for i, err := range errors {
		id := ids[i]
		collated[id] = append(collated[id], err)
	}
	return collated, errors, ids
}

// IsMultiError reports whether err was created with NewMultiError().
func IsMultiError(err error) bool {
	err = Cause(err)
	_, ok := err.(*MultiError)
	return ok
}
