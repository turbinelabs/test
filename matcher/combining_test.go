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
	"testing"

	"github.com/turbinelabs/test/assert"
)

type anyMatcher struct{}
type noneMatcher struct{}
type panickyMatcher struct{}

func (*anyMatcher) Matches(x interface{}) bool     { return true }
func (*anyMatcher) String() string                 { return "always matches" }
func (*noneMatcher) Matches(x interface{}) bool    { return false }
func (*noneMatcher) String() string                { return "never matches" }
func (*panickyMatcher) Matches(x interface{}) bool { panic("unexpected matcher check") }
func (*panickyMatcher) String() string             { return "panicky" }

func any() Matcher     { return &anyMatcher{} }
func none() Matcher    { return &noneMatcher{} }
func panicky() Matcher { return &panickyMatcher{} }

func TestAnd(t *testing.T) {
	assert.True(t, And(any(), any()).Matches(0))
	assert.False(t, And(any(), none()).Matches(0))
	assert.False(t, And(none(), any()).Matches(0))
	assert.False(t, And(none(), none()).Matches(0))

	defer func() {
		if x := recover(); x != nil {
			assert.Failed(t, fmt.Sprintf("test panic: %v", x))
		}
	}()

	// Test short-circuit
	assert.False(t, And(none(), panicky()).Matches(0))
}

func TestAndString(t *testing.T) {
	a := And(any(), none())
	assert.Equal(t, a.String(), "(always matches) and (never matches)")
}

func TestOr(t *testing.T) {
	assert.True(t, Or(any(), any()).Matches(0))
	assert.True(t, Or(any(), none()).Matches(0))
	assert.True(t, Or(none(), any()).Matches(0))
	assert.False(t, Or(none(), none()).Matches(0))

	defer func() {
		if x := recover(); x != nil {
			assert.Failed(t, fmt.Sprintf("test panic: %v", x))
		}
	}()

	// Test short-circuit
	assert.True(t, Or(any(), panicky()).Matches(0))
}

func TestOrString(t *testing.T) {
	o := Or(any(), none())
	assert.Equal(t, o.String(), "(always matches) or (never matches)")
}
