// Simple checks that don't don't make any assertions but are most likely of
// use in tests.
package check

import "reflect"

// Checks to see a variable is nil in the colloquial sense. This returns true
// true if the parameter is == nil or if the parameter has a fixed dynamic type
// and a nil value.
func IsNil(got interface{}) bool {
	if got == nil {
		return true
	}

	v := reflect.ValueOf(got)
	kind := v.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && v.IsNil() {
		return true
	}

	return false
}
