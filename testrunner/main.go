/*
Copyright 2017 Turbine Labs, Inc.

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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/turbinelabs/test/testrunner/junit"
	"github.com/turbinelabs/test/testrunner/parser"
	"github.com/turbinelabs/test/testrunner/results"
)

const (
	ENV_ROOT_PACKAGE = "TEST_RUNNER_ROOT_PACKAGE"
	ENV_OUTPUT_DIR   = "TEST_RUNNER_OUTPUT"
)

var RootPackage = getEnv(ENV_ROOT_PACKAGE, "github.com/turbinelabs")

var TestOutput = getEnv(ENV_OUTPUT_DIR, "testresults")

type lockedWriter struct {
	lock       *sync.Mutex
	underlying io.Writer
}

func (w *lockedWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.underlying.Write(p)
}

func newLockedWriter(underlying io.Writer) io.Writer {
	return &lockedWriter{
		lock:       &sync.Mutex{},
		underlying: underlying,
	}
}

func main() {
	testExecutable := os.Args[1]
	pkgName := extractPackageFromTestExecutable(testExecutable)
	pkgFileName := strings.Replace(pkgName, "/", ".", -1)

	testArgs := parser.GoLangParser.FlagFn(os.Args[2:])

	test := exec.Command(testExecutable, testArgs...)

	output := new(bytes.Buffer)
	outputWriter := newLockedWriter(output)

	test.Stdout = io.MultiWriter(outputWriter, os.Stdout)
	test.Stderr = io.MultiWriter(outputWriter, os.Stderr)

	start := time.Now()
	exitStatus := 0
	if err := test.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			switch sysStatus := e.ProcessState.Sys().(type) {
			case syscall.WaitStatus:
				exitStatus = sysStatus.ExitStatus()
			default:
				exitStatus = 1
			}
		default:
			panic(err)
		}
	}

	duration := time.Since(start)

	pkgs, err := parser.GoLangParser.ParseFn(pkgName, duration, output)
	if err != nil {
		panic(err)
	}

	// Parsing errors may result in the package being marked as a
	// failure even though the test binary reported success.
	// Convert exit status to failure since it may mean some test
	// results were not properly parsed.
	for _, pkg := range pkgs {
		if pkg.Result == results.Failed && exitStatus == 0 {
			exitStatus = 1
		}
	}

	report := openFile(pkgFileName)
	defer report.Close()

	junit.WriteReport(report, pkgs)

	os.Exit(exitStatus)
}

func openFile(pkgFileName string) *os.File {
	dirInfo, err := os.Stat(TestOutput)
	if err == nil {
		if !dirInfo.IsDir() {
			panic(
				fmt.Sprintf(
					"Env var %s=%s is not a directory",
					ENV_OUTPUT_DIR,
					TestOutput))
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(TestOutput, 0755); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	reportFile := filepath.Join(TestOutput, fmt.Sprintf("%s.xml", pkgFileName))
	report, err := os.OpenFile(reportFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	return report
}

func getEnv(name, defaultValue string) string {
	if value, present := os.LookupEnv(name); present {
		return value
	}
	return defaultValue
}

func extractPackageFromTestExecutable(exec string) string {
	start := strings.Index(exec, RootPackage)

	var pkg string
	if start == -1 {
		// expect ugly report output, but should still work
		pkg = exec
	} else {
		pkg = exec[start:]
	}

	// strip executable and test directory
	for {
		dir, base := filepath.Split(pkg)
		if strings.HasSuffix(base, ".test") || base == "_test" {
			pkg = filepath.Dir(dir)
		} else {
			break
		}
	}

	return pkg
}
