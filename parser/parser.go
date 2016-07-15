// Package parser provides a test output parser for go test.
package parser

import (
	"bytes"
	"time"

	"github.com/turbinelabs/test/testrunner/results"
)

type Parser struct {
	// Modifies command line arguments to include any test executable flags required by
	// the parser.
	FlagFn func([]string) []string

	// Parses the test executable's output, returning 1 or more test package results or
	// an error if the output could not be parsed.
	ParseFn func(
		packageName string,
		duration time.Duration,
		testOutput *bytes.Buffer,
	) ([]*results.TestPackage, error)
}

var (
	GoLangParser = Parser{ForceVerboseFlag, ParseTestOutput}
)
