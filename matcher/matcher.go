// Package matcher provides helpful implementations of the gomock.Matcher interface.
package matcher

import (
	"fmt"
	"io"
	"reflect"
)

// Matcher duplicates gomock.Matcher interface. It is not referenced directly so
// as to not create a dependency between the test parent package and gomock.
type Matcher interface {
	Matches(interface{}) bool
	String() string
}

// AnyWriter implements gomock.Matcher, and can be used
// to mock a call that passes an io.Writer, filling that
// writer with the provided byte array
type AnyWriter struct {
	Data []byte
}

var _ Matcher = &AnyWriter{}

func (a AnyWriter) Matches(x interface{}) bool {
	if writer, ok := x.(io.Writer); ok {
		writer.Write(a.Data)
		return true
	}
	return false
}

func (a AnyWriter) String() string {
	return fmt.Sprintf("AnyWriter(%q)", a.Data)
}

type PredicateMatcher struct {
	Test func(interface{}) bool
	Name string
}

var _ Matcher = &PredicateMatcher{}

func (em PredicateMatcher) Matches(x interface{}) bool {
	return em.Test(x)
}

func (em PredicateMatcher) String() string {
	return fmt.Sprintf("PredicateMatcher(%s)", em.Name)
}

type IsOfType struct {
	Type reflect.Type
}

var _ Matcher = &IsOfType{}

func (iot IsOfType) Matches(x interface{}) bool {
	return reflect.TypeOf(x) == iot.Type
}

func (iot IsOfType) String() string {
	return fmt.Sprintf("IsOfType(%s)", iot.Type.String())
}
