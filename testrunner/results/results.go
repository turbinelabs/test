// Package results contains types that represent test packages, tests,
// and their results.
package results

import (
	"bytes"
)

// passed/skipped/failed
type TestResult int

const (
	Passed TestResult = iota
	Failed
	Skipped
)

type TestPackage struct {
	Name     string
	Result   TestResult
	Duration float64
	Tests    []*Test
	Output   string
}

type Test struct {
	Name     string
	Result   TestResult
	Duration float64
	Failure  bytes.Buffer
	Output   bytes.Buffer
}
