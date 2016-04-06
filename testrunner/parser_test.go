package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/turbinelabs/tbn/test/assert"
)

const (
	testPackageName = "github.com/turbinelabs/tbn/agent/sdagentctl/helper"

	SinglePackageVerbose = `=== RUN   TestWrapAndCallAgentNoErr
--- PASS: TestWrapAndCallAgentNoErr (0.01s)
=== RUN   TestWrapAndCallAgentValidationFails
--- PASS: TestWrapAndCallAgentValidationFails (0.02s)
=== RUN   TestWrapAndCallAgentAgentFails
--- PASS: TestWrapAndCallAgentAgentFails (0.03s)
=== RUN   TestWrapAndCallAgentInvocationFails
--- PASS: TestWrapAndCallAgentInvocationFails (0.04s)
=== RUN   TestWrapAndCallAgentMarshalFails
--- PASS: TestWrapAndCallAgentMarshalFails (0.05s)
=== RUN   TestDiffNoErr
--- PASS: TestDiffNoErr (0.06s)
=== RUN   TestDiffDecodeFails
--- PASS: TestDiffDecodeFails (0.07s)
=== RUN   TestDiffReadAllFails
--- PASS: TestDiffReadAllFails (0.08s)
=== RUN   TestPatchNoErr
--- PASS: TestPatchNoErr (0.09s)
=== RUN   TestPatchAgentErr
--- PASS: TestPatchAgentErr (0.10s)
=== RUN   TestPatchDryRun
--- PASS: TestPatchDryRun (0.11s)
PASS
`

	SinglePackageVerboseCoverage = `=== RUN   TestWrapAndCallAgentNoErr
--- PASS: TestWrapAndCallAgentNoErr (0.01s)
=== RUN   TestWrapAndCallAgentValidationFails
--- PASS: TestWrapAndCallAgentValidationFails (0.02s)
=== RUN   TestWrapAndCallAgentAgentFails
--- PASS: TestWrapAndCallAgentAgentFails (0.03s)
=== RUN   TestWrapAndCallAgentInvocationFails
--- PASS: TestWrapAndCallAgentInvocationFails (0.04s)
=== RUN   TestWrapAndCallAgentMarshalFails
--- PASS: TestWrapAndCallAgentMarshalFails (0.05s)
=== RUN   TestDiffNoErr
--- PASS: TestDiffNoErr (0.06s)
=== RUN   TestDiffDecodeFails
--- PASS: TestDiffDecodeFails (0.07s)
=== RUN   TestDiffReadAllFails
--- PASS: TestDiffReadAllFails (0.08s)
=== RUN   TestPatchNoErr
--- PASS: TestPatchNoErr (0.09s)
=== RUN   TestPatchAgentErr
--- PASS: TestPatchAgentErr (0.10s)
=== RUN   TestPatchDryRun
--- PASS: TestPatchDryRun (0.11s)
PASS
coverage: 57.1% of statements
`

	SinglePackageVerboseCoverageOutput = `=== RUN   TestWrapAndCallAgentNoErr
this is some test output
--- PASS: TestWrapAndCallAgentNoErr (0.01s)
=== RUN   TestWrapAndCallAgentValidationFails
--- PASS: TestWrapAndCallAgentValidationFails (0.02s)
=== RUN   TestWrapAndCallAgentAgentFails
--- PASS: TestWrapAndCallAgentAgentFails (0.03s)
=== RUN   TestWrapAndCallAgentInvocationFails
--- PASS: TestWrapAndCallAgentInvocationFails (0.04s)
=== RUN   TestWrapAndCallAgentMarshalFails
--- PASS: TestWrapAndCallAgentMarshalFails (0.05s)
=== RUN   TestDiffNoErr
--- PASS: TestDiffNoErr (0.06s)
=== RUN   TestDiffDecodeFails
--- PASS: TestDiffDecodeFails (0.07s)
=== RUN   TestDiffReadAllFails
--- PASS: TestDiffReadAllFails (0.08s)
=== RUN   TestPatchNoErr
--- PASS: TestPatchNoErr (0.09s)
=== RUN   TestPatchAgentErr
--- PASS: TestPatchAgentErr (0.10s)
=== RUN   TestPatchDryRun
this is some test output
--- PASS: TestPatchDryRun (0.11s)
PASS
coverage: 57.1% of statements
`

	SinglePackageVerboseFailure = `=== RUN   TestWrapAndCallAgentNoErr
--- PASS: TestWrapAndCallAgentNoErr (1.00s)
=== RUN   TestWrapAndCallAgentValidationFails
--- PASS: TestWrapAndCallAgentValidationFails (1.00s)
=== RUN   TestWrapAndCallAgentAgentFails
--- PASS: TestWrapAndCallAgentAgentFails (1.00s)
=== RUN   TestWrapAndCallAgentInvocationFails
--- PASS: TestWrapAndCallAgentInvocationFails (1.00s)
=== RUN   TestWrapAndCallAgentMarshalFails
--- PASS: TestWrapAndCallAgentMarshalFails (1.00s)
=== RUN   TestDiffNoErr
--- PASS: TestDiffNoErr (1.00s)
=== RUN   TestDiffDecodeFails
--- PASS: TestDiffDecodeFails (1.00s)
=== RUN   TestDiffReadAllFails
--- PASS: TestDiffReadAllFails (1.00s)
=== RUN   TestPatchNoErr
this is some test output
--- PASS: TestPatchNoErr (1.00s)
=== RUN   TestPatchAgentErr
this is some other test output
--- FAIL: TestPatchAgentErr (1.00s)
	assert.go:143: got: (bool) true, want (bool) false in
		function          line file
		TestPatchAgentErr 220  agent/sdagentctl/helper/helper_test.go
		tRunner           473  go/src/testing/testing.go
		goexit            1998 go/src/runtime/asm_amd64.s
=== RUN   TestPatchDryRun
--- PASS: TestPatchDryRun (1.00s)
FAIL
exit status 1
`

	SinglePackageVerboseCoverageFailure = `=== RUN   TestWrapAndCallAgentNoErr
--- PASS: TestWrapAndCallAgentNoErr (1.00s)
=== RUN   TestWrapAndCallAgentValidationFails
--- PASS: TestWrapAndCallAgentValidationFails (1.00s)
=== RUN   TestWrapAndCallAgentAgentFails
--- PASS: TestWrapAndCallAgentAgentFails (1.00s)
=== RUN   TestWrapAndCallAgentInvocationFails
--- PASS: TestWrapAndCallAgentInvocationFails (1.00s)
=== RUN   TestWrapAndCallAgentMarshalFails
--- PASS: TestWrapAndCallAgentMarshalFails (1.00s)
=== RUN   TestDiffNoErr
--- PASS: TestDiffNoErr (1.00s)
=== RUN   TestDiffDecodeFails
--- PASS: TestDiffDecodeFails (1.00s)
=== RUN   TestDiffReadAllFails
--- PASS: TestDiffReadAllFails (1.00s)
=== RUN   TestPatchNoErr
this is some test output
--- PASS: TestPatchNoErr (1.00s)
=== RUN   TestPatchAgentErr
this is some other test output
--- FAIL: TestPatchAgentErr (1.00s)
	assert.go:143: got: (bool) true, want (bool) false in
		function          line file
		TestPatchAgentErr 220  agent/sdagentctl/helper/helper_test.go
		tRunner           473  go/src/testing/testing.go
		goexit            1998 go/src/runtime/asm_amd64.s
=== RUN   TestPatchDryRun
--- PASS: TestPatchDryRun (1.00s)
FAIL
coverage: 57.1% of statements
exit status 1
`
)

func testSuccess(t *testing.T, testdata string) {
	duration := 880 * time.Millisecond
	pkg := parseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(testdata)))

	assert.Equal(t, pkg.name, "github.com/turbinelabs/tbn/agent/sdagentctl/helper")
	assert.Equal(t, pkg.result, passed)
	assert.Equal(t, pkg.duration, 0.88)
	assert.Equal(t, len(pkg.tests), 11)
	assert.Equal(t, pkg.output, testdata)

	names := make([]string, len(pkg.tests))
	for i, test := range pkg.tests {
		names[i] = test.name
		assert.Equal(t, test.result, passed)
		assert.Equal(t, test.duration, 0.01*(1.0+float64(i)))
		assert.Equal(t, test.output.String(), "")
		assert.Equal(t, test.failure.String(), "")
	}

	assert.DeepEqual(t, names, []string{
		"TestWrapAndCallAgentNoErr",
		"TestWrapAndCallAgentValidationFails",
		"TestWrapAndCallAgentAgentFails",
		"TestWrapAndCallAgentInvocationFails",
		"TestWrapAndCallAgentMarshalFails",
		"TestDiffNoErr",
		"TestDiffDecodeFails",
		"TestDiffReadAllFails",
		"TestPatchNoErr",
		"TestPatchAgentErr",
		"TestPatchDryRun",
	})
}

func TestParseOutputOnVerboseSuccess(t *testing.T) {
	testSuccess(t, SinglePackageVerbose)
}

func TestParseOutputOnVerboseCoverageSuccess(t *testing.T) {
	testSuccess(t, SinglePackageVerboseCoverage)
}

func TestParseOutputOnVerboseCoverageOutputSuccess(t *testing.T) {
	duration := 880 * time.Millisecond
	pkg := parseTestOutput(
		testPackageName,
		duration,
		bytes.NewBuffer([]byte(SinglePackageVerboseCoverageOutput)))

	assert.Equal(t, len(pkg.tests), 11)

	for i, test := range pkg.tests {
		assert.Equal(t, test.result, passed)
		switch i {
		case 0, 10:
			assert.Equal(t, test.output.String(), "this is some test output\n")
		default:
			assert.Equal(t, test.output.String(), "")
		}

	}
}

func testFailure(t *testing.T, testdata string) {
	duration := 11 * time.Second
	pkg := parseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(testdata)))

	assert.Equal(t, pkg.name, "github.com/turbinelabs/tbn/agent/sdagentctl/helper")
	assert.Equal(t, pkg.result, failed)
	assert.Equal(t, pkg.duration, 11.0)
	assert.Equal(t, len(pkg.tests), 11)

	for i, test := range pkg.tests {
		if i == 9 {
			assert.Equal(t, test.result, failed)
			assert.MatchesRegex(t, test.failure.String(), `got: .*, want .*`)
			assert.True(t, len(strings.Split(test.failure.String(), "\n")) > 1)
		} else {
			assert.Equal(t, test.result, passed)
			assert.Equal(t, test.failure.String(), "")
		}

		switch i {
		case 8:
			assert.Equal(t, test.output.String(), "this is some test output\n")
		case 9:
			assert.Equal(t, test.output.String(), "this is some other test output\n")
		default:
			assert.Equal(t, test.output.String(), "")
		}

		assert.Equal(t, test.duration, 1.0)
	}
}

func TestParseOutputOnVerboseFailure(t *testing.T) {
	testFailure(t, SinglePackageVerboseFailure)
}

func TestParseOutputOnVerboseCoverageFailure(t *testing.T) {
	testFailure(t, SinglePackageVerboseCoverageFailure)
}
