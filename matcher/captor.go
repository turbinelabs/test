package matcher

import (
	"fmt"
)

func CaptureAny() *ValueCaptor {
	return &ValueCaptor{nil, nil}
}

func CaptureMatching(m Matcher) *ValueCaptor {
	return &ValueCaptor{m, nil}
}

// TODO: write CaptureType(interface{}) that will capture anything matching the
// the provided type. Construct a matcher.IsTypeOf and use that.

type ValueCaptor struct {
	mustMatch Matcher
	V         interface{}
}

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

// Matcher duplicates gomock.Matcher interface. It is not referenced directly so
// as to not create a dependency between the test parent package and gomock.
type Matcher interface {
	Matches(interface{}) bool
	String() string
}
