package matcher

import (
	"reflect"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type teststruct struct {
	a int
	b string
}

func TestCaptureAny(t *testing.T) {
	cap := CaptureAny()
	passed := 1234
	assert.True(t, cap.Matches(1234))

	asInt, ok := cap.V.(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, passed)
}

type maybeMatcher struct {
	shouldMatch bool
}

func (mm maybeMatcher) Matches(_ interface{}) bool {
	return mm.shouldMatch
}

func (mm maybeMatcher) String() string {
	return ""
}

func TestCaptureMatching(t *testing.T) {
	want := teststruct{1234, "aoeu"}
	cap := CaptureMatching(maybeMatcher{true})

	assert.True(t, cap.Matches(want))

	asTs, ok := cap.V.(teststruct)
	assert.True(t, ok)
	assert.Equal(t, asTs, want)
}

func TestCaptureMatchingDidNotMatch(t *testing.T) {
	doNotWant := teststruct{1234, "aoeu"}
	cap := CaptureMatching(maybeMatcher{false})

	assert.False(t, cap.Matches(doNotWant))
	assert.Nil(t, cap.V)
}

func TestCaptureType(t *testing.T) {
	rightType := teststruct{1234, "whee"}
	wrongType := struct{ x string }{x: "nope"}

	cap := CaptureType(reflect.TypeOf(teststruct{}))
	assert.True(t, cap.Matches(rightType))
	assert.DeepEqual(t, cap.V, rightType)

	cap = CaptureType(reflect.TypeOf(teststruct{}))
	assert.False(t, cap.Matches(wrongType))
	assert.Nil(t, cap.V)
}
