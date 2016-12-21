/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
