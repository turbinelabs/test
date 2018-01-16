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
package check

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

type testCase struct {
	got, want    interface{}
	expectEqual  bool
	expectReason string
}

func (tc testCase) runTest(t *testing.T, idx int) {
	ok, reason := DeepEqual(tc.got, tc.want)
	if ok != tc.expectEqual {
		s := "expected equal, got not equal"
		if !tc.expectEqual {
			s = "expected not equal, got equal"
		}
		t.Errorf(
			"testCases[%d]: %#v (%T) <=> %#v (%T), %s",
			idx,
			tc.got,
			tc.got,
			tc.want,
			tc.want,
			s,
		)
		return
	}

	if ok {
		if reason != "" {
			t.Errorf("testCases[%d]: objects equal, got unexpected reason %q", idx, reason)
		}
		return
	} else if reason != tc.expectReason {
		t.Errorf(
			"testCases[%d]: objects not equal for wrong reason\ngot reason: %s\nexp reason: %s",
			idx,
			reason,
			tc.expectReason,
		)
	}
}

func runTests(t *testing.T, testCases []testCase) {
	for i, tc := range testCases {
		tc.runTest(t, i)
	}
}

func TestDeepEqualSimpleTypes(t *testing.T) {
	runTests(
		t,
		[]testCase{
			// int
			{
				got:         1,
				want:        1,
				expectEqual: true,
			},
			{
				got:          1,
				want:         0,
				expectEqual:  false,
				expectReason: "got is 1, want is 0",
			},
			// uint
			{
				got:         uint(1),
				want:        uint(1),
				expectEqual: true,
			},
			{
				got:          uint(2),
				want:         uint(1),
				expectEqual:  false,
				expectReason: "got is 2, want is 1",
			},
			// float
			{
				got:         1.2,
				want:        1.2,
				expectEqual: true,
			},
			{
				got:          1.2,
				want:         2.4,
				expectEqual:  false,
				expectReason: "got is 1.2, want is 2.4",
			},
			// bool
			{
				got:         true,
				want:        true,
				expectEqual: true,
			},
			{
				got:          true,
				want:         false,
				expectEqual:  false,
				expectReason: "got is true, want is false",
			},
			// string
			{
				got:         "yes",
				want:        "yes",
				expectEqual: true,
			},
			{
				got:          "no",
				want:         "yes",
				expectEqual:  false,
				expectReason: `got is "no", want is "yes"`,
			},
			// complex
			{
				got:         complex(2, -2),
				want:        complex(2, -2),
				expectEqual: true,
			},
			{
				got:          complex(2, 2),
				want:         complex(2, -2),
				expectEqual:  false,
				expectReason: `got is (2+2i), want is (2-2i)`,
			},
		},
	)
}

func TestDeepEqualChanFundAndUnsafePointer(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)
	var nilC chan string

	f1 := func() string { return "x" }
	f2 := func() string { return "x" }
	var nilF func() string

	i1 := 100
	i2 := 100

	p1 := unsafe.Pointer(&i1)
	p2 := unsafe.Pointer(&i2)
	var nilP unsafe.Pointer

	runTests(
		t,
		[]testCase{
			// channels
			{
				got:         c1,
				want:        c1,
				expectEqual: true,
			},
			{
				got:         nilC,
				want:        nilC,
				expectEqual: true,
			},
			{
				got:          c1,
				want:         c2,
				expectEqual:  false,
				expectReason: fmt.Sprintf("got is chan %p, want is chan %p", c1, c2),
			},

			// funcs
			{
				got:         f1,
				want:        f1,
				expectEqual: true,
			},
			{
				got:         nilF,
				want:        nilF,
				expectEqual: true,
			},
			{
				got:          f1,
				want:         f2,
				expectEqual:  false,
				expectReason: fmt.Sprintf("got is func %p, want is func %p", f1, f2),
			},

			// pointers
			{
				got:         p1,
				want:        p1,
				expectEqual: true,
			},
			{
				got:         nilP,
				want:        nilP,
				expectEqual: true,
			},
			{
				got:         p1,
				want:        p2,
				expectEqual: false,
				expectReason: fmt.Sprintf(
					"got is unsafe.Pointer %p, want is unsafe.Pointer %p", p1, p2),
			},
		},
	)
}

func TestDeepEqualMaps(t *testing.T) {
	var nil1, nil2 map[string]string

	m1 := map[string]string{"x": "y"}
	m2 := map[string]string{"x": "y"}
	m3 := map[string]string{"not x": "y"}
	m4 := map[string]string{"x": "not y"}
	m5 := map[string]string{"x": "y", "extra": "y"}

	runTests(
		t,
		[]testCase{
			{
				got:         nil1,
				want:        nil2,
				expectEqual: true,
			},
			{
				got:          m1,
				want:         nil1,
				expectEqual:  false,
				expectReason: "got is not nil, want is nil",
			},
			{
				got:          nil1,
				want:         m1,
				expectEqual:  false,
				expectReason: "got is nil, want is not nil",
			},
			// equal by pointer
			{
				got:         m1,
				want:        m1,
				expectEqual: true,
			},
			// equal by value
			{
				got:         m1,
				want:        m2,
				expectEqual: true,
			},
			// not equal: missing key
			{
				got:          m1,
				want:         m3,
				expectEqual:  false,
				expectReason: `got["x"] is "y", want["x"] is missing`,
			},
			// not equal: different value
			{
				got:          m1,
				want:         m4,
				expectEqual:  false,
				expectReason: `got["x"] is "y", want["x"] is "not y"`,
			},
			// not equal: different size
			{
				got:          m1,
				want:         m5,
				expectEqual:  false,
				expectReason: "got is a map[string]string with 1 entries, want has 2 entries",
			},
		},
	)
}

func TestDeepEqualPointers(t *testing.T) {
	i1 := 100
	i2 := 100
	i3 := 200
	var nilI *int

	runTests(
		t,
		[]testCase{
			// equal by value
			{
				got:         &i1,
				want:        &i2,
				expectEqual: true,
			},
			// equal by pointer
			{
				got:         &i1,
				want:        &i1,
				expectEqual: true,
			},
			{
				got:          &i1,
				want:         &i3,
				expectEqual:  false,
				expectReason: "got is 100, want is 200",
			},
			{
				got:         nilI,
				want:        nilI,
				expectEqual: true,
			},
			{
				got:          &i1,
				want:         nilI,
				expectEqual:  false,
				expectReason: "got is not nil, want is nil",
			},
		},
	)
}

type testStruct struct {
	i      int
	child  *testStruct
	parent *testStruct
}

func TestDeepEqualStruct(t *testing.T) {
	s1 := testStruct{i: 1}
	s2 := testStruct{i: 1}
	s3 := testStruct{i: 2}

	s4 := testStruct{i: 1}
	nested4 := testStruct{i: 100, parent: &s4}
	s4.child = &nested4

	s5 := testStruct{i: 1}
	nested5 := testStruct{i: 100, parent: &s5}
	s5.child = &nested5

	s6 := testStruct{i: 1}
	nested6 := testStruct{i: 200, parent: &s6}
	s6.child = &nested6

	runTests(
		t,
		[]testCase{
			{
				got:         s1,
				want:        s1,
				expectEqual: true,
			},
			{
				got:         s1,
				want:        s2,
				expectEqual: true,
			},
			{
				got:          s1,
				want:         s3,
				expectEqual:  false,
				expectReason: "got.i is 1, want.i is 2",
			},
			{
				got:         s4,
				want:        s4,
				expectEqual: true,
			},
			{
				got:         s4,
				want:        s5,
				expectEqual: true,
			},
			{
				got:          s4,
				want:         s6,
				expectEqual:  false,
				expectReason: "got.child.i is 100, want.child.i is 200",
			},
		},
	)
}

type testIface interface {
	F() string
}

type testIfaceStrImpl string

func (s testIfaceStrImpl) F() string { return "F:" + string(s) }

type testIfaceStructImpl struct {
	fvalue string
}

func (s testIfaceStructImpl) F() string { return s.fvalue }

var _ testIface = testIfaceStrImpl("")
var _ testIface = testIfaceStructImpl{}

func TestDeepEqualInterfaces(t *testing.T) {
	var i1 testIface = testIfaceStrImpl("string")
	var i2 testIface = testIfaceStrImpl("string")
	var i3 testIface = testIfaceStrImpl("string2")
	var i4 testIface = testIfaceStructImpl{fvalue: "struct"}
	var i5 testIface = testIfaceStructImpl{fvalue: "struct"}
	var i6 testIface = testIfaceStructImpl{fvalue: "struct2"}

	// Use slices because this forces the array element Kind to be reflect.Interface where
	// passing the interfaces directly does not.
	runTests(
		t,
		[]testCase{
			{
				got:         []testIface{i1},
				want:        []testIface{i2},
				expectEqual: true,
			},
			{
				got:          []testIface{i1},
				want:         []testIface{i3},
				expectEqual:  false,
				expectReason: `got[0].(check.testIfaceStrImpl) is "string", want[0].(check.testIfaceStrImpl) is "string2"`,
			},
			{
				got:          []testIface{i1},
				want:         []testIface{i4},
				expectEqual:  false,
				expectReason: "got[0].(check.testIfaceStrImpl) has type check.testIfaceStrImpl, want[0].(check.testIfaceStrImpl) has type check.testIfaceStructImpl",
			},
			{
				got:          []testIface{i1},
				want:         []testIface{nil},
				expectEqual:  false,
				expectReason: "got[0] is not nil, want[0] is nil",
			},
			{
				got:         []testIface{i4},
				want:        []testIface{i5},
				expectEqual: true,
			},
			{
				got:          []testIface{i4},
				want:         []testIface{i6},
				expectEqual:  false,
				expectReason: `got[0].(check.testIfaceStructImpl).fvalue is "struct", want[0].(check.testIfaceStructImpl).fvalue is "struct2"`,
			},
		},
	)
}

func TestDeepEqualArrays(t *testing.T) {
	a1 := [2]string{"a", "b"}
	a2 := [2]string{"a", "b"}
	a3 := [2]string{"a", "not b"}
	a4 := [1]string{"a"}

	runTests(
		t,
		[]testCase{
			{
				got:         a1,
				want:        a2,
				expectEqual: true,
			},
			{
				got:          a1,
				want:         a3,
				expectEqual:  false,
				expectReason: `got[1] is "b", want[1] is "not b"`,
			},

			// nest arrays in a struct to bypass checks in DeepEqual
			{
				got:          struct{ a interface{} }{a1},
				want:         struct{ a interface{} }{a4},
				expectEqual:  false,
				expectReason: "got.a.([2]string) has type [2]string, want.a.([2]string) has type [1]string",
			},
		},
	)
}

func TestDeepEqualSlices(t *testing.T) {
	s1 := []string{"a", "b"}
	s2 := []string{"a", "b"}
	s3 := []string{"a", "not b"}
	s4 := []string{"a"}
	var s5 []string
	s6 := []string{"a", "b", "c", "d"}
	s7 := []string{"x", "b", "y", "d", "z"}
	s8 := []string{"x", "y", "z", "x", "y", "not z"}

	runTests(
		t,
		[]testCase{
			// equal by value
			{
				got:         s1,
				want:        s2,
				expectEqual: true,
			},
			// equal by pointer
			{
				got:         s1,
				want:        s1,
				expectEqual: true,
			},
			{
				got:          s1,
				want:         s3,
				expectEqual:  false,
				expectReason: `got[1] is "b", want[1] is "not b"`,
			},
			{
				got:          s1,
				want:         s4,
				expectEqual:  false,
				expectReason: `got[1] is "b", no want[1] given`,
			},
			{
				got:          s4,
				want:         s1,
				expectEqual:  false,
				expectReason: `got[1] missing, want[1] is "b"`,
			},
			{
				got:          s1,
				want:         s5,
				expectEqual:  false,
				expectReason: "got is not nil, want is nil",
			},
			{
				got:         s6,
				want:        s7,
				expectEqual: false,
				expectReason: strings.Join([]string{
					`got[0] is "a", want[0] is "x"`,
					`got[2] is "c", want[2] is "y"`,
					`got[4] missing, want[4] is "z"`,
				}, "\n"),
			},
			{
				got:         s7,
				want:        s6,
				expectEqual: false,
				expectReason: strings.Join([]string{
					`got[0] is "x", want[0] is "a"`,
					`got[2] is "y", want[2] is "c"`,
					`got[4] is "z", no want[4] given`,
				}, "\n"),
			},

			// Verify that we don't mistakenly detect subslices as equal by pointer
			{
				got:          s6,
				want:         s6[0:3],
				expectEqual:  false,
				expectReason: `got[3] is "d", no want[3] given`,
			},
			{
				got:          s8[0:3],
				want:         s8[3:6],
				expectEqual:  false,
				expectReason: `got[2] is "z", want[2] is "not z"`,
			},
			{
				got:         s8[0:3],
				want:        s8[0:3],
				expectEqual: true,
			},
		},
	)
}

func TestDeepEqualShortCircuits(t *testing.T) {
	runTests(
		t,
		[]testCase{
			{
				got:          nil,
				want:         "x",
				expectEqual:  false,
				expectReason: "got is nil, but want is non-nil",
			},
			{
				got:          "x",
				want:         nil,
				expectEqual:  false,
				expectReason: "got is non-nil, but want is nil",
			},
			{
				got:         nil,
				want:        nil,
				expectEqual: true,
			},
			{
				got:          100,
				want:         100.0,
				expectEqual:  false,
				expectReason: "got is of type int, want is of type float64",
			},
		},
	)
}

func TestDeepEqualInternalsInvalidValues(t *testing.T) {
	var p1, p2 *int
	v1 := reflect.ValueOf(p1).Elem()
	if v1.IsValid() {
		t.Error("expected invalid value for v1")
	}
	i := 1
	p2 = &i
	v2 := reflect.ValueOf(p2).Elem()
	if !v2.IsValid() {
		t.Error("expected valid value for v2")
	}

	ok, reason := deepEqual(v1, v2, nil, "")
	if ok {
		t.Error("expected not equal result, but got equal")
	}
	expectedReason := "got is not valid, want is valid"
	if reason != expectedReason {
		t.Errorf("expected reason to be: %q, but it was %q", expectedReason, reason)
	}
}

func TestDeepEqualInternalsTrackedTypes(t *testing.T) {
	s1 := &testStruct{i: 1}
	s2 := &testStruct{i: 2}

	v1 := reflect.Indirect(reflect.ValueOf(&s1))
	v2 := reflect.Indirect(reflect.ValueOf(&s2))

	visited := map[visit]struct{}{}
	ok, _ := deepEqual(v1, v2, visited, "")
	if ok {
		t.Error("expected not equal result, but got equal result")
	}

	// Invoking deepEqual again trigger the reference cycle breaking
	// behavior deepEqual. Reverse the order to prove order doesn't
	// matter.
	ok, _ = deepEqual(v2, v1, visited, "")
	if !ok {
		t.Error("expected equal result, but got not equal result")
	}
}

func TestRender(t *testing.T) {
	a := "abc"
	b := "def"
	var c *string
	d := struct{ i int }{i: 123}

	expected := fmt.Sprintf("%q", a)
	rendered := render(reflect.TypeOf(a), reflect.ValueOf(a))
	if rendered != expected {
		t.Errorf("expected %s, got %s", expected, rendered)
	}

	expected = fmt.Sprintf("%q", b)
	rendered = render(reflect.TypeOf(&b), reflect.ValueOf(&b))
	if rendered != expected {
		t.Errorf("expected %s, got %s", expected, rendered)
	}

	expected = "<nil>"
	rendered = render(reflect.TypeOf(c), reflect.ValueOf(c))
	if rendered != expected {
		t.Errorf("expected %s, got %s", expected, rendered)
	}

	expected = "{i:123}"
	rendered = render(reflect.TypeOf(d), reflect.ValueOf(d))
	if rendered != expected {
		t.Errorf("expected %s, got %s", expected, rendered)
	}

	expected = "&{i:123}"
	rendered = render(reflect.TypeOf(&d), reflect.ValueOf(&d))
	if rendered != expected {
		t.Errorf("expected %s, got %s", expected, rendered)
	}
}
