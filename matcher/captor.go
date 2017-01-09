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

// CaptureAny will match and capture any value.
func CaptureAny() *ValueCaptor {
	return &ValueCaptor{nil, nil}
}

// CaptureMatching is a ValueCaptor that proxies the given Matcher to match
// and capture matching values.
func CaptureMatching(m Matcher) *ValueCaptor {
	return &ValueCaptor{m, nil}
}

// CaptureType is a ValueCaptor that uses an IsOfType matcher to capture
// values of the given reflect.Type.
func CaptureType(t reflect.Type) *ValueCaptor {
	return &ValueCaptor{IsOfType{t}, nil}
}

// A ValueCaptor can be wrapped around a Matcher to produce a Matcher
// which also captures the most recent matched value.
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
