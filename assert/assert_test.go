package assert

import (
	"fmt"
	"reflect"
	"testing"
)

type complexStruct struct {
	x int
	y *string
}

type equalityKind int

const (
	notEqual equalityKind = iota
	justEqual
	justDeepEqual
	justJsonEqual
	equalAndDeepEqual
)

type nilnessKind bool

const (
	notNilish nilnessKind = false
	nilish                = true
)

type equalTestCase struct {
	name string
	a    interface{}
	b    interface{}
	kind equalityKind
}

func (e equalTestCase) run(
	t *testing.T,
	f func(testing.TB, interface{}, interface{}) bool,
	expectEqual func(equalTestCase) bool,
) {
	testT := new(testing.T)

	result := f(testT, e.a, e.b)
	expectedResult := expectEqual(e)

	if result != expectedResult {
		comparison := "equal"
		if !expectedResult {
			comparison = "not equal"
		}
		t.Errorf("%s: expected %+v to %s %+v", e.name, e.a, comparison, e.b)
	}
}

type nilTestCase struct {
	name string
	v    interface{}
	kind nilnessKind
}

func (n nilTestCase) run(
	t *testing.T,
	f func(testing.TB, interface{}) bool,
	expectNil func(nilTestCase) bool,
) {
	testT := new(testing.T)

	result := f(testT, n.v)
	expectedResult := expectNil(n)

	if result != expectedResult {
		t.Errorf("%s: expected %t, got %t for %+v", n.name, expectedResult, result, n.v)
	}
}

type logEntry struct {
	op   string
	args string
}

type mockT struct {
	testing.TB

	log []logEntry
}

func (t *mockT) record(op string, args string) {
	entry := logEntry{op, args}
	if t.log == nil {
		t.log = make([]logEntry, 1)
		t.log[0] = entry
	} else {
		t.log = append(t.log, entry)
	}
}

// For any `testing.TB` method invoked on `mockT`, you'll need to
// override the version inherited from embedding a `testing.TB` in
// `mockT`. (The embedded versions of the methods will fail due to the
// `TB` field being nil.)
func (t *mockT) Errorf(format string, args ...interface{}) {
	t.record("Errorf", fmt.Sprintf(format, args...))
}

type moreComplexStruct struct {
	A string               `json:"a"`
	C lessComplexSubstruct `json:"c"`
}
type lessComplexSubstruct struct {
	D string `json:"d"`
}

var (
	nilStringPtr *string   = nil
	nilStructPtr *struct{} = nil
	nilChannel   chan<- bool
	nilFunction  func()
	nilSlice     []string

	channel  = make(chan<- bool)
	function = func() {}

	int1a = 123
	int1b = 123
	int2  = 456

	string1a = "string"
	string1b = "string"
	string2  = "other string"

	cs1a = complexStruct{1, &string1a}
	cs1b = complexStruct{1, &string1a}
	cs2a = complexStruct{1, &string1a}
	cs2b = complexStruct{1, &string1b}
	cs3  = complexStruct{1, &string1a}
	cs4  = complexStruct{1, &string2}

	slice1a = []string{"a", "b", "c"}
	slice1b = []string{"a", "b", "c"}
	slice2  = []string{"X", "Y", "Z"}

	map1    = map[string]interface{}{"a": "b", "c": map[string]string{"d": "e"}}
	struct1 = moreComplexStruct{A: "b", C: lessComplexSubstruct{D: "e"}}
	map2    = map[string]interface{}{"a": "b", "c": map[string]string{"d": "z"}}
	struct2 = moreComplexStruct{A: "b", C: lessComplexSubstruct{D: "z"}}

	nilnessTestCases = []nilTestCase{
		nilTestCase{"nil", nil, nilish},
		nilTestCase{"*string-nil", nilStringPtr, nilish},
		nilTestCase{"*struct-nil", nilStructPtr, nilish},
		nilTestCase{"*int", &int1a, notNilish},
		nilTestCase{"*string-1a", &string1a, notNilish},
		nilTestCase{"*struct-1a", &cs1a, notNilish},
		nilTestCase{"[]string-nilish", nilSlice, nilish},
		nilTestCase{"[]string-notnilish", slice1a, notNilish},
		nilTestCase{"chan-nilish", nilChannel, nilish},
		nilTestCase{"chan-notnilish", channel, notNilish},
		nilTestCase{"func-nilish", nilFunction, nilish},
		nilTestCase{"func-notnilish", function, notNilish},
	}

	equalityTestCasesJsonOk = []equalTestCase{
		equalTestCase{"string-1a-1b", string1a, string1b, equalAndDeepEqual},
		equalTestCase{"string-1a-2", string1a, string2, notEqual},
		equalTestCase{"*string-nil-nil", nilStringPtr, nilStringPtr, equalAndDeepEqual},
		equalTestCase{"*string-1a-1b", &string1a, &string1b, justDeepEqual},
		equalTestCase{"*string-1a-2", &string1a, &string2, notEqual},
		equalTestCase{"string & *-1a-1b", string1a, &string1b, justJsonEqual},
		equalTestCase{"string & *-1a-1b", &string1a, string1b, justJsonEqual},
		equalTestCase{"int-1a-1b", int1a, int1b, equalAndDeepEqual},
		equalTestCase{"int-1a-2", int1a, int2, notEqual},
		equalTestCase{"*int-1a-1b", &int1a, &int1b, justDeepEqual},
		equalTestCase{"*int-1a-2", &int1a, &int2, notEqual},
		equalTestCase{"struct-1a-1b", cs1a, cs1b, equalAndDeepEqual},
		equalTestCase{"struct-2a-2b", cs2a, cs2b, justDeepEqual},
		equalTestCase{"struct-3-4", cs3, cs4, justJsonEqual},
		equalTestCase{"*struct-1a-1b", &cs1a, &cs1b, justDeepEqual},
		equalTestCase{"*struct-2a-2b", &cs2a, &cs2b, justDeepEqual},
		equalTestCase{"*struct-3-4", &cs3, &cs4, justJsonEqual},
		equalTestCase{"*[]string-1a-1b", &slice1a, &slice1b, justDeepEqual},
		equalTestCase{"*[]string-1a-2", &slice1a, &slice2, notEqual},
	}

	equalityTestCases = append(
		equalityTestCasesJsonOk,
		// these types cannot be json marshalled
		equalTestCase{"chan", channel, channel, equalAndDeepEqual},
	)

	// these types cannot be compared with == (runtime panic)
	deepEqualityTestCasesJsonOk = []equalTestCase{
		equalTestCase{"[]string-1a-1b", slice1a, slice1b, justDeepEqual},
		equalTestCase{"[]string-1a-2", slice1a, slice2, notEqual},
	}

	deepEqualityTestCases = append(
		deepEqualityTestCasesJsonOk,
		// these types cannot be json marshalled
		equalTestCase{"func", function, function, notEqual},
	)

	justJsonEqualTestCases = []equalTestCase{
		equalTestCase{"map1a-struct", map1, struct1, justJsonEqual},
		equalTestCase{"map1a-map2", map1, map2, notEqual},
		equalTestCase{"map2-struct", map2, struct1, notEqual},
		equalTestCase{"struct1-cs3", struct1, cs3, notEqual},
		equalTestCase{"struct1-struct2", struct1, struct2, notEqual},
	}
)

func TestNonNil(t *testing.T) {
	for _, test := range nilnessTestCases {
		test.run(
			t,
			NonNil,
			func(test nilTestCase) bool {
				return test.kind == notNilish
			})
	}
}

func TestNil(t *testing.T) {
	for _, test := range nilnessTestCases {
		test.run(
			t,
			Nil,
			func(test nilTestCase) bool {
				return test.kind == nilish
			})
	}
}

func TestEqual(t *testing.T) {
	for _, test := range equalityTestCases {
		test.run(
			t,
			Equal,
			func(test equalTestCase) bool {
				return test.kind == justEqual || test.kind == equalAndDeepEqual
			})
	}
}

func TestNotEqual(t *testing.T) {
	for _, test := range equalityTestCases {
		test.run(
			t,
			NotEqual,
			func(test equalTestCase) bool {
				return test.kind == notEqual || test.kind == justDeepEqual || test.kind == justJsonEqual
			})
	}
}

func TestDeepEqual(t *testing.T) {
	testCases := append(equalityTestCases, deepEqualityTestCases...)
	for _, test := range testCases {
		test.run(
			t,
			DeepEqual,
			func(test equalTestCase) bool {
				return test.kind == justDeepEqual || test.kind == equalAndDeepEqual
			})
	}
}

func TestNotDeepEqual(t *testing.T) {
	testCases := append(equalityTestCases, deepEqualityTestCases...)
	for _, test := range testCases {
		test.run(
			t,
			NotDeepEqual,
			func(test equalTestCase) bool {
				return test.kind == notEqual || test.kind == justEqual || test.kind == justJsonEqual
			})
	}
}

func TestEqualJson(t *testing.T) {
	testCases := append(append(equalityTestCasesJsonOk, deepEqualityTestCasesJsonOk...), justJsonEqualTestCases...)
	for _, test := range testCases {
		test.run(
			t,
			EqualJson,
			func(test equalTestCase) bool {
				return test.kind == justDeepEqual || test.kind == equalAndDeepEqual || test.kind == justEqual || test.kind == justJsonEqual
			})
	}
}

func TestNotEqualJson(t *testing.T) {
	testCases := append(append(equalityTestCasesJsonOk, deepEqualityTestCasesJsonOk...), justJsonEqualTestCases...)
	for _, test := range testCases {
		test.run(
			t,
			NotEqualJson,
			func(test equalTestCase) bool {
				return test.kind == notEqual
			})
	}
}

func TestMatchesRegex(t *testing.T) {
	mockT := new(testing.T)

	if !MatchesRegex(mockT, "xyzpdq", "^xyz") {
		t.Errorf("expected 'xyzpdq' to match '^xyz'")
	}
	if !MatchesRegex(mockT, "xyzpdq", "pdq$") {
		t.Errorf("expected 'xyzpdq' to match 'pdq$'")
	}
	if !MatchesRegex(mockT, "xyzpdq", "zp") {
		t.Errorf("expected 'xyzpdq' to match 'zp'")
	}
	if !MatchesRegex(mockT, "xyzpdq", "^xy.+dq$") {
		t.Errorf("expected 'xyzpdq' to match '^xy.+dq$'")
	}

	if MatchesRegex(mockT, "xyzpdq", "a+") {
		t.Errorf("expected 'xyzpdq' to not match 'a+'")
	}
}

func TestDoesNotMatchRegex(t *testing.T) {
	mockT := new(testing.T)

	if DoesNotMatchRegex(mockT, "xyzpdq", "^xyz") {
		t.Errorf("expected 'xyzpdq' to fail by matching '^xyz'")
	}

	if !DoesNotMatchRegex(mockT, "xyzpdq", "a+") {
		t.Errorf("expected 'xyzpdq' to not match 'a+'")
	}
}

func TestErrorContains(t *testing.T) {
	mockT := new(testing.T)

	err := fmt.Errorf("this error contains: magic!")

	if !ErrorContains(mockT, err, "magic") {
		t.Errorf("expected '%s' to contain 'magic'", err.Error())
	}
	if ErrorContains(mockT, err, "special sauce") {
		t.Errorf("expected '%s' not to contain 'special sauce'", err.Error())
	}
	if ErrorContains(mockT, nil, "anything") {
		t.Errorf("expected nil error not to pass check")
	}
}

func TestErrorDoesNotContain(t *testing.T) {
	mockT := new(testing.T)

	err := fmt.Errorf("this error contains: magic!")

	if ErrorDoesNotContain(mockT, err, "magic") {
		t.Errorf("expected '%s' to contain 'magic', but it did", err.Error())
	}
	if !ErrorDoesNotContain(mockT, err, "special sauce") {
		t.Errorf("expected '%s' not to contain 'special sauce'", err.Error())
	}
	if ErrorDoesNotContain(mockT, nil, "anything") {
		t.Errorf("expected nil error not to pass check")
	}
}

func TestHasSameElements(t *testing.T) {
	mockT := new(testing.T)

	expectSame := func(a, b interface{}) {
		if !HasSameElements(mockT, a, b) {
			t.Errorf("expected '%v' to have same elements as '%v'", a, b)
		}
	}

	expectDifferent := func(a, b interface{}) {
		if HasSameElements(mockT, a, b) {
			t.Errorf("expected '%v' to not have same elements as '%v'", a, b)
		}
	}

	a1 := []int{1, 2, 3}
	a2 := []int{3, 2, 1}
	a3 := []int{1, 1, 1}
	a4 := []int{1, 2, 3, 4}
	a5 := []int{1, 1, 1, 2, 2, 2}
	a6 := []int{1, 2, 1, 2, 1, 2}
	a7 := []int{1, 1, 2, 2}

	expectSame(a1, a2)
	expectDifferent(a3, a1)
	expectDifferent(a1, a3)
	expectDifferent(a1, a4)
	expectDifferent(a4, a1)
	expectSame(a5, a6)
	expectDifferent(a5, a7)

	a8 := []complexStruct{cs1a, cs2b}
	a9 := []complexStruct{cs2b, cs1a}
	a10 := []complexStruct{cs1a, cs4}

	expectSame(a8, a9)
	expectDifferent(a8, a10)

	big_array := []int{1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1}
	s1 := big_array[0:5]
	s2 := big_array[6:]
	s3 := big_array[3:9]

	expectSame(s1, s2)
	expectDifferent(s1, s3)

	c1 := make(chan string, 10)
	c2 := make(chan string, 10)
	c3 := make(chan string, 10)

	for _, ch := range []string{"a", "b", "c"} {
		c1 <- ch
		c2 <- ch + ch
		c3 <- ch
	}
	close(c1)
	close(c2)
	// do not close c3

	expectSame(c1, []string{"a", "b", "c"})
	expectDifferent(c2, []string{"a", "b", "c"})
	expectSame(c3, []string{"a", "b", "c"})
}

func TestHasSameElementsInternals(t *testing.T) {
	mockT := new(testing.T)
	g := defaultGroup(mockT)

	strType := reflect.TypeOf("x")

	intArray := []int{1, 2, 3}
	intChan := make(chan int, 1)

	intArrayType := reflect.TypeOf(intArray)
	intSliceType := reflect.TypeOf(intArray[0:1])
	intChanType := reflect.TypeOf(intChan)

	strArray := []string{"a", "b", "c"}
	strChan := make(chan string, 1)
	var strSendChan chan<- string
	strSendChan = strChan

	strArrayType := reflect.TypeOf(strArray)
	strSliceType := reflect.TypeOf(strArray[0:1])
	strChanType := reflect.TypeOf(strChan)
	strSendChanType := reflect.TypeOf(strSendChan)

	acceptableCases := [][]reflect.Type{
		[]reflect.Type{intArrayType, intArrayType},
		[]reflect.Type{intSliceType, intArrayType},
		[]reflect.Type{intArrayType, intSliceType},
		[]reflect.Type{intSliceType, intSliceType},
		[]reflect.Type{intChanType, intArrayType},
		[]reflect.Type{intChanType, intSliceType},
		[]reflect.Type{strArrayType, strArrayType},
		[]reflect.Type{strSliceType, strArrayType},
		[]reflect.Type{strArrayType, strSliceType},
		[]reflect.Type{strSliceType, strSliceType},
		[]reflect.Type{strChanType, strArrayType},
		[]reflect.Type{strChanType, strSliceType},
	}

	unacceptableCases := [][]reflect.Type{
		[]reflect.Type{strType, strArrayType},
		[]reflect.Type{strType, strSliceType},
		[]reflect.Type{strArrayType, intArrayType},
		[]reflect.Type{intArrayType, strArrayType},
		[]reflect.Type{strChanType, strChanType},
		[]reflect.Type{strArrayType, strType},
		[]reflect.Type{strSendChanType, strChanType},
		[]reflect.Type{strSendChanType, strArrayType},
	}

	for i, testcase := range acceptableCases {
		gotType := testcase[0]
		wantType := testcase[1]
		if !g.checkContainerTypes(gotType, wantType) {
			t.Errorf(
				"expected '%v' and '%v' to be accepted, but was not (case %d)",
				gotType,
				wantType,
				i)
		}
	}

	for i, testcase := range unacceptableCases {
		gotType := testcase[0]
		wantType := testcase[1]
		if g.checkContainerTypes(gotType, wantType) {
			t.Errorf(
				"expected '%v' and '%v' to be rejected, but was not (case %d)",
				gotType,
				wantType,
				i)
		}
	}
}

func TestGroupPassing(t *testing.T) {
	mockT := new(mockT)

	Group("name", mockT, func(g *G) {
		g.AssertTrue(true)
		g.Group("sub-group", func(g *G) {
			g.AssertFalse(false)
		})
	})

	if len(mockT.log) != 0 {
		t.Errorf("Expected no testing.T operations, got: %v", mockT.log)
	}
}

func TestGroupFailing(t *testing.T) {
	mockT := new(mockT)

	Group("name", mockT, func(g *G) {
		g.AssertTrue(false)
		g.Group("sub-group", func(g *G) {
			g.AssertFalse(false)
		})
	})

	if len(mockT.log) != 1 || mockT.log[0].op != "Errorf" {
		t.Errorf("expected single Errorf op, got %+v", mockT.log)
	}

	expectedPrefix := "name: "
	if mockT.log[0].args[0:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("got '%s', expected prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestNestedGroupFailing(t *testing.T) {
	mockT := new(mockT)

	Group("main-group", mockT, func(g *G) {
		g.AssertTrue(true)
		g.Group("sub-group", func(g *G) {
			g.AssertTrue(false)
		})
	})

	if len(mockT.log) != 1 || mockT.log[0].op != "Errorf" {
		t.Errorf("expected single Errorf op, got %+v", mockT.log)
	}

	expectedPrefix := "main-group sub-group: "
	if mockT.log[0].args[0:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("got '%s', expected prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestUngrouped(t *testing.T) {
	mockT := new(mockT)

	True(mockT, false)

	if len(mockT.log) != 1 || mockT.log[0].op != "Errorf" {
		t.Errorf("expected single Errorf op, got %+v", mockT.log)
	}

	expectedPrefix := "got: (bool) false, want (bool) true"
	if mockT.log[0].args[0:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("got '%s', expected prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}
