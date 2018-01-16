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

package check

import "testing"

type foo struct {
	Field1 *int    `json:"field1"`
	Field2 *string `json:"field2"`
}

func testIsNil(t *testing.T, i interface{}, expected bool) {
	got := IsNil(i)
	if got != expected {
		t.Errorf("IsNil(%+v) == %t, expected %t", i, got, expected)
	}
}

func TestIsNil(t *testing.T) {
	testIsNil(t, nil, true)

	var i *int
	testIsNil(t, 1, false)
	testIsNil(t, i, true)
	x := 2
	i = &x
	testIsNil(t, i, false)

	// c.f. https://github.com/turbinelabs/golang-gotchas/blob/master/example_nil_interfaces_test.go
	var f *foo
	testIsNil(t, f, true)
	f = &foo{}
	testIsNil(t, f, false)

	var f2 *foo
	var i2 interface{}
	i2 = f2
	testIsNil(t, i2, true)
	testIsNil(t, f2, true)
}
