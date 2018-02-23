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

package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/testrunner/results"
)

const (
	testPackageName = "foo/bar/baz"

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
		TestPatchAgentErr 220  agent/differctl/helper/helper_test.go
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
		TestPatchAgentErr 220  agent/differctl/helper/helper_test.go
		tRunner           473  go/src/testing/testing.go
		goexit            1998 go/src/runtime/asm_amd64.s
=== RUN   TestPatchDryRun
--- PASS: TestPatchDryRun (1.00s)
FAIL
coverage: 57.1% of statements
exit status 1
`
	SkippedTests = `=== RUN   TestPatchAgentErr
--- PASS: TestPatchAgentErr (0.10s)
=== RUN   TestPatchDryRun
--- SKIP: TestPatchDryRun (0.0s)
	some_path.go:101: reason
PASS
coverage: 57.1% of statements
`

	NoTests = `PASS
`

	NoPackageResultFailure = `=== RUN   TestPatchAgentErr
--- PASS: TestPatchAgentErr (0.10s)
`

	MultiplePackageFailure = `=== RUN   TestPatchAgentErr1
--- PASS: TestPatchAgentErr1 (0.10s)
PASS
=== RUN   TestPatchAgentErr2
--- PASS: TestPatchAgentErr2 (0.10s)
PASS
`
)

func testSuccess(t *testing.T, testdata string) {
	duration := 880 * time.Millisecond
	pkgs, err := ParseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(testdata)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	pkg := pkgs[0]
	assert.Equal(t, pkg.Name, "foo/bar/baz")
	assert.Equal(t, pkg.Result, results.Passed)
	assert.Equal(t, pkg.Duration, 0.88)
	assert.Equal(t, len(pkg.Tests), 11)
	assert.Equal(t, pkg.Output, testdata)

	names := make([]string, len(pkg.Tests))
	for i, test := range pkg.Tests {
		names[i] = test.Name
		assert.Equal(t, test.Result, results.Passed)
		assert.Equal(t, test.Duration, 0.01*(1.0+float64(i)))
		assert.Equal(t, test.Output.String(), "")
		assert.Equal(t, test.Failure.String(), "")
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
	pkgs, err := ParseTestOutput(
		testPackageName,
		duration,
		bytes.NewBuffer([]byte(SinglePackageVerboseCoverageOutput)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	pkg := pkgs[0]
	assert.Equal(t, len(pkg.Tests), 11)

	for i, test := range pkg.Tests {
		assert.Equal(t, test.Result, results.Passed)
		switch i {
		case 0, 10:
			assert.Equal(t, test.Output.String(), "this is some test output\n")
		default:
			assert.Equal(t, test.Output.String(), "")
		}
	}
}

func testFailure(t *testing.T, testdata string) {
	duration := 11 * time.Second
	pkgs, err := ParseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(testdata)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	pkg := pkgs[0]
	assert.Equal(t, pkg.Name, "foo/bar/baz")
	assert.Equal(t, pkg.Result, results.Failed)
	assert.Equal(t, pkg.Duration, 11.0)
	assert.Equal(t, len(pkg.Tests), 11)

	for i, test := range pkg.Tests {
		if i == 9 {
			assert.Equal(t, test.Result, results.Failed)
			assert.MatchesRegex(t, test.Failure.String(), `got: .*, want .*`)
			assert.GreaterThan(t, len(strings.Split(test.Failure.String(), "\n")), 1)
		} else {
			assert.Equal(t, test.Result, results.Passed)
			assert.Equal(t, test.Failure.String(), "")
		}

		switch i {
		case 8:
			assert.Equal(t, test.Output.String(), "this is some test output\n")
		case 9:
			assert.Equal(t, test.Output.String(), "this is some other test output\n")
		default:
			assert.Equal(t, test.Output.String(), "")
		}

		assert.Equal(t, test.Duration, 1.0)
	}
}

func TestParseOutputOnVerboseFailure(t *testing.T) {
	testFailure(t, SinglePackageVerboseFailure)
}

func TestParseOutputOnVerboseCoverageFailure(t *testing.T) {
	testFailure(t, SinglePackageVerboseCoverageFailure)
}

func TestParseOutputOnSkippedTest(t *testing.T) {
	duration := 11 * time.Second
	pkgs, err := ParseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(SkippedTests)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	pkg := pkgs[0]
	assert.Equal(t, pkg.Name, "foo/bar/baz")
	assert.Equal(t, pkg.Result, results.Passed)
	assert.Equal(t, pkg.Duration, 11.0)
	assert.Equal(t, len(pkg.Tests), 2)

	test1 := pkg.Tests[0]
	assert.Equal(t, test1.Result, results.Passed)

	test2 := pkg.Tests[1]
	assert.Equal(t, test2.Result, results.Skipped)
	assert.MatchesRegex(t, test2.Failure.String(), `reason`)
}

func TestParseOutputOnNoTests(t *testing.T) {
	duration := 11 * time.Second
	pkgs, err := ParseTestOutput(testPackageName, duration, bytes.NewBuffer([]byte(NoTests)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	pkg := pkgs[0]
	assert.Equal(t, pkg.Name, "foo/bar/baz")
	assert.Equal(t, pkg.Result, results.Passed)
	assert.Equal(t, pkg.Duration, 11.0)
	assert.Equal(t, len(pkg.Tests), 0)
}

func TestParseOutputOnNoPackageResult(t *testing.T) {
	duration := 11 * time.Second
	pkgs, err := ParseTestOutput(
		testPackageName,
		duration,
		bytes.NewBuffer([]byte(NoPackageResultFailure)))
	assert.Nil(t, err)
	assert.Equal(t, len(pkgs), 1)

	// package failed
	pkg := pkgs[0]
	assert.Equal(t, pkg.Name, "foo/bar/baz")
	assert.Equal(t, pkg.Result, results.Failed)
	assert.Equal(t, pkg.Duration, 11.0)
	assert.Equal(t, len(pkg.Tests), 1)

	// test succeeded
	test1 := pkg.Tests[0]
	assert.Equal(t, test1.Result, results.Passed)
}

func TestParseOutputOnMultiplePackages(t *testing.T) {
	duration := 11 * time.Second
	pkgs, err := ParseTestOutput(
		testPackageName,
		duration,
		bytes.NewBuffer([]byte(MultiplePackageFailure)))
	assert.Nil(t, pkgs)
	assert.NonNil(t, err)
	assert.Equal(t, err.Error(), "expected only a single package")

}

func TestForceVerboseFlag(t *testing.T) {
	nonVerboseArgs := []string{"-test.timeout=4s"}
	result := ForceVerboseFlag(nonVerboseArgs)
	assert.Equal(t, result[len(result)-1], "-test.v=true")

	verboseArgs := []string{"-test.v=true", "-test.timeout=4s"}
	result = ForceVerboseFlag(verboseArgs)
	assert.DeepEqual(t, result, verboseArgs)

	verboseOffArgs := []string{"-test.v=false"}
	result = ForceVerboseFlag(verboseOffArgs)
	assert.DeepEqual(t, result, []string{"-test.v=true"})
}
