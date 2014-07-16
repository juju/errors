// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"unsafe"
)

// interfaceInternals reflects the internal state of an interface.
type interfaceInternals struct {
	typ uintptr
	val uintptr
}

// sameError checks to see if the two errors are the same error. This works
// with uncomparable types by actually looking at the interface values
// internally. Somewhat a hack.
func sameError(e1, e2 error) bool {
	i1 := *(*interfaceInternals)(unsafe.Pointer(&e1))
	i2 := *(*interfaceInternals)(unsafe.Pointer(&e2))
	return i1 == i2
}
