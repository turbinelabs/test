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

// The assert package include various simple test assertions, and a mechanism
// to group assertions together.
package assert

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/turbinelabs/test/check"
	tbnstr "github.com/turbinelabs/test/strings"
)

const (
	goPath      = "/usr/local/"
	tbnHomePath = "TBN_FULL_HOME"
)

// Nil asserts the nilness of got.
func Nil(t testing.TB, got interface{}) bool {
	if !check.IsNil(got) {
		Tracing(t).Errorf("got (%T) %s, want <nil>", got, tbnstr.Stringify(got))
		return false
	}
	return true
}

// NonNil asserts the non-nilness of got.
func NonNil(t testing.TB, got interface{}) bool {
	if check.IsNil(got) {
		Tracing(t).Errorf("got (%T) %s, want <non-nil>", got, tbnstr.Stringify(got))
		return false
	}
	return true
}

func mkErrorMsg(got, want interface{}) string {
	return mkErrorMsgWithExp(got, want, "want")
}

func mkErrorMsgWithExp(got, want interface{}, expectation string) string {
	return fmt.Sprintf(
		"got (%T) %s, %s (%T) %s",
		got,
		tbnstr.Stringify(got),
		expectation,
		want,
		tbnstr.Stringify(want),
	)
}

// Equal asserts that got == want, and will panic for types that can't
// be compared with ==.
func Equal(t testing.TB, got, want interface{}) bool {
	if got != want {
		Tracing(t).Error(mkErrorMsg(got, want))
		return false
	}
	return true
}

// Equal asserts that got != want, and will panic for types that can't
// be compared with !=.
func NotEqual(t testing.TB, got, want interface{}) bool {
	if got == want {
		Tracing(t).Error(mkErrorMsgWithExp(got, want, "want !="))
		return false
	}
	return true
}

// EqualWithin asserts that two floating point numbers are within a
// given epsilon of each other. If both numbers are NaN, +Inf, or
// -Inf, they are equal irrespective of the value of epsilon.
func EqualWithin(t testing.TB, got, want, epsilon float64) bool {
	if !check.EqualWithin(got, want, epsilon) {
		Tracing(t).Error(fmt.Sprintf("got %g, want %g (within %g)", got, want, epsilon))
		return false
	}
	return true
}

// NotEqualWithin asserts that two floating point numbers are not
// within a given epsilon of each other. If both numbers are Nan,
// +Inf, -Inf they are considered equal irrespective of epsilon (and
// this assertion will fail).
func NotEqualWithin(t testing.TB, got, want, epsilon float64) bool {
	if check.EqualWithin(got, want, epsilon) {
		Tracing(t).Error(fmt.Sprintf("got %g, want != %g (within %g)", got, want, epsilon))
		return false
	}
	return true
}

func isInt(i interface{}) (int64, bool) {
	switch x := i.(type) {
	case int:
		return int64(x), true

	case int8:
		return int64(x), true

	case int16:
		return int64(x), true

	case int32:
		return int64(x), true

	case int64:
		return x, true

	case time.Duration:
		return int64(x), true

	default:
		return 0, false
	}
}

func isUint(i interface{}) (uint64, bool) {
	switch x := i.(type) {
	case uint:
		return uint64(x), true

	case uint8:
		return uint64(x), true

	case uint16:
		return uint64(x), true

	case uint32:
		return uint64(x), true

	case uint64:
		return x, true

	default:
		return 0, false
	}
}

func isFloat(i interface{}) (float64, bool) {
	switch x := i.(type) {
	case float32:
		f64 := float64(x)
		if math.IsNaN(f64) {
			return 0.0, false
		}
		return f64, true

	case float64:
		if math.IsNaN(x) {
			return 0.0, false
		}
		return x, true

	default:
		return 0.0, false
	}
}

func compare(got, cmp interface{}, f func(int) bool) (bool, error) {
	cmpInt := func(a, b int64) bool {
		if a < b {
			return f(-1)
		}
		if a > b {
			return f(1)
		}
		return f(0)
	}

	cmpUint := func(a, b uint64) bool {
		if a < b {
			return f(-1)
		}
		if a > b {
			return f(1)
		}
		return f(0)
	}

	cmpFlt := func(a, b float64) bool {
		if a < b {
			return f(-1)
		}
		if a > b {
			return f(1)
		}
		return f(0)
	}

	if gotIVal, isGotIVal := isInt(got); isGotIVal {
		if cmpIVal, isCmpIVal := isInt(cmp); isCmpIVal {
			return cmpInt(gotIVal, cmpIVal), nil
		}

		if cmpUVal, isCmpUVal := isUint(cmp); isCmpUVal {
			if gotIVal < 0 {
				// gotIVal < cmpUVal by range
				return f(-1), nil
			}

			return cmpUint(uint64(gotIVal), cmpUVal), nil
		}

		if cmpFVal, isCmpFVal := isFloat(cmp); isCmpFVal {
			return cmpFlt(float64(gotIVal), cmpFVal), nil
		}
	} else if gotUVal, isGotUVal := isUint(got); isGotUVal {
		if cmpUVal, isCmpUVal := isUint(cmp); isCmpUVal {
			return cmpUint(gotUVal, cmpUVal), nil
		}

		if cmpIVal, isCmpIVal := isInt(cmp); isCmpIVal {
			if cmpIVal < 0 {
				// gotUVal > cmpIVal by range
				return f(1), nil
			}

			return cmpUint(gotUVal, uint64(cmpIVal)), nil
		}

		if cmpFVal, isCmpFVal := isFloat(cmp); isCmpFVal {
			return cmpFlt(float64(gotUVal), cmpFVal), nil
		}
	} else if gotFVal, isGotFVal := isFloat(got); isGotFVal {
		if cmpFVal, isCmpFVal := isFloat(cmp); isCmpFVal {
			return cmpFlt(gotFVal, cmpFVal), nil
		}

		if cmpIVal, isCmpIVal := isInt(cmp); isCmpIVal {
			return cmpFlt(gotFVal, float64(cmpIVal)), nil
		}

		if cmpUVal, isCmpUVal := isUint(cmp); isCmpUVal {
			return cmpFlt(gotFVal, float64(cmpUVal)), nil
		}
	}

	return false, fmt.Errorf("cannot compare %T with %T", got, cmp)
}

// GreaterThan compares two numeric values to see if got >
// comparator. Integer types may be mixed: the types are coerced to
// 64-bit integers for comparison. Similarly, floating point types may
// be mixed. Finally, integer and floating point types may be mixed:
// the got value is converted to the comparator's type. It is possible
// to write assertions that cannot fail: For example, if the
// comparator's value is greater than the got value's type allows.
func GreaterThan(t testing.TB, got, comparator interface{}) bool {
	isGreater, err := compare(got, comparator, func(c int) bool { return c > 0 })
	if err != nil {
		Tracing(t).Error(err.Error())
		return false
	}

	if !isGreater {
		Tracing(t).Error(mkErrorMsgWithExp(got, comparator, "want >"))
		return false
	}

	return true
}

// GreaterThanEqual compares two numeric values to see if got >=
// comparator. See discussion of numeric types in GreaterThan.
func GreaterThanEqual(t testing.TB, got, comparator interface{}) bool {
	isGreaterEq, err := compare(got, comparator, func(c int) bool { return c >= 0 })
	if err != nil {
		Tracing(t).Error(err.Error())
		return false
	}

	if !isGreaterEq {
		Tracing(t).Error(mkErrorMsgWithExp(got, comparator, "want >="))
		return false
	}

	return true
}

// LessThan compares two numeric values to see if got <
// comparator. See discussion of numeric types in GreaterThan.
func LessThan(t testing.TB, got, comparator interface{}) bool {
	isLess, err := compare(got, comparator, func(c int) bool { return c < 0 })
	if err != nil {
		Tracing(t).Error(err.Error())
		return false
	}

	if !isLess {
		Tracing(t).Error(mkErrorMsgWithExp(got, comparator, "want <"))
		return false
	}

	return true
}

// LessThanEqual compares two numeric values to see if got <=
// comparator. See discussion of numeric types in GreaterThan.
func LessThanEqual(t testing.TB, got, comparator interface{}) bool {
	isLessEq, err := compare(got, comparator, func(c int) bool { return c <= 0 })
	if err != nil {
		Tracing(t).Error(err.Error())
		return false
	}

	if !isLessEq {
		Tracing(t).Error(mkErrorMsgWithExp(got, comparator, "want <="))
		return false
	}

	return true
}

func isArrayLike(i interface{}) bool {
	t := reflect.TypeOf(i)
	if t == nil {
		return false
	}
	kind := t.Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

// panics if a is not an array
func arrayValues(a interface{}) []reflect.Value {
	aValue := reflect.ValueOf(a)
	if aValue.Kind() != reflect.Array && aValue.IsNil() {
		return nil
	}
	valueArray := make([]reflect.Value, aValue.Len())
	for i := range valueArray {
		valueArray[i] = aValue.Index(i)
	}
	return valueArray
}

// ArrayEqual compares two arrays for equality. Arrays may be compared
// with slices and vice versa. Nil arrays/slices are not equal to
// empty arrays/slices. Array/slice elements are compared as in
// DeepEqual. The assertion error messages indicate the index at which
// an inequality occurred and report extra or missing values.
func ArrayEqual(t testing.TB, got, want interface{}) bool {
	if !isArrayLike(got) || !isArrayLike(want) {
		Tracing(t).Error(mkErrorMsg(got, want))
		return false
	}

	gotValues := arrayValues(got)
	wantValues := arrayValues(want)

	if gotValues == nil && wantValues != nil {
		Tracing(t).Errorf("got (%T) nil, want (%T) %s", got, want, tbnstr.Stringify(want))
		return false
	} else if wantValues == nil && gotValues != nil {
		Tracing(t).Errorf("got (%T) %s, want (%T) nil", got, tbnstr.Stringify(got), want)
		return false
	}

	gotLen := len(gotValues)
	wantLen := len(wantValues)

	errors := []string{}
	for i := 0; i < gotLen || i < wantLen; i++ {
		var gotIface, wantIface interface{}
		gotValid := i < gotLen
		if gotValid {
			gotIface = gotValues[i].Interface()
		}

		wantValid := i < wantLen
		if wantValid {
			wantIface = wantValues[i].Interface()
		}

		var err string
		if gotValid && wantValid {
			if !reflect.DeepEqual(gotIface, wantIface) {
				err = fmt.Sprintf(
					"index %d: %s",
					i,
					mkErrorMsg(gotIface, wantIface),
				)
			}
		} else if gotValid {
			err = fmt.Sprintf(
				"index %d: got extra value: (%T) %s",
				i,
				gotIface,
				tbnstr.Stringify(gotIface),
			)
		} else if wantValid {
			err = fmt.Sprintf(
				"index %d: missing wanted value: (%T) %s",
				i,
				wantIface,
				tbnstr.Stringify(wantIface),
			)
		}

		if err != "" {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		Tracing(t).Errorf("arrays not equal:\n%s", strings.Join(errors, "\n"))
		return false
	}
	return true
}

func isMap(i interface{}) bool {
	t := reflect.TypeOf(i)
	return t != nil && t.Kind() == reflect.Map
}

// MapEqual compares to maps for equality. Nil maps are not equal to
// empty maps. Values are retrieved from both maps for each key in the
// want map. Values are compared as in DeepEqual. Missing and extra
// entries are reported in the error messages.
func MapEqual(t testing.TB, got, want interface{}) bool {
	if !isMap(got) || !isMap(want) {
		Tracing(t).Error(mkErrorMsg(got, want))
		return false
	}

	wantValue := reflect.ValueOf(want)
	wantKeys := wantValue.MapKeys()

	gotValue := reflect.ValueOf(got)
	gotKeys := gotValue.MapKeys()

	if gotValue.IsNil() && !wantValue.IsNil() {
		Tracing(t).Errorf("got (%T) nil, want (%T) %s", got, want, tbnstr.Stringify(want))
		return false
	} else if wantValue.IsNil() && !gotValue.IsNil() {
		Tracing(t).Errorf("got (%T) %s, want (%T) nil", got, tbnstr.Stringify(got), want)
		return false
	}

	errors := []string{}
	for _, wantKey := range wantKeys {
		wantIface := wantValue.MapIndex(wantKey).Interface()

		var err string
		gotMapValue := gotValue.MapIndex(wantKey)
		if gotMapValue.IsValid() {
			gotIface := gotMapValue.Interface()
			if !reflect.DeepEqual(gotIface, wantIface) {
				err = fmt.Sprintf(
					"key %s: %s",
					tbnstr.Stringify(wantKey.Interface()),
					mkErrorMsg(gotIface, wantIface),
				)
			}
		} else {
			err = fmt.Sprintf(
				"missing key %s: wanted value: (%T) %s",
				tbnstr.Stringify(wantKey.Interface()),
				wantIface,
				tbnstr.Stringify(wantIface),
			)
		}
		if err != "" {
			errors = append(errors, err)
		}
	}

	for _, gotKey := range gotKeys {
		wantMapValue := wantValue.MapIndex(gotKey)
		if !wantMapValue.IsValid() {
			gotIface := gotValue.MapIndex(gotKey).Interface()
			err := fmt.Sprintf(
				"extra key %s: unwanted value: (%T) %s",
				tbnstr.Stringify(gotKey.Interface()),
				gotIface,
				tbnstr.Stringify(gotIface),
			)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		Tracing(t).Errorf("maps not equal:\n%s", strings.Join(errors, "\n"))
		return false
	}

	return true
}

// DeepEqual asserts reflect.DeepEqual(got, want).
func DeepEqual(t testing.TB, got, want interface{}) bool {
	if isArrayLike(got) && isArrayLike(want) {
		return ArrayEqual(t, got, want)
	} else if isMap(got) && isMap(want) {
		return MapEqual(t, got, want)
	} else if !reflect.DeepEqual(got, want) {
		Tracing(t).Error(mkErrorMsg(got, want))
		return false
	}
	return true
}

// NotDeepEqual asserts !reflect.DeepEqual(got, want).
func NotDeepEqual(t testing.TB, got, want interface{}) bool {
	if reflect.DeepEqual(got, want) {
		Tracing(t).Error(mkErrorMsgWithExp(got, want, "want !="))
		return false
	}
	return true
}

func sameInstance(got, want interface{}) bool {
	gotType := reflect.TypeOf(got)
	if gotType != reflect.TypeOf(want) {
		return false
	}
	switch gotType.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:

		gotVal := reflect.ValueOf(got)
		wantVal := reflect.ValueOf(want)
		if gotVal.Pointer() != wantVal.Pointer() {
			return false
		}
		// slices of different lengths can still share a pointer
		if gotType.Kind() == reflect.Slice && gotVal.Len() != wantVal.Len() {
			return false
		}
		return true

	default:
		panic(fmt.Sprintf(
			"cannot determine instance equality for non-pointer type: %T",
			got,
		))
	}
}

// SameInstance asserts that got and want are the same instance of a
// pointer type, or are the same value of a literal type.
func SameInstance(t testing.TB, got, want interface{}) bool {
	if !sameInstance(got, want) {
		Tracing(t).Error(mkErrorMsgWithExp(got, want, "want same instance as"))
		return false
	}
	return true
}

// NotSameInstance asserts that got and want are not the same instance of a
// pointer type, and are not the same value of a literal type.
func NotSameInstance(t testing.TB, got, want interface{}) bool {
	if sameInstance(got, want) {
		Tracing(t).Error(mkErrorMsgWithExp(got, want, "want not same instance as"))
		return false
	}
	return true
}

func encodeJson(t testing.TB, got, want interface{}) (string, string) {
	gotJson, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("could not marshal json: %#v", err)
	}
	wantJson, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("could not marshal json: %#v", err)
	}
	return string(gotJson), string(wantJson)
}

// EqualJson asserts that got and want encode to the same JSON value.
func EqualJson(t testing.TB, got, want interface{}) bool {
	tr := Tracing(t)
	gotJson, wantJson := encodeJson(tr, got, want)
	return Equal(tr, gotJson, wantJson)
}

// NotEqualJson asserts that got and want do not encode to the same JSON value.
func NotEqualJson(t testing.TB, got, want interface{}) bool {
	tr := Tracing(t)
	gotJson, wantJson := encodeJson(tr, got, want)
	return NotEqual(tr, gotJson, wantJson)
}

func matchRegex(t testing.TB, got, wantRegex string) bool {
	matched, err := regexp.MatchString(wantRegex, got)
	if err != nil {
		t.Fatalf("invalid regular expression `%s`: %#v", wantRegex, err)
	}

	return matched
}

// MatchesRegex asserts that got matches the provided regular expression.
func MatchesRegex(t testing.TB, got, wantRegex string) bool {
	tr := Tracing(t)
	if !matchRegex(tr, got, wantRegex) {
		tr.Errorf("got %q, did not match `%s`", got, wantRegex)
		return false
	}

	return true
}

// DoesNotMatchRegex asserts that got does not match the provided regular expression.
func DoesNotMatchRegex(t testing.TB, got, wantRegex string) bool {
	tr := Tracing(t)
	if matchRegex(tr, got, wantRegex) {
		tr.Errorf("got %q, matched `%s`", got, wantRegex)
		return false
	}

	return true
}

// True asserts that value is true.
func True(t testing.TB, value bool) bool {
	return Equal(Tracing(t), value, true)
}

// False asserts that value is false.
func False(t testing.TB, value bool) bool {
	return Equal(Tracing(t), value, false)
}

// Failed logs msg and aborts the current test.
func Failed(t testing.TB, msg string) {
	Tracing(t).Fatalf("Failed: %s", msg)
}

// ErrorContains asserts that got.Error() contains want.
func ErrorContains(t testing.TB, got error, want string) bool {
	tr := Tracing(t)
	if got == nil {
		tr.Errorf("got nil error, wanted message containing %s", tbnstr.Stringify(want))
		return false
	} else if !strings.Contains(got.Error(), want) {
		tr.Errorf(
			"got error %s, wanted message containing %s",
			tbnstr.Stringify(got.Error()),
			tbnstr.Stringify(want),
		)
		return false
	}

	return true
}

// ErrorDoesNotContain asserts that got.Error() does not contain want.
func ErrorDoesNotContain(t testing.TB, got error, want string) bool {
	tr := Tracing(t)
	if got == nil {
		tr.Errorf("got nil error, wanted message not containing %s", tbnstr.Stringify(want))
		return false
	} else if strings.Contains(got.Error(), want) {
		tr.Errorf(
			"got error %s, wanted message not containing %s",
			tbnstr.Stringify(got.Error()),
			tbnstr.Stringify(want),
		)
		return false
	}

	return true
}

// StringContains asserts that got contains want.
func StringContains(t testing.TB, got, want string) bool {
	tr := Tracing(t)
	if !strings.Contains(got, want) {
		tr.Errorf(
			"got %s, wanted message containing %s",
			tbnstr.Stringify(got),
			tbnstr.Stringify(want),
		)
		return false
	}

	return true
}

// StringDoesNotContain asserts that got does not contain want.
func StringDoesNotContain(t testing.TB, got, want string) bool {
	tr := Tracing(t)
	if strings.Contains(got, want) {
		tr.Errorf(
			"got %s, wanted message not containing %s",
			tbnstr.Stringify(got),
			tbnstr.Stringify(want),
		)
		return false
	}

	return true
}

// HasSameElements compares two container-like values. The got
// parameter may be an array, slice, or channel. The want parameter
// must be an array or slice whose element type is the same as that of
// got. If got is a channel, all available values are consumed (until
// the channel either blocks or indicates it was closed). The got and
// want values are then compared (as in DeepEqual) without respect to
// order.
func HasSameElements(t testing.TB, got, want interface{}) bool {
	if err := check.HasSameElements(got, want); err != nil {
		Tracing(t).Error(err.Error())
		return false
	}

	return true
}

// v must be a zero-arg function
func checkPanic(v reflect.Value) (i interface{}) {
	defer func() {
		if x := recover(); x != nil {
			i = x
		}
	}()

	v.Call(nil)
	return
}

// Panic asserts that the given function panics. The f parameter must
// be a function that takes no arguments. It may, however, return any
// number of arguments.
func Panic(t testing.TB, f interface{}) bool {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		Tracing(t).Errorf("parameter to Panic must be a function: %+v", f)
		return false
	}
	if fType.NumIn() != 0 {
		Tracing(t).Errorf("function passed to Panic may not take arguments: %+v", f)
		return false
	}

	if checkPanic(reflect.ValueOf(f)) == nil {
		Tracing(t).Error("expected panic")
		return false
	}
	return true
}
