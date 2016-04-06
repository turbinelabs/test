package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type junitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []junitTestSuite
}

type junitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Duration   string          `xml:"time,attr"`
	Properties []junitProperty `xml:"properties>property,omitempty"`
	TestCases  []junitTestCase
	Output     *junitOutput `xml:"system-out,omitempty"`
}

type junitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"name,attr"`
}

type junitTestCase struct {
	XMLName   xml.Name          `xml:"testcase"`
	Classname string            `xml:"classname,attr"`
	Name      string            `xml:"name,attr"`
	Duration  string            `xml:"time,attr"`
	Skipped   *junitSkipMessage `xml:"skipped,omitempty"`
	Failure   *junitFailure     `xml:"failure,omitempty"`
	Output    *junitOutput      `xml:"system-out,omitempty"`
}

type junitSkipMessage struct {
	Message string `xml:",cdata"`
}

type junitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",cdata"`
}

type junitOutput struct {
	Contents string `xml:",cdata"`
}

func writeReport(out io.Writer, pkg *testPackage) {
	suites := generateReport(pkg)

	xml, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		panic(err)
	}
	out.Write(xml)
	out.Write([]byte{'\n'})
}

func generateReport(pkg *testPackage) junitTestSuites {
	suite := junitTestSuite{
		Name:       pkg.name,
		Tests:      len(pkg.tests),
		Duration:   formatDuration(pkg.duration),
		Properties: []junitProperty{},
		TestCases:  []junitTestCase{},
	}

	if pkg.output != "" {
		suite.Output = &junitOutput{pkg.output}
	}

	classname := pkg.name
	if i := strings.LastIndex(classname, "/"); i >= 0 && i < len(classname) {
		classname = classname[i+1:]
	}

	for _, test := range pkg.tests {
		output := test.output.String()
		var testCaseOutput *junitOutput
		if output != "" {
			testCaseOutput = &junitOutput{output}
		}

		testCase := junitTestCase{
			Classname: classname,
			Name:      test.name,
			Duration:  formatDuration(test.duration),
			Output:    testCaseOutput,
		}

		switch test.result {
		case failed:
			suite.Failures++
			testCase.Failure = &junitFailure{
				Message:  "Failed",
				Contents: test.failure.String(),
			}
		case skipped:
			testCase.Skipped = &junitSkipMessage{output}
		}

		suite.TestCases = append(suite.TestCases, testCase)
	}

	return junitTestSuites{Suites: []junitTestSuite{suite}}
}

func formatDuration(f float64) string {
	return fmt.Sprintf("%.3f", f)
}
