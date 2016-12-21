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
