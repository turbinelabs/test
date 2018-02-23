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
