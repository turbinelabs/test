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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func withTrappedOutput(f func()) string {
	buffer := &bytes.Buffer{}
	saved := out
	defer func() { out = saved }()

	out = buffer
	f()
	return buffer.String()
}

func resetFlags() {
	ports = nil
	portsList = ""
	errorRate = 0
	latencyMeanMs = 0
	latencyStdDevMs = 0
	verbose = false
	help = false
}

func testParse(args []string, f func(rc int)) {
	defer resetFlags()

	fs := configureFlags()
	f(parseFlags(fs, args))
}

func testRun(t *testing.T, args []string, f func(rc int)) {
	defer resetFlags()

	fs := configureFlags()
	rc := parseFlags(fs, args)
	assert.Equal(t, rc, 0)
	f(run(fs))
}

func TestHelp(t *testing.T) {
	output := withTrappedOutput(func() {
		testParse([]string{"--help"}, func(rc int) {
			assert.Equal(t, rc, 1)
		})
	})

	assert.StringContains(t, output, "USAGE")
}

func TestBadFlags(t *testing.T) {
	output := withTrappedOutput(func() {
		testParse([]string{"--ports=1234", "--blah"}, func(rc int) {
			assert.Equal(t, rc, 1)
		})
	})

	assert.StringContains(t, output, "flag provided but not defined")
	assert.StringContains(t, output, "USAGE")
}

func TestTooManyArguments(t *testing.T) {
	output := withTrappedOutput(func() {
		testParse([]string{"--ports=1234", "blah"}, func(rc int) {
			assert.Equal(t, rc, 1)
		})
	})

	assert.StringContains(t, output, "too many arguments")
	assert.StringContains(t, output, "USAGE")
}

func TestNoPorts(t *testing.T) {
	output := withTrappedOutput(func() {
		testParse([]string{"--ports="}, func(rc int) {
			assert.Equal(t, rc, 1)
		})
	})

	assert.StringContains(t, output, "no listener port")
	assert.StringContains(t, output, "USAGE")
}

func TestEmptyPorts(t *testing.T) {
	output := withTrappedOutput(func() {
		testParse([]string{"--ports=,, 1 ,2,,3, ,"}, func(rc int) {
			assert.Equal(t, rc, 0)
			assert.ArrayEqual(t, ports, []string{"1", "2", "3"})
		})
	})

	assert.Equal(t, output, "")
}

func TestRunError(t *testing.T) {
	output := withTrappedOutput(func() {
		testRun(t, []string{"--error-rate=999"}, func(rc int) {
			assert.NotEqual(t, rc, 0)
		})
	})

	assert.StringContains(t, output, "error rate must be between")
}
