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

package junit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/testrunner/results"
)

var (
	passingSuite = []*results.TestPackage{
		{
			Name:     "github.com/turbinelabs/something",
			Result:   results.Passed,
			Duration: 1.234,
			Tests: []*results.Test{
				{
					Name:     "TestFoo",
					Result:   results.Passed,
					Duration: 1.2,
					Output:   makeBuffer("some output"),
				},
				{
					Name:     "TestBar",
					Result:   results.Passed,
					Duration: 0.3,
				},
			},
		},
	}

	failingSuite = []*results.TestPackage{
		{
			Name:     "github.com/turbinelabs/something",
			Result:   results.Failed,
			Duration: 1.234,
			Tests: []*results.Test{
				{
					Name:     "TestFoo",
					Result:   results.Failed,
					Duration: 1.2,
					Failure:  makeBuffer("some assertion"),
					Output:   makeBuffer("some output"),
				},
				{
					Name:     "TestBar",
					Result:   results.Passed,
					Duration: 0.3,
				},
			},
		},
	}

	skippedSuite = []*results.TestPackage{
		{
			Name:     "github.com/turbinelabs/something",
			Result:   results.Passed,
			Duration: 1.234,
			Tests: []*results.Test{
				{
					Name:     "TestFoo",
					Result:   results.Passed,
					Duration: 1.2,
				},
				{
					Name:     "TestBar",
					Result:   results.Skipped,
					Duration: 0.0,
					Output:   makeBuffer("skipped it"),
				},
			},
		},
	}

	suiteOutput = []*results.TestPackage{
		{
			Name:     "github.com/turbinelabs/tbn/something",
			Result:   results.Passed,
			Duration: 1.234,
			Output:   "the output, sanitized: \x00",
		},
	}
)

func makeBuffer(s string) bytes.Buffer {
	b := bytes.NewBufferString(s)
	return *b
}

func TestGenerateReportSuccess(t *testing.T) {
	suites := GenerateReport(passingSuite)
	assert.Equal(t, len(suites.Suites), 1)

	suite := suites.Suites[0]
	assert.Equal(t, suite.Name, "github.com/turbinelabs/something")
	assert.Equal(t, suite.Tests, 2)
	assert.Equal(t, suite.Failures, 0)
	assert.Equal(t, suite.Duration, "1.234")
	assert.Equal(t, len(suite.Properties), 0)
	assert.Equal(t, len(suite.TestCases), 2)

	testCase1 := suite.TestCases[0]
	assert.Equal(t, testCase1.Classname, "something")
	assert.Equal(t, testCase1.Name, "TestFoo")
	assert.Equal(t, testCase1.Duration, "1.200")
	assert.Nil(t, testCase1.Skipped)
	assert.Nil(t, testCase1.Failure)
	assert.DeepEqual(t, testCase1.Output, &JunitOutput{"some output"})

	testCase2 := suite.TestCases[1]
	assert.Equal(t, testCase2.Classname, "something")
	assert.Equal(t, testCase2.Name, "TestBar")
	assert.Equal(t, testCase2.Duration, "0.300")
	assert.Nil(t, testCase2.Skipped)
	assert.Nil(t, testCase2.Failure)
}

func TestGenerateReportFailure(t *testing.T) {
	suites := GenerateReport(failingSuite)
	assert.Equal(t, len(suites.Suites), 1)

	suite := suites.Suites[0]
	assert.Equal(t, suite.Name, "github.com/turbinelabs/something")
	assert.Equal(t, suite.Tests, 2)
	assert.Equal(t, suite.Failures, 1)
	assert.Equal(t, suite.Duration, "1.234")
	assert.Equal(t, len(suite.Properties), 0)
	assert.Equal(t, len(suite.TestCases), 2)

	testCase1 := suite.TestCases[0]
	assert.Equal(t, testCase1.Classname, "something")
	assert.Equal(t, testCase1.Name, "TestFoo")
	assert.Equal(t, testCase1.Duration, "1.200")
	assert.Nil(t, testCase1.Skipped)
	assert.DeepEqual(t, testCase1.Failure, &JunitFailure{
		Message:  "Failed",
		Contents: "some assertion",
	})
	assert.DeepEqual(t, testCase1.Output, &JunitOutput{"some output"})

	testCase2 := suite.TestCases[1]
	assert.Equal(t, testCase2.Classname, "something")
	assert.Equal(t, testCase2.Name, "TestBar")
	assert.Equal(t, testCase2.Duration, "0.300")
	assert.Nil(t, testCase2.Skipped)
	assert.Nil(t, testCase2.Failure)
}

func TestGenerateReportSkipped(t *testing.T) {
	suites := GenerateReport(skippedSuite)
	assert.Equal(t, len(suites.Suites), 1)

	suite := suites.Suites[0]
	assert.Equal(t, suite.Name, "github.com/turbinelabs/something")
	assert.Equal(t, suite.Tests, 2)
	assert.Equal(t, suite.Failures, 0)
	assert.Equal(t, suite.Duration, "1.234")
	assert.Equal(t, len(suite.Properties), 0)
	assert.Equal(t, len(suite.TestCases), 2)

	testCase1 := suite.TestCases[0]
	assert.Equal(t, testCase1.Classname, "something")
	assert.Equal(t, testCase1.Name, "TestFoo")
	assert.Equal(t, testCase1.Duration, "1.200")
	assert.Nil(t, testCase1.Skipped)
	assert.Nil(t, testCase1.Failure)

	testCase2 := suite.TestCases[1]
	assert.Equal(t, testCase2.Classname, "something")
	assert.Equal(t, testCase2.Name, "TestBar")
	assert.Equal(t, testCase2.Duration, "0.000")
	assert.DeepEqual(t, testCase2.Skipped, &JunitSkipMessage{
		Message: "skipped it",
	})
	assert.Nil(t, testCase2.Failure)
}

func TestGenerateReportCombined(t *testing.T) {
	combinedInput := append([]*results.TestPackage{}, passingSuite...)
	combinedInput = append(combinedInput, failingSuite...)
	combinedInput = append(combinedInput, skippedSuite...)

	suites := GenerateReport(combinedInput)
	assert.Equal(t, len(suites.Suites), 3)
}

func TestGenerateReportSuiteOutput(t *testing.T) {
	suites := GenerateReport(suiteOutput)
	assert.Equal(t, len(suites.Suites), 1)

	suite := suites.Suites[0]
	assert.NonNil(t, suite.Output)
	assert.Equal(t, suite.Output.Contents, `the output, sanitized: \x00`)
}

func TestWriteReportSuccess(t *testing.T) {
	var buf bytes.Buffer
	WriteReport(&buf, passingSuite)

	s := strings.Replace(buf.String(), "\n", "", -1)

	assert.MatchesRegex(t, s, `^<testsuites>.*</testsuites>$`)
	assert.MatchesRegex(t, s, `<testsuite .*name="github.com/turbinelabs/something".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*tests="2".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*failures="0".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*time="1.234".*>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestFoo".*>\s*<system-out><!\[CDATA\[some output\]\]></system-out>\s*</testcase>`)
	assert.MatchesRegex(t, s, `<testcase .*classname="something".*>`)
	assert.MatchesRegex(t, s, `<testcase .*time="1.200".*>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestBar".*></testcase>`)
	assert.MatchesRegex(t, s, `<testcase .*time="0.300".*></testcase>`)
}

func TestWriteReportFailure(t *testing.T) {
	var buf bytes.Buffer
	WriteReport(&buf, failingSuite)

	s := strings.Replace(buf.String(), "\n", "", -1)

	assert.MatchesRegex(t, s, `^<testsuites>.*</testsuites>$`)
	assert.MatchesRegex(t, s, `<testsuite .*name="github.com/turbinelabs/something".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*tests="2".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*failures="1".*>`)
	assert.MatchesRegex(t, s, `<testcase .*classname="something".*>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestFoo".*>`)
	assert.MatchesRegex(t, s, `<failure .*message="Failed".*><!\[CDATA\[some assertion\]\]></failure>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestBar".*></testcase>`)
}

func TestWriteReportSkipped(t *testing.T) {
	var buf bytes.Buffer
	WriteReport(&buf, skippedSuite)

	s := strings.Replace(buf.String(), "\n", "", -1)

	assert.MatchesRegex(t, s, `^<testsuites>.*</testsuites>$`)
	assert.MatchesRegex(t, s, `<testsuite .*name="github.com/turbinelabs/something".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*tests="2".*>`)
	assert.MatchesRegex(t, s, `<testsuite .*failures="0".*>`)
	assert.MatchesRegex(t, s, `<testcase .*classname="something".*>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestFoo".*></testcase>`)
	assert.MatchesRegex(t, s, `<testcase .*name="TestBar".*>`)
	assert.MatchesRegex(t, s, `<skipped><!\[CDATA\[skipped it\]\]></skipped>`)
}

func TestSanitize(t *testing.T) {
	testcases := [][]string{
		{"abc is ok by me", "abc is ok by me"},
		{`\ is safe`, `\ is safe`},
		{"123\t\r\n", "123\t\r\n"},               // CR, LF, and TAB are ok
		{"nope \ufdd0 nope", `nope \ufdd0 nope`}, // private use characters are not allowed
		{"null \x00 null", `null \x00 null`},     // null is not allowed
		{"esc \x1b esc", `esc \x1b esc`},         // escape is not allowed
		{"\x1f", `\x1f`},                         // whatever this is: not allowed
		{"del \x7f del", `del \x7f del`},         // delete is not allowed
		{"\x80", `\x80`},                         // bare continuation byte
		{"\xbf", `\xbf`},                         // bare continuation byte
		{"\x80\xbf", `\x80\xbf`},                 // pair of bare continuation bytes
		{"\xc0 ", `\xc0 `},                       // missing continuation byte
		{"\xf0\x80\x80 ", `\xf0\x80\x80 `},       // missing last continuation byte
		{"\U0001FFFF", `\U0001ffff`},             // non-characters
	}

	for _, testcase := range testcases {
		input := testcase[0]
		expected := testcase[1]

		assert.Equal(t, sanitize(input), expected)
	}
}
