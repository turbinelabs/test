package matcher

import (
	"testing"
)

type teststruct struct {
	a int
	b string
}

func TestCaptureAny(t *testing.T) {
	cap := CaptureAny()
	passed := 1234
	if !cap.Matches(1234) {
		t.Errorf("got: does not match, want: matches")
	}

	asInt, ok := cap.V.(int)
	if !ok {
		t.Errorf("got: %T, want: int", cap.V)
	}

	if asInt != passed {
		t.Errorf("got: %d, want: %d", asInt, passed)
	}
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

	if !cap.Matches(want) {
		t.Errorf("got: does not match, want: matches")
	}

	asTs, ok := cap.V.(teststruct)
	if !ok {
		t.Errorf("got: %T, want: teststruct", cap.V)
	}

	if asTs != want {
		t.Errorf("got: %#v, want: %#v\n", cap.V, want)
	}
}

func TestCaptureMatchingDidNotMatch(t *testing.T) {
	doNotWant := teststruct{1234, "aoeu"}
	cap := CaptureMatching(maybeMatcher{false})

	if cap.Matches(doNotWant) {
		t.Errorf("got: matches, want: no match")
	}
	if cap.V != nil {
		t.Errorf("got: %#v, want: nil\n", cap.V)
	}
}
