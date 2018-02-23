/*
Copyright 2018 Turbine Labs, Inc.

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
package check

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// During deepEqual, we keep track of checks that are in progress and
// prevent circular checks. This type is used as a map key.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func trackedType(k reflect.Kind) bool {
	return k == reflect.Map || k == reflect.Slice || k == reflect.Ptr || k == reflect.Interface
}

func boolResult(got, want bool, path, reason string) (bool, string) {
	if got == want {
		return true, ""
	}

	if got {
		return false, fmt.Sprintf("got%s is %s, want%s is not %s", path, reason, path, reason)
	}

	return false, fmt.Sprintf("got%s is not %s, want%s is %s", path, reason, path, reason)
}

func render(t reflect.Type, v reflect.Value) string {
	switch t.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", v.String())

	case reflect.Ptr:
		if t.Elem().Kind() == reflect.String {
			if v.IsNil() {
				return "<nil>"
			}
			return fmt.Sprintf("%q", reflect.Indirect(v).String())
		}
		fallthrough
	default:
		return fmt.Sprintf("%+v", v)
	}
}

func deepEqual(v1, v2 reflect.Value, visited map[visit]struct{}, path string) (bool, string) {
	if !v1.IsValid() || !v2.IsValid() {
		return boolResult(v1.IsValid(), v2.IsValid(), path, "valid")
	}

	if v1.Type() != v2.Type() {
		return false, fmt.Sprintf(
			"got%s has type %s, want%s has type %s",
			path,
			v1.Type().String(),
			path,
			v2.Type().String(),
		)
	}

	v1Kind := v1.Kind()
	if v1.CanAddr() && v2.CanAddr() && trackedType(v1Kind) {
		v1Addr := unsafe.Pointer(v1.UnsafeAddr())
		v2Addr := unsafe.Pointer(v2.UnsafeAddr())
		// Construct visit with smaller address first to handle nested
		// comparisons with got/want reversed.
		if uintptr(v1Addr) > uintptr(v2Addr) {
			v1Addr, v2Addr = v2Addr, v1Addr
		}

		// Short-circuit if we've already done (or are currently
		// doing) this comparison.
		v := visit{v1Addr, v2Addr, v1.Type()}
		if _, ok := visited[v]; ok {
			return true, ""
		}

		visited[v] = struct{}{}
	}

	switch v1Kind {
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return boolResult(v1.IsNil(), v2.IsNil(), path, "nil")
		}

		if v1.Pointer() == v2.Pointer() && v1.Len() == v2.Len() {
			// same instance
			return true, ""
		}
		fallthrough

	case reflect.Array:
		allMatched := true
		fullReason := ""
		i := 0
		for ; i < v1.Len() && i < v2.Len(); i++ {
			ok, reason := deepEqual(
				v1.Index(i),
				v2.Index(i),
				visited,
				fmt.Sprintf("%s[%d]", path, i),
			)
			if !ok {
				fullReason += reason + "\n"
				allMatched = false
			}
		}

		if i < v1.Len() {
			for j := i; j < v1.Len(); j++ {
				fullReason += fmt.Sprintf(
					"got%s[%d] is %s, no want%s[%d] given\n",
					path,
					j,
					render(v1.Type().Elem(), v1.Index(j)),
					path,
					j,
				)
			}
			allMatched = false
		}

		if i < v2.Len() {
			for j := i; j < v2.Len(); j++ {
				fullReason += fmt.Sprintf(
					"got%s[%d] missing, want%s[%d] is %s\n",
					path,
					j,
					path,
					j,
					render(v2.Type().Elem(), v2.Index(j)),
				)
			}
			allMatched = false
		}

		return allMatched, strings.TrimSpace(fullReason)

	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return boolResult(v1.IsNil(), v2.IsNil(), path, "nil")
		}
		return deepEqual(
			v1.Elem(),
			v2.Elem(),
			visited,
			fmt.Sprintf("%s.(%s)", path, v1.Elem().Type().String()),
		)

	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			ok, reason := deepEqual(
				v1.Field(i),
				v2.Field(i),
				visited,
				fmt.Sprintf("%s.%s", path, v1.Type().Field(i).Name),
			)
			if !ok {
				return false, reason
			}
		}
		return true, ""

	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			// same instance
			return true, ""
		}
		if v1.IsNil() || v2.IsNil() {
			return boolResult(v1.IsNil(), v2.IsNil(), path, "nil")
		}
		return deepEqual(v1.Elem(), v2.Elem(), visited, path)

	case reflect.Map:
		if v1.IsNil() || v2.IsNil() {
			return boolResult(v1.IsNil(), v2.IsNil(), path, "nil")
		}

		if v1.Len() != v2.Len() {
			return false, fmt.Sprintf(
				"got%s is a map[%s]%s with %d entries, want has %d entries",
				path,
				v1.Type().Key().String(),
				v1.Type().Elem().String(),
				v1.Len(),
				v2.Len(),
			)
		}

		if v1.Pointer() == v2.Pointer() {
			// same instance
			return true, ""
		}

		for _, k := range v1.MapKeys() {
			mapV1 := v1.MapIndex(k)
			if !mapV1.IsValid() {
				// This shouldn't be possible, since we got this key
				// from v1, but go's implementation checks this.
				return false, fmt.Sprintf(
					"got%s[%s]: no value for key",
					path,
					k.String(),
				)
			}
			mapV2 := v2.MapIndex(k)
			if !mapV2.IsValid() {
				// v2 doesn't have this key.
				return false, fmt.Sprintf(
					"got%s[%#v] is %#v, want%s[%#v] is missing",
					path,
					k,
					mapV1,
					path,
					k,
				)
			}
			ok, reason := deepEqual(mapV1, mapV2, visited, fmt.Sprintf("%s[%#v]", path, k))
			if !ok {
				return false, reason
			}
		}
		return true, ""

	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		// N.B. go's reflect.DeepEqual only indicates equality for
		// Funcs when both are nil. We allow pointer-identical
		// functions to be equal.
		if v1.Pointer() != v2.Pointer() {
			return false, fmt.Sprintf(
				"got%s is %s %#x, want%s is %s %#x",
				path,
				v1.Kind().String(),
				v1.Pointer(),
				path,
				v1.Kind().String(),
				v2.Pointer(),
			)
		}
		return true, ""

	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			return false, fmt.Sprintf("got%s is %t, want%s is %t", path, v1.Bool(), path, v2.Bool())
		}
		return true, ""

	case reflect.String:
		if v1.String() != v2.String() {
			return false, fmt.Sprintf(
				"got%s is %q, want%s is %q",
				path,
				v1.String(),
				path,
				v2.String(),
			)
		}
		return true, ""

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v1.Int() != v2.Int() {
			return false, fmt.Sprintf("got%s is %d, want%s is %d", path, v1.Int(), path, v2.Int())
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		if v1.Uint() != v2.Uint() {
			return false, fmt.Sprintf("got%s is %d, want%s is %d", path, v1.Uint(), path, v2.Uint())
		}
		return true, ""

	case reflect.Float32, reflect.Float64:
		if v1.Float() != v2.Float() {
			return false, fmt.Sprintf(
				"got%s is %g, want%s is %g",
				path,
				v1.Float(),
				path,
				v2.Float(),
			)
		}
		return true, ""

	case reflect.Complex64, reflect.Complex128:
		if v1.Complex() != v2.Complex() {
			return false, fmt.Sprintf(
				"got%s is %g, want%s is %g",
				path,
				v1.Complex(),
				path,
				v2.Complex(),
			)
		}
		return true, ""

	default:
		panic(fmt.Sprintf("unknown kind %s", v1Kind.String()))
	}
}

// DeepEqual compares two objects as in reflect.DeepEqual. If the
// result is false, the returned string contains a description of how
// the objects differed. The only notable difference between DeepEqual
// and reflect.DeepEqual is that this function will return true for
// pointer-identical (e.g., same instance) channels and functions.
func DeepEqual(got, want interface{}) (bool, string) {
	if got == nil {
		if want == nil {
			return true, ""
		}

		return false, "got is nil, but want is non-nil"
	} else if want == nil {
		return false, "got is non-nil, but want is nil"
	}

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	if gotValue.Type() != wantValue.Type() {
		return false, fmt.Sprintf("got is of type %T, want is of type %T", got, want)
	}

	return deepEqual(gotValue, wantValue, map[visit]struct{}{}, "")
}
