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

package matcher

import (
	"fmt"
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
	assert.Equal(t, cap.String(), "ValueCaptor(mustMatch: <nil>)")
	assert.True(t, cap.Matches(1234))

	asInt, ok := cap.V.(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, 1234)
}

func TestCaptureAll(t *testing.T) {
	cap := CaptureAll()
	assert.Equal(t, cap.String(), "AllValueCaptor(mustMatch: <nil>)")
	assert.True(t, cap.Matches(1234))
	assert.True(t, cap.Matches(5678))

	assert.Equal(t, len(cap.V), 2)

	asInt, ok := cap.V[0].(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, 1234)

	asInt, ok = cap.V[1].(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, 5678)
}

type maybeMatcher struct {
	shouldMatch bool
}

func (mm maybeMatcher) Matches(_ interface{}) bool {
	return mm.shouldMatch
}

func (mm maybeMatcher) String() string {
	return fmt.Sprintf("maybe=%t", mm.shouldMatch)
}

func TestCaptureMatching(t *testing.T) {
	want := teststruct{1234, "aoeu"}
	cap := CaptureMatching(maybeMatcher{true})

	assert.Equal(t, cap.String(), "ValueCaptor(mustMatch: maybe=true)")
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

type oddMatcher struct{}

func (om oddMatcher) Matches(x interface{}) bool {
	if i, ok := x.(int); ok {
		return i&1 == 1
	}

	return false
}

func (om oddMatcher) String() string {
	return "odd"
}

func TestCaptureAllMatching(t *testing.T) {
	cap := CaptureAllMatching(oddMatcher{})

	assert.Equal(t, cap.String(), "AllValueCaptor(mustMatch: odd)")
	assert.True(t, cap.Matches(123))
	assert.False(t, cap.Matches(456))
	assert.True(t, cap.Matches(789))

	assert.Equal(t, len(cap.V), 2)

	asInt, ok := cap.V[0].(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, 123)

	asInt, ok = cap.V[1].(int)
	assert.True(t, ok)
	assert.Equal(t, asInt, 789)
}

func TestCaptureAllMatchingDidNotMatch(t *testing.T) {
	cap := CaptureMatching(oddMatcher{})

	assert.False(t, cap.Matches(2))
	assert.False(t, cap.Matches(4))
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
