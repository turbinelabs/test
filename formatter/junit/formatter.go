// Package junit provides a formatter for generating junit-style test results.
package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/turbinelabs/test/testrunner/results"
)

type JunitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JunitTestSuite
}

type JunitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Duration   string          `xml:"time,attr"`
	Properties []JunitProperty `xml:"properties>property,omitempty"`
	TestCases  []JunitTestCase
	Output     *JunitOutput `xml:"system-out,omitempty"`
}

type JunitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"name,attr"`
}

type JunitTestCase struct {
	XMLName   xml.Name          `xml:"testcase"`
	Classname string            `xml:"classname,attr"`
	Name      string            `xml:"name,attr"`
	Duration  string            `xml:"time,attr"`
	Skipped   *JunitSkipMessage `xml:"skipped,omitempty"`
	Failure   *JunitFailure     `xml:"failure,omitempty"`
	Output    *JunitOutput      `xml:"system-out,omitempty"`
}

type JunitSkipMessage struct {
	Message string `xml:",cdata"`
}

type JunitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",cdata"`
}

type JunitOutput struct {
	Contents string `xml:",cdata"`
}

func WriteReport(out io.Writer, pkgs []*results.TestPackage) {
	suites := GenerateReport(pkgs)

	xml, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		panic(err)
	}
	out.Write(xml)
	out.Write([]byte{'\n'})
}

func GenerateReport(pkgs []*results.TestPackage) JunitTestSuites {
	suites := make([]JunitTestSuite, len(pkgs))

	for i, pkg := range pkgs {
		suite := JunitTestSuite{
			Name:       pkg.Name,
			Tests:      len(pkg.Tests),
			Duration:   formatDuration(pkg.Duration),
			Properties: []JunitProperty{},
			TestCases:  []JunitTestCase{},
		}

		if pkg.Output != "" {
			suite.Output = &JunitOutput{pkg.Output}
		}

		classname := pkg.Name
		if i := strings.LastIndex(classname, "/"); i >= 0 && i < len(classname) {
			classname = classname[i+1:]
		}

		for _, test := range pkg.Tests {
			output := test.Output.String()
			var testCaseOutput *JunitOutput
			if output != "" {
				testCaseOutput = &JunitOutput{output}
			}

			testCase := JunitTestCase{
				Classname: classname,
				Name:      test.Name,
				Duration:  formatDuration(test.Duration),
				Output:    testCaseOutput,
			}

			switch test.Result {
			case results.Failed:
				suite.Failures++
				testCase.Failure = &JunitFailure{
					Message:  "Failed",
					Contents: test.Failure.String(),
				}
			case results.Skipped:
				testCase.Skipped = &JunitSkipMessage{output}
			}

			suite.TestCases = append(suite.TestCases, testCase)
		}

		suites[i] = suite
	}
	return JunitTestSuites{Suites: suites}
}

func formatDuration(f float64) string {
	return fmt.Sprintf("%.3f", f)
}
