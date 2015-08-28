// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"strings"
)

// BulkErrors represents when a bulk request fails for one or more of
// the items.
type BulkErrors struct {
	ids    []string
	errors map[string]error
	count  int
}

// NewBulkErrors returns a new BulkErrors primed for the provided IDs.
// It also returns a function that sets the error for one of the IDs.
// That function returns false if the provided ID is not recognized and
// true otherwise.
func NewBulkErrors(ids ...string) (*BulkErrors, func(string, error) bool) {
	be := &BulkErrors{
		ids:    ids,
		errors: make(map[string]error, len(ids)),
	}
	for _, id := range ids {
		be.errors[id] = nil
	}
	return be, be.setError
}

func (be *BulkErrors) setError(id string, err error) bool {
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

// Error returns the error string for the error.
func (be BulkErrors) Error() string {
	msg := fmt.Sprintf("%d/%d items failed a bulk request", be.count, len(be.ids))
	if be.count == 0 {
		return msg
	}

	var errors []string
	for _, id := range be.ids {
		if err := be.errors[id]; err != nil {
			errors = append(errors, fmt.Sprintf("(%q) %v", id, err))
		}
	}
	msg += ": " + strings.Join(errors, ",")
	return msg
}

// IDs returns a new list containing the IDs in the originally provided order.
func (be BulkErrors) IDs() []string {
	ids := make([]string, len(be.ids))
	copy(ids, be.ids)
	return ids
}

// Enumerate returns the list of errors (or nils) corresponding to the
// original IDs.
func (be BulkErrors) Enumerate() []error {
	errors := make([]error, len(be.ids))
	for i, id := range be.ids {
		errors[i] = be.errors[id]
	}
	return errors
}

// IsBulkErrors reports whether err was created with NewBulkErrors().
func IsBulkErrors(err error) bool {
	err = Cause(err)
	_, ok := err.(*BulkErrors)
	return ok
}
