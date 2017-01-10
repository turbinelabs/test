package check

import (
	"math"
)

// Returns true if got is within epsilon of want. Useful for comparing
// the results of floating point arithmetic. If got and want are both
// NaN, +Inf, or -Inf this function returns true.
func EqualWithin(got, want, epsilon float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	} else if math.IsNaN(got) {
		return false
	}

	if math.IsInf(want, 1) {
		return math.IsInf(got, 1)
	} else if math.IsInf(got, 1) {
		return false
	}

	if math.IsInf(want, -1) {
		return math.IsInf(got, -1)
	} else if math.IsInf(got, -1) {
		return false
	}

	if got == want {
		return true
	}

	return math.Abs(got-want) <= math.Abs(epsilon)
}
