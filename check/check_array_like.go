package check

import (
	"fmt"
	"reflect"
	"strings"

	tbnstr "github.com/turbinelabs/test/strings"
)

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

func checkContainerTypes(gotType, wantType reflect.Type) error {
	gotKind := gotType.Kind()
	wantKind := wantType.Kind()

	switch gotKind {
	case reflect.Array, reflect.Slice:
		// ok

	case reflect.Chan:
		if gotType.ChanDir()&reflect.RecvDir == 0 {
			return fmt.Errorf("got type '%v', a non-receiving channel", gotType)
		}

	default:
		return fmt.Errorf(
			"got type '%v', can only compare arrays, slices, or channels",
			gotType,
		)
	}

	if wantKind != reflect.Array && wantKind != reflect.Slice {
		// We only compare with Array/Slices
		return fmt.Errorf(
			"got type '%v', want type must be an array or slice of %s, not '%v'",
			gotType,
			gotType.Elem(),
			wantType)
	}

	// The Array/Slice/Chan element types must match
	if gotType.Elem() != wantType.Elem() {
		return fmt.Errorf(
			"got type '%v', wanted type '%v': contains types do not match",
			gotType,
			wantType)
	}

	return nil
}

type multiArray interface {
	Len() int
	InnerLen(int) int
	Value(int, int) interface{}
}

type valueMultiArray [][]reflect.Value

func (a valueMultiArray) Len() int                   { return len(a) }
func (a valueMultiArray) InnerLen(i int) int         { return len(a[i]) }
func (a valueMultiArray) Value(i, j int) interface{} { return a[i][j] }

var _ multiArray = valueMultiArray{}

type ifaceMultiArray [][]interface{}

func (a ifaceMultiArray) Len() int                   { return len(a) }
func (a ifaceMultiArray) InnerLen(i int) int         { return len(a[i]) }
func (a ifaceMultiArray) Value(i, j int) interface{} { return a[i][j] }

var _ multiArray = ifaceMultiArray{}

func formattedIfaceArrayStrings(
	format func(i interface{}) string,
	arrays multiArray,
) []string {
	if arrays.Len() == 0 {
		return []string{}
	}

	inline := true
	ifaceStrs := make([][]string, arrays.Len())
	for i := 0; i < arrays.Len(); i++ {
		strs := make([]string, arrays.InnerLen(i))
		length := 0

		for j := 0; j < arrays.InnerLen(i); j++ {
			s := format(arrays.Value(i, j))
			strs[j] = s
			length += len(s)
		}

		if length > 40 {
			inline = false
		}
		ifaceStrs[i] = strs
	}

	results := make([]string, arrays.Len())
	for i, strs := range ifaceStrs {
		if inline {
			results[i] = fmt.Sprintf("[%s]", strings.Join(strs, ", "))
		} else {
			results[i] = fmt.Sprintf("[\n%s\n]", strings.Join(strs, ",\n"))
		}
	}

	return results
}

func ifaceArrayStrings(ifaceArrays ...[]interface{}) []string {
	return formattedIfaceArrayStrings(
		func(i interface{}) string {
			return fmt.Sprintf("(%T) %s", i, tbnstr.Stringify(i))
		},
		ifaceMultiArray(ifaceArrays),
	)
}

func ifaceArrayString(ifaceArrays []interface{}) string {
	return ifaceArrayStrings(ifaceArrays)[0]
}

func valueArrayStrings(valueArrays ...[]reflect.Value) []string {
	return formattedIfaceArrayStrings(
		func(i interface{}) string {
			v := i.(reflect.Value)
			return fmt.Sprintf(
				"(%s) %s",
				v.Type().Name(),
				tbnstr.Stringify(v.Interface()),
			)
		},
		valueMultiArray(valueArrays),
	)
}

func assertSameArray(gotValue, wantValue []reflect.Value) error {
	gotLen := len(gotValue)
	wantLen := len(wantValue)

	unusedGotIndicies := make([]int, gotLen)
	for i := 0; i < gotLen; i++ {
		unusedGotIndicies[i] = i
	}

	unusedWantIndicies := make([]int, wantLen)
	for i := 0; i < wantLen; i++ {
		unusedWantIndicies[i] = i
	}

	for gotIndex, v := range gotValue {
		for _, wantIndex := range unusedWantIndicies {
			if wantIndex != -1 {
				w := wantValue[wantIndex]
				if reflect.DeepEqual(v.Interface(), w.Interface()) {
					unusedWantIndicies[wantIndex] = -1
					unusedGotIndicies[gotIndex] = -1
					break
				}
			}
		}
	}

	extra := []interface{}{}
	for _, gotIndex := range unusedGotIndicies {
		if gotIndex != -1 {
			extra = append(extra, gotValue[gotIndex].Interface())
		}
	}

	missing := []interface{}{}
	for _, wantIndex := range unusedWantIndicies {
		if wantIndex != -1 {
			missing = append(missing, wantValue[wantIndex].Interface())
		}
	}

	if gotLen != wantLen || len(extra) > 0 || len(missing) > 0 {
		missingStr := ""
		if len(missing) > 0 {
			missingStr =
				fmt.Sprintf(";\n missing elements: %s", ifaceArrayString(missing))
		}

		extraStr := ""
		if len(extra) > 0 {
			extraStr = fmt.Sprintf(";\nextra elements: %s", ifaceArrayString(extra))
		}

		gotWantStrs := valueArrayStrings(gotValue, wantValue)
		gotValueStr := gotWantStrs[0]
		wantValueStr := gotWantStrs[1]

		return fmt.Errorf(
			"got %s (len %d),\nwanted %s (len %d)%s%s",
			gotValueStr,
			gotLen,
			wantValueStr,
			wantLen,
			missingStr,
			extraStr)
	}

	return nil
}

// Compares two container-like values. The got parameter may be an
// array, slice, or channel. The want parameter must be an array or
// slice whose element type is the same as that of got. If got is a
// channel, all available values are consumed (until the channel
// either blocks or indicates it was closed). The got and want values
// are then compared without respect to order. Returns nil if the
// arrays were comparable and contained the same elements
func HasSameElements(got, want interface{}) error {
	gotType := reflect.TypeOf(got)
	wantType := reflect.TypeOf(want)
	if err := checkContainerTypes(gotType, wantType); err != nil {
		return err
	}

	gotValue := reflect.ValueOf(got)

	wantValueArray := arrayValues(want)

	switch gotType.Kind() {
	case reflect.Array, reflect.Slice:
		gotValueArray := arrayValues(got)
		return assertSameArray(gotValueArray, wantValueArray)

	case reflect.Chan:
		gotValueArray := []reflect.Value{}
		for {
			v, ok := gotValue.TryRecv()
			if !ok {
				// blocked or closed
				break
			}
			gotValueArray = append(gotValueArray, v)
		}
		return assertSameArray(gotValueArray, wantValueArray)

	default:
		return fmt.Errorf(
			"internal error: unexpected kind %v",
			gotType.Kind())
	}
}
