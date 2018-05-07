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

package strings

import (
	"fmt"
	"testing"
)

type testCase struct {
	input    interface{}
	expected string
}

type stringer struct {
	s string
}

func (x stringer) String() string { return x.s }

var _ fmt.Stringer = stringer{}
var _ fmt.Stringer = &stringer{}

func strp(s string) *string { return &s }

var testCases = []testCase{
	// string
	{"simple string", "`simple string`"},
	{"string with backquotes: `x`", "\"string with backquotes: `x`\""},
	{`string with quotes: "x"`, "`string with quotes: \"x\"`"},
	{"", "``"},

	// *string
	{strp("simple string"), "`simple string`"},
	{strp("string with backquotes: `x`"), "\"string with backquotes: `x`\""},
	{strp(`string with quotes: "x"`), "`string with quotes: \"x\"`"},
	{strp(""), "``"},
	{(*string)(nil), "<nil>"},

	// fmt.Stringer
	{stringer{"simple string"}, "`simple string`"},
	{stringer{"string with backquotes: `x`"}, "\"string with backquotes: `x`\""},
	{stringer{`string with quotes: "x"`}, "`string with quotes: \"x\"`"},
	{stringer{""}, "``"},
	{&stringer{"simple string"}, "`simple string`"},
	{&stringer{"string with backquotes: `x`"}, "\"string with backquotes: `x`\""},
	{&stringer{`string with quotes: "x"`}, "`string with quotes: \"x\"`"},
	{&stringer{""}, "``"},
	{(*stringer)(nil), "<nil>"},

	// other types
	{123, "123"},
	{12.3, "12.3"},
	{true, "true"},
	{struct{ i int }{123}, "{i:123}"},
	{nil, "<nil>"},
}

func TestStringify(t *testing.T) {
	for i, tc := range testCases {
		got := Stringify(tc.input)

		if got != tc.expected {
			t.Errorf("testCases[%d]: expected %s, got %s", i, tc.expected, got)
		}
	}
}
