// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
)

// BulkError represents the case where multiple items are handled and at
// least one of them failed. The [ordered] set of items, each identified
// by a unique string, is intrinsic to the BulkError. Consequently,
// those IDs must be provided when the BulkError is created and will not
// change.
//
// An error for any given ID may be set, reset, or unset using the
// function returned from NewBulkError. All IDs and any associated
// errors are accessible via BulkError methods.
//
// BulkError is relevant for several use cases, including bulk API
// requests.
type BulkError struct {
	ids    []string
	errors map[string]error
	count  int
}

// NewBulkError returns a new BulkError primed for the provided IDs.
// It also returns a function that sets the error for any of the IDs.
// That function returns false if the provided ID is not recognized and
// true otherwise.
func NewBulkError(ids ...string) (*BulkError, func(string, error) bool) {
	be := &BulkError{
		ids:    ids,
		errors: make(map[string]error, len(ids)),
	}
	for _, id := range ids {
		be.errors[id] = nil
	}
	return be, be.setError
}

func (be *BulkError) setError(id string, err error) bool {
	existing, ok := be.errors[id]
	if !ok {
		return false
	}
	if existing == nil && err != nil {
		be.count += 1
	}
	if existing != nil && err == nil {
		be.count -= 1
	}
	be.errors[id] = err
	return true
}

// TODO(ericsnow) Follow the precedent of ErrorResults.Combine()?

// Error returns the error string for the error.
func (be BulkError) Error() string {
	return fmt.Sprintf("%d/%d items failed a bulk request", be.count, len(be.ids))
}

// NoErrors determines whether or not the BulkError has any item errors set.
func (be BulkError) NoErrors() bool {
	return be.count == 0
}

// TODO(ericsnow) Add a OneError() method a la ErrorResults?

// IDs returns a new list containing the IDs in the originally provided order.
func (be BulkError) IDs() []string {
	ids := make([]string, len(be.ids))
	copy(ids, be.ids)
	return ids
}

// Enumerate returns the list of errors (or nils) corresponding to the
// original IDs.
func (be BulkError) Enumerate() []error {
	errors := make([]error, len(be.ids))
	for i, id := range be.ids {
		errors[i] = be.errors[id]
	}
	return errors
}

// IsBulkError reports whether err was created with NewBulkError().
func IsBulkError(err error) bool {
	err = Cause(err)
	_, ok := err.(*BulkError)
	return ok
}
