package main

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// passed/skipped/failed
type testResult int

const (
	passed testResult = iota
	failed
	skipped
)

type testPackage struct {
	name     string
	result   testResult
	duration float64
	tests    []*test
	output   string
}

type test struct {
	name     string
	result   testResult
	duration float64
	failure  bytes.Buffer
	output   bytes.Buffer
}

var (
	resultRegex  = regexp.MustCompile(`^--- (PASS|FAIL|SKIP): (.+) \((\d+\.\d+)(?: seconds|s)\)$`)
	summaryRegex = regexp.MustCompile(`^(?:PASS|FAIL)$`)
)

func parseTestOutput(pkgName string, duration time.Duration, output *bytes.Buffer) *testPackage {
	eof := false
	var t *test
	testPkg := testPackage{
		name:     pkgName,
		result:   skipped,
		duration: duration.Seconds(),
		tests:    make([]*test, 0),
		output:   string(output.Bytes()),
	}

	for !eof {
		lineBytes, err := output.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				panic(err)
			}
		}

		line := strings.TrimRightFunc(string(lineBytes), unicode.IsSpace)

		if t == nil {
			if strings.HasPrefix(line, "=== RUN ") {
				// start of test
				t = new(test)
				t.name = strings.TrimSpace(line[8:])
			} else if m := summaryRegex.FindStringSubmatch(line); len(m) == 1 {
				// End of package
				if testPkg.result != skipped {
					panic("expected only a single package")
				}
				switch line {
				case "PASS":
					testPkg.result = passed
				default:
					testPkg.result = failed
				}
			} else if len(testPkg.tests) > 0 && strings.HasPrefix(line, "\t") {
				// test failure output
				lastTest := testPkg.tests[len(testPkg.tests)-1]
				lastTest.failure.Write(lineBytes)
			}
		} else {
			if m := resultRegex.FindStringSubmatch(line); len(m) == 4 {
				// end of test
				switch m[1] {
				case "PASS":
					t.result = passed
				case "SKIP":
					t.result = skipped
				default:
					t.result = failed
				}
				t.duration = parseDuration(m[3])

				testPkg.tests = append(testPkg.tests, t)
				t = nil
			} else {
				t.output.Write(lineBytes)
			}
		}
	}

	if testPkg.result == skipped {
		testPkg.result = failed
		testPkg.output += "\n[Did not find package result: marking package as failed.]\n"
	}

	return &testPkg
}

func parseDuration(d string) float64 {
	f, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0.0
	}
	return f
}
