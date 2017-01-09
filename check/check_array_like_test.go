package check

import (
	"reflect"
	"testing"
)

type complexStruct struct {
	x int
	y *string
}

var (
	string1a = "string"
	string1b = "string"
	string2  = "other string"

	cs1a = complexStruct{1, &string1a}
	cs1b = complexStruct{1, &string1a}
	cs2a = complexStruct{1, &string1a}
	cs2b = complexStruct{1, &string1b}
	cs3  = complexStruct{1, &string1a}
	cs4  = complexStruct{1, &string2}
)

func TestHasSameElementsInternals(t *testing.T) {
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
		{intArrayType, intArrayType},
		{intSliceType, intArrayType},
		{intArrayType, intSliceType},
		{intSliceType, intSliceType},
		{intChanType, intArrayType},
		{intChanType, intSliceType},
		{strArrayType, strArrayType},
		{strSliceType, strArrayType},
		{strArrayType, strSliceType},
		{strSliceType, strSliceType},
		{strChanType, strArrayType},
		{strChanType, strSliceType},
	}

	unacceptableCases := [][]reflect.Type{
		{strType, strArrayType},
		{strType, strSliceType},
		{strArrayType, intArrayType},
		{intArrayType, strArrayType},
		{strChanType, strChanType},
		{strArrayType, strType},
		{strSendChanType, strChanType},
		{strSendChanType, strArrayType},
	}

	for i, testcase := range acceptableCases {
		gotType := testcase[0]
		wantType := testcase[1]
		if err := checkContainerTypes(gotType, wantType); err != nil {
			t.Errorf(
				"expected '%v' and '%v' to be accepted, but was not (case %d): %s",
				gotType,
				wantType,
				i,
				err.Error(),
			)
		}
	}

	for i, testcase := range unacceptableCases {
		gotType := testcase[0]
		wantType := testcase[1]
		if err := checkContainerTypes(gotType, wantType); err == nil {
			t.Errorf(
				"expected '%v' and '%v' to be rejected, but was not (case %d)",
				gotType,
				wantType,
				i)
		}
	}
}

func TestHasSameElements(t *testing.T) {
	expectSame := func(a, b interface{}) {
		if err := HasSameElements(a, b); err != nil {
			t.Errorf(
				"expected '%v' to have same elements as '%v': %s",
				a,
				b,
				err.Error(),
			)
		}
	}

	expectDifferent := func(a, b interface{}) {
		if HasSameElements(a, b) == nil {
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
