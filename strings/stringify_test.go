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

func (x *stringer) String() string { return x.s }

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
