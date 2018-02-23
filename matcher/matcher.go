/*
Copyright 2018 Turbine Labs, Inc.

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

// Package matcher provides helpful implementations of the gomock.Matcher
// interface.
package matcher

import (
	"fmt"
	"io"
	"reflect"

	"github.com/turbinelabs/test/check"
)

// Matcher duplicates gomock.Matcher interface. It is not referenced directly
// so as to not create a dependency between the test parent package and gomock.
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

// PredicateMatcher implements gomock.Matcher and can be used to mock a
// call that receives any value that passes a specified predicate.
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

// IsOfType implements gomock.Matcher and can be used to mock a
// call that receives any value of a specified reflect.Type.
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

// SameElements implements gomock.Matcher and can be used to mock a
// call that receives an array-like parameter where the order of
// elements may vary. Array-like parameters are arrays, slices, or
// channels. The Elems field is expected to be an array or slice. See
// check.HasSameElements for details on the comparison used.
type SameElements struct {
	Elems interface{}
}

func (se SameElements) Matches(x interface{}) bool {
	return check.HasSameElements(x, se.Elems) == nil
}

func (se SameElements) String() string {
	return fmt.Sprintf("SameElements(%+v)", se.Elems)
}
