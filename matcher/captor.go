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

package matcher

import (
	"fmt"
	"reflect"
)

// CaptureAny will match any value and captures the most recent one.
func CaptureAny() *ValueCaptor {
	return &ValueCaptor{}
}

// CaptureAll will match and capture multiple values of any type.
func CaptureAll() *AllValueCaptor {
	return &AllValueCaptor{}
}

// CaptureMatching is a ValueCaptor that proxies the given Matcher to
// match and capture the most recent matching value.
func CaptureMatching(m Matcher) *ValueCaptor {
	return &ValueCaptor{mustMatch: m}
}

// CaptureAllMatching is an AllValueCaptor that proxies the given
// Matcher to match and capture all matching values.
func CaptureAllMatching(m Matcher) *AllValueCaptor {
	return &AllValueCaptor{mustMatch: m}
}

// CaptureType is a ValueCaptor that uses an IsOfType matcher to
// capture the most recent value of the given reflect.Type.
func CaptureType(t reflect.Type) *ValueCaptor {
	return &ValueCaptor{mustMatch: IsOfType{t}}
}

// A ValueCaptor can be wrapped around a Matcher to produce a Matcher
// which also captures the most recent matched value.
type ValueCaptor struct {
	mustMatch Matcher
	V         interface{}
}

var _ Matcher = &ValueCaptor{}

// Matches compares x to the underlying Matcher, if any. If x matches
// (or there is no underlying Matcher), it returns true and captures
// x. If a previous value was captured, it is forgotten. If x does not
// match, it returns false.
func (vc *ValueCaptor) Matches(x interface{}) bool {
	if vc.mustMatch != nil && !vc.mustMatch.Matches(x) {
		return false
	}

	vc.V = x
	return true
}

func (vc *ValueCaptor) String() string {
	return fmt.Sprintf("ValueCaptor(mustMatch: %v)", vc.mustMatch)
}

// An AllValueCaptor can be wrapped around a Matcher to produce a
// Matcher which also captures all matched values.
type AllValueCaptor struct {
	mustMatch Matcher
	V         []interface{}
}

// Matches compares x to the underlying Matcher, if any. If x matches
// (or there is no underlying Matcher), it returns true and captures
// x. If x does not match, it returns false.
func (avc *AllValueCaptor) Matches(x interface{}) bool {
	if avc.mustMatch != nil && !avc.mustMatch.Matches(x) {
		return false
	}

	if avc.V == nil {
		avc.V = make([]interface{}, 1)
		avc.V[0] = x
	} else {
		avc.V = append(avc.V, x)
	}

	return true
}

func (avc *AllValueCaptor) String() string {
	return fmt.Sprintf("AllValueCaptor(mustMatch: %v)", avc.mustMatch)
}
