package matcher

import (
	"fmt"
	"reflect"
)

func CaptureAny() *ValueCaptor {
	return &ValueCaptor{nil, nil}
}

func CaptureMatching(m Matcher) *ValueCaptor {
	return &ValueCaptor{m, nil}
}

func CaptureType(t reflect.Type) *ValueCaptor {
	return &ValueCaptor{IsOfType{t}, nil}
}

type ValueCaptor struct {
	mustMatch Matcher
	V         interface{}
}

var _ Matcher = &ValueCaptor{}

func (vc *ValueCaptor) Matches(x interface{}) bool {
	if vc.mustMatch != nil && !vc.mustMatch.Matches(x) {
		return false
	}

	vc.V = x
	return true
}

func (vc *ValueCaptor) String() string {
	return fmt.Sprintf("valueCaptor(mustMatch: %s)", vc.mustMatch)
}
