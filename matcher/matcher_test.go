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
	"bytes"
	"reflect"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestAnyWriter(t *testing.T) {
	aw := AnyWriter{[]byte("yep")}
	buf := &bytes.Buffer{}

	assert.False(t, aw.Matches("nope"))

	assert.True(t, aw.Matches(buf))
	assert.Equal(t, buf.String(), "yep")

	assert.Equal(t, aw.String(), `AnyWriter("yep")`)
}

func TestPredicateMatcher(t *testing.T) {
	eqm := PredicateMatcher{
		Name: "string value tester",
		Test: func(x interface{}) bool {
			v, ok := x.(string)
			if !ok {
				return false
			}

			return v == "matched!"
		},
	}

	assert.False(t, eqm.Matches("nope"))

	assert.True(t, eqm.Matches("matched!"))
	assert.Equal(t, eqm.String(), "PredicateMatcher(string value tester)")
}

func TestIsOfTypeMatcher(t *testing.T) {
	iot := IsOfType{reflect.TypeOf(teststruct{})}
	assert.False(t, iot.Matches("nope"))
	assert.True(t, iot.Matches(teststruct{1234, "ok"}))
	assert.Equal(t, iot.String(), "IsOfType(matcher.teststruct)")
}

func TestSameElementsMatcher(t *testing.T) {
	se := SameElements{[]int{1, 2, 3}}

	assert.False(t, se.Matches("nope"))
	assert.False(t, se.Matches([]string{"1", "2", "3"}))
	assert.True(t, se.Matches([]int{1, 2, 3}))
	assert.True(t, se.Matches([]int{3, 1, 2}))
}
