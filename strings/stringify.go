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

// String conversions that are most likely of use in tests.
package strings

import (
	"fmt"
	"reflect"
	"strconv"
)

// Generically converts an interface into string. Useful for test
// error messages.  Nils are encoded as "<nil>". Objects conforming to
// fmt.Stringer are encoded using the result of String(). Non-string
// types are encoded via fmt.Sprintf's %+v format specifier. Strings
// are returned as backticked strings if possible, or else they are
// returned as double-quoted string literals with appropriate escape
// sequences (see strconv.Quote).
func Stringify(i interface{}) string {
	var s string
	switch t := i.(type) {
	case string:
		s = t
	case *string:
		if t == nil {
			return "<nil>"
		}
		s = *t
	case fmt.Stringer:
		if reflect.ValueOf(t).IsNil() {
			return "<nil>"
		}

		s = t.String()
	default:
		return fmt.Sprintf("%+v", i)
	}

	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	} else {
		return strconv.Quote(s)
	}
}
