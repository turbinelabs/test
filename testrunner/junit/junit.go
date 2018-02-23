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

// Package junit provides a formatter for generating junit-style test results.
package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

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
	Value string `xml:"value,attr"`
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

// WriteReport writes a []*results.TestPackage to a writer using
// the junit-standard format for test results.
func WriteReport(out io.Writer, pkgs []*results.TestPackage) {
	suites := GenerateReport(pkgs)

	xml, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		panic(err)
	}
	out.Write(xml)
	out.Write([]byte{'\n'})
}

// Escapes strings containing characters disallowed in XML (even in
// CDATA sections). The escaping is lossy and is meant only to prevent
// XML parsers from failing to decode the document.
func sanitize(s string) string {
	buf := make([]byte, 0, len(s))
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			// invalid UTF-8 encoding
			buf = append(buf, fmt.Sprintf(`\x%02x`, s[0])...)
			continue
		}

		// c.f. https://www.w3.org/TR/REC-xml/#charsets
		switch {
		case r == '\t' || r == '\r' || r == '\n':
			buf = append(buf, s[0])
		case r < ' ' || (r >= 0x7f && r <= 0x9f):
			buf = append(buf, fmt.Sprintf(`\x%02x`, s[0])...)
		case r >= 0xfdd0 && r <= 0xfdef:
			buf = append(buf, fmt.Sprintf(`\u%04x`, uint16(r))...)
		case r > utf8.MaxRune || r&0xFFFE == 0xFFFE:
			buf = append(buf, fmt.Sprintf(`\U%08x`, uint32(r))...)
		default:
			buf = append(buf, s[0:width]...)
		}
	}
	return string(buf)
}

// GenerateReport takes a []*results.TestPackage and produces a
// unitTestSuites suitable for marshalling as standard junit test results.
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
			suite.Output = &JunitOutput{sanitize(pkg.Output)}
		}

		classname := pkg.Name
		if i := strings.LastIndex(classname, "/"); i >= 0 && i < len(classname) {
			classname = classname[i+1:]
		}

		for _, test := range pkg.Tests {
			output := test.Output.String()
			var testCaseOutput *JunitOutput
			if output != "" {
				testCaseOutput = &JunitOutput{sanitize(output)}
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
