package check

import (
	"fmt"
	"math"
	"testing"
)

type ewTestCase struct {
	name     string
	a        float64
	b        float64
	epsilon  float64
	expected bool
}

func (tc *ewTestCase) test(t *testing.T) {
	result := EqualWithin(tc.a, tc.b, tc.epsilon)
	if result != tc.expected {
		t.Errorf(
			"testcase %s: expected EqualWithin(%g, %g, %g) == %t, got %t",
			tc.name,
			tc.a,
			tc.b,
			tc.epsilon,
			tc.expected,
			result,
		)
	}
}

var ewTestCases = []ewTestCase{
	{"zeroes", 0.0, 0.0, math.SmallestNonzeroFloat64, true},
	{"ones", 1.0, 1.0, math.SmallestNonzeroFloat64, true},
	{"min diff", 0.0, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
	{"min eq", math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
	{"max eq", math.MaxFloat64, math.MaxFloat64, math.SmallestNonzeroFloat64, true},
	{"-max eq", -math.MaxFloat64, -math.MaxFloat64, math.SmallestNonzeroFloat64, true},
	{"biggest", math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64, false},
	{"biggest-", -math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, false},
	{"simple miss", 100.0, 101.0, 0.9, false},
	{"nans", math.NaN(), math.NaN(), 0.0, true},
	{"got nan", math.NaN(), 0.0, 0.0, false},
	{"want nan", 0.0, math.NaN(), 0.0, false},
	{"+infs", math.Inf(1), math.Inf(1), 0.0, true},
	{"-infs", math.Inf(-1), math.Inf(-1), 0.0, true},
	{"+/-inf", math.Inf(1), math.Inf(-1), math.MaxFloat64, false},
	{"-/+inf", math.Inf(-1), math.Inf(1), math.MaxFloat64, false},
	{"got +inf", math.Inf(1), math.MaxFloat64, math.MaxFloat64, false},
	{"got -inf", math.Inf(-1), -math.MaxFloat64, math.MaxFloat64, false},
	{"want +inf", math.MaxFloat64, math.Inf(1), math.MaxFloat64, false},
	{"want -inf", -math.MaxFloat64, math.Inf(-1), math.MaxFloat64, false},
}

func TestEqualWithin(t *testing.T) {
	for _, tc := range ewTestCases {
		tc.test(t)
	}
}

func TestEqualWithinForRoundingErrors(t *testing.T) {
	v := float64(0.0)
	roundingErrors := 0

	for i := 1; i <= 10; i++ {
		tenth := float64(i) / 10.0

		v += 0.1

		exact := true
		if v != tenth {
			roundingErrors++
			exact = false

		}

		tc := &ewTestCase{
			fmt.Sprintf("exact tenths: %g (%d)", v, i),
			v,
			tenth,
			0.0,
			exact,
		}
		tc.test(t)

		tc = &ewTestCase{
			fmt.Sprintf("inexact tenths: %g (%d)", v, i),
			v,
			tenth,
			0.0001,
			true,
		}
		tc.test(t)
	}

	if roundingErrors == 0 {
		t.Error("expected at least one floating-point rounding error, got none")
	}
}
