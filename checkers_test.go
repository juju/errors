// Copyright 2022 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors_test

import (
	"fmt"
	"reflect"
	"strings"

	gc "gopkg.in/check.v1"
)

// containsChecker is a copy of the containsChecker from juju/testing
type containsChecker struct {
	*gc.CheckerInfo
}

// satisfiesChecker is a copy of the satisfiesChecker from juju/testing
type satisfiesChecker struct {
	*gc.CheckerInfo
}

// Contains is a copy of the Contains checker from juju/testing
var Contains gc.Checker = &containsChecker{
	&gc.CheckerInfo{Name: "Contains", Params: []string{"obtained", "expected"}},
}

// Satisfies is a copy of the Satisfies checker from juju/testing
var Satisfies gc.Checker = &satisfiesChecker{
	&gc.CheckerInfo{
		Name:   "Satisfies",
		Params: []string{"obtained", "func(T) bool"},
	},
}

// canBeNil is copied from juju/testing
func canBeNil(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return true
	}
	return false
}

// Check is copied from juju/testing containsChecker
func (checker *containsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	expected, ok := params[1].(string)
	if !ok {
		return false, "expected must be a string"
	}

	obtained, isString := stringOrStringer(params[0])
	if isString {
		return strings.Contains(obtained, expected), ""
	}

	return false, "Obtained value is not a string and has no .String()"
}

// Check is copied from juju/testing satisfiesChecker
func (checker *satisfiesChecker) Check(params []interface{}, names []string) (result bool, error string) {
	f := reflect.ValueOf(params[1])
	ft := f.Type()
	if ft.Kind() != reflect.Func ||
		ft.NumIn() != 1 ||
		ft.NumOut() != 1 ||
		ft.Out(0) != reflect.TypeOf(true) {
		return false, fmt.Sprintf("expected func(T) bool, got %s", ft)
	}
	v := reflect.ValueOf(params[0])
	if !v.IsValid() {
		if !canBeNil(ft.In(0)) {
			return false, fmt.Sprintf("cannot assign nil to argument %T", ft.In(0))
		}
		v = reflect.Zero(ft.In(0))
	}
	if !v.Type().AssignableTo(ft.In(0)) {
		return false, fmt.Sprintf("wrong argument type %s for %s", v.Type(), ft)
	}
	return f.Call([]reflect.Value{v})[0].Interface().(bool), ""
}

// stringOrStringer is copied from juju/testing
func stringOrStringer(value interface{}) (string, bool) {
	result, isString := value.(string)
	if !isString {
		if stringer, isStringer := value.(fmt.Stringer); isStringer {
			result, isString = stringer.String(), true
		}
	}
	return result, isString
}
