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

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/doc"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/turbinelabs/test/server"
)

const (
	desc = `
Starts an HTTP test server that responds to HTTP requests on one or more ports
with configurable latency and error rate.

Every response from the server includes a header
("` + server.TestServerIDHeader + `") which reports the listener address that
served the request.

The status code for error responses is 503 (see below). All responses contain a
simple payload.

If the query parameter "` + server.TestServerForceResponseCode + `" is set and
contains an integer value greater than 0, the value is used as the response
code. The configured error rate continues to determine the payload of the
response.

If the query parameter "` + server.TestServerEchoHeadersWithPrefix + `" is set,
then success responses contain additional payload data displaying the name and
value of each HTTP request header that starts with the specified prefix. The
query parameter may be repeated to display headers with multiple prefixes.`
)

var (
	ports           []string
	portsList       string
	errorStatus     int
	errorRate       float64
	latencyMeanMs   float64
	latencyStdDevMs float64
	verbose         bool
	help            bool

	usageSections = []struct {
		text   string
		indent int
	}{
		{"NAME", 0},
		{"testserver - an HTTP test server", 4},
		{},
		{"USAGE", 0},
		{"testserver [OPTIONS]", 4},
		{},
		{"DESCRIPTION", 0},
		{desc, 4},
		{},
		{"OPTIONS", 0},
	}

	out io.Writer = os.Stderr
)

func stderr(s string, args ...interface{}) {
	fmt.Fprintf(out, s, args...)
}

func wrap(indent int, s string, args ...interface{}) {
	str := fmt.Sprintf(s, args...)
	indentStr := strings.Repeat(" ", indent)
	buffer := &bytes.Buffer{}
	doc.ToText(buffer, str, indentStr, "", 80)
	stderr(buffer.String())
}

func usage(fs *flag.FlagSet, err error) int {
	if err != nil && err != flag.ErrHelp {
		wrap(0, "Error: %s\n", err.Error())
		stderr("\n")
	}

	for _, u := range usageSections {
		if u.text == "" {
			stderr("\n")
		} else {
			wrap(u.indent, u.text+"\n")
		}
	}

	fs.VisitAll(func(f *flag.Flag) {
		name, usage := flag.UnquoteUsage(f)
		if name == "" {
			wrap(4, "--%s\n", f.Name)
		} else {
			wrap(4, "--%s=%s\n", f.Name, name)
		}
		if f.DefValue != "" {
			wrap(8, "(default: %s)", f.DefValue)
		}
		wrap(8, usage)
		stderr("\n")
	})

	return 1
}

func configureFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("testserver", flag.ContinueOnError)
	fs.StringVar(
		&portsList,
		"ports",
		"8889",
		"A comma-separated list of listener `ports` for the test server. The server listens on all interfaces.",
	)

	fs.IntVar(
		&errorStatus,
		"error-status",
		server.DefaultErrorStatus,
		"The HTTP status code the test server returns periodically when error-rate is non-zero.",
	)

	fs.Float64Var(
		&errorRate,
		"error-rate",
		0.01,
		"The test server's error rate as a `percentage`. A value of 100 means always fail.",
	)

	fs.Float64Var(
		&latencyMeanMs,
		"latency-mean",
		4.0,
		"The test server's mean latency in `milliseconds`. Note that latency skews slightly high because this setting effectively controls only the lower bound.",
	)

	fs.Float64Var(
		&latencyStdDevMs,
		"latency-stddev",
		1.0,
		"The test server's standard deviation from its mean latency in `milliseconds`.",
	)

	fs.BoolVar(
		&verbose,
		"verbose",
		false,
		"Enable verbose logging from the test server. Reports the success or failure of each response, and the latency for each request.",
	)

	fs.BoolVar(
		&help,
		"help",
		false,
		"Prints this help message.",
	)

	fs.SetOutput(ioutil.Discard)
	return fs
}

func parseFlags(fs *flag.FlagSet, args []string) int {
	if err := fs.Parse(args); err != nil {
		return usage(fs, err)
	}

	if fs.NArg() > 0 {
		return usage(fs, errors.New("too many arguments"))
	}

	if help {
		return usage(fs, nil)
	}

	ports = strings.Split(portsList, ",")
	i := 0
	for i < len(ports) {
		port := strings.TrimSpace(ports[i])
		if port == "" {
			copy(ports[i:], ports[i+1:])
			ports = ports[0 : len(ports)-1]
			continue
		}
		ports[i] = port
		i++
	}

	if len(ports) == 0 {
		return usage(fs, errors.New("no listener port(s) specified"))
	}

	return 0
}

func run(fs *flag.FlagSet) int {
	ts, err := server.NewTestServer(
		ports,
		errorRate,
		time.Duration(latencyMeanMs*float64(time.Millisecond)),
		time.Duration(latencyStdDevMs*float64(time.Millisecond)),
		verbose,
		nil,
	)
	if err != nil {
		return usage(fs, err)
	}

	if err := ts.SetErrorStatus(errorStatus); err != nil {
		return usage(fs, err)
	}

	// Blocks forever since there's no way to stop the server.
	ts.ServeAsync().Await()
	return 0
}

func main() {
	fs := configureFlags()
	rc := parseFlags(fs, os.Args[1:])
	if rc == 0 {
		rc = run(fs)
	}
	os.Exit(rc)
}
