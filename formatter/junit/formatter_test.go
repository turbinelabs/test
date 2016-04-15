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
			Output:   "the output",
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
	assert.Equal(t, suite.Output.Contents, "the output")
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
