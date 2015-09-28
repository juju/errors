// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

// Since variables are declared before the init block, in order to get the goPath
// we need to return it rather than just reference it.
func GoPath() string {
	return goPath
}

var TrimGoPath = trimGoPath

func ExposeMultiError(err error) ([]error, []string) {
	multi := err.(*MultiError)
	var errors []error
	var ids []string
	for _, item := range multi.errors {
		errors = append(errors, item.err)
		ids = append(ids, item.id)
	}
	return errors, ids
}
