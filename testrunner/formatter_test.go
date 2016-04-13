package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

var (
	passingSuite = &testPackage{
		name:     "github.com/turbinelabs/something",
		result:   passed,
		duration: 1.234,
		tests: []*test{
			{
				name:     "TestFoo",
				result:   passed,
				duration: 1.2,
				output:   makeBuffer("some output"),
			},
			{
				name:     "TestBar",
				result:   passed,
				duration: 0.3,
			},
		},
	}

	failingSuite = &testPackage{
		name:     "github.com/turbinelabs/something",
		result:   failed,
		duration: 1.234,
		tests: []*test{
			{
				name:     "TestFoo",
				result:   failed,
				duration: 1.2,
				failure:  makeBuffer("some assertion"),
				output:   makeBuffer("some output"),
			},
			{
				name:     "TestBar",
				result:   passed,
				duration: 0.3,
			},
		},
	}

	skippedSuite = &testPackage{
		name:     "github.com/turbinelabs/something",
		result:   passed,
		duration: 1.234,
		tests: []*test{
			{
				name:     "TestFoo",
				result:   passed,
				duration: 1.2,
			},
			{
				name:     "TestBar",
				result:   skipped,
				duration: 0.0,
				output:   makeBuffer("skipped it"),
			},
		},
	}
)

func makeBuffer(s string) bytes.Buffer {
	b := bytes.NewBufferString(s)
	return *b
}

func TestGenerateReportSuccess(t *testing.T) {
	suites := generateReport(passingSuite)
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
	assert.DeepEqual(t, testCase1.Output, &junitOutput{"some output"})

	testCase2 := suite.TestCases[1]
	assert.Equal(t, testCase2.Classname, "something")
	assert.Equal(t, testCase2.Name, "TestBar")
	assert.Equal(t, testCase2.Duration, "0.300")
	assert.Nil(t, testCase2.Skipped)
	assert.Nil(t, testCase2.Failure)
}

func TestGenerateReportFailure(t *testing.T) {
	suites := generateReport(failingSuite)
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
	assert.DeepEqual(t, testCase1.Failure, &junitFailure{
		Message:  "Failed",
		Contents: "some assertion",
	})
	assert.DeepEqual(t, testCase1.Output, &junitOutput{"some output"})

	testCase2 := suite.TestCases[1]
	assert.Equal(t, testCase2.Classname, "something")
	assert.Equal(t, testCase2.Name, "TestBar")
	assert.Equal(t, testCase2.Duration, "0.300")
	assert.Nil(t, testCase2.Skipped)
	assert.Nil(t, testCase2.Failure)
}

func TestGenerateReportSkipped(t *testing.T) {
	suites := generateReport(skippedSuite)
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
	assert.DeepEqual(t, testCase2.Skipped, &junitSkipMessage{
		Message: "skipped it",
	})
	assert.Nil(t, testCase2.Failure)
}

func TestWriteReportSuccess(t *testing.T) {
	var buf bytes.Buffer
	writeReport(&buf, passingSuite)

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
	writeReport(&buf, failingSuite)

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
	writeReport(&buf, skippedSuite)

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
