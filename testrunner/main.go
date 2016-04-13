package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Invocation:
// go test -exec testrunner ...
//
// go test will execute testrunner once for each package, passing the
// package test executable and all test flags as arguments.

const (
	ENV_ROOT_PACKAGE = "TEST_RUNNER_ROOT_PACKAGE"
	ENV_OUTPUT_DIR   = "TEST_RUNNER_OUTPUT"
)

var RootPackage = getEnv(ENV_ROOT_PACKAGE, "github.com/turbinelabs")

var TestOutput = getEnv(ENV_OUTPUT_DIR, "testresults")

func main() {
	testExecutable := os.Args[1]
	testArgs := forceVerboseFlag(os.Args[2:])

	pkgName := extractPackageFromTestExecutable(testExecutable)
	pkgFileName := strings.Replace(pkgName, "/", ".", -1)

	test := exec.Command(testExecutable, testArgs...)

	var output bytes.Buffer

	test.Stdout = io.MultiWriter(&output, os.Stdout)
	test.Stderr = io.MultiWriter(&output, os.Stderr)

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

	pkg := parseTestOutput(pkgName, duration, &output)

	// Parsing errors may result in the package being marked as a
	// failure even though the test binary reported success.
	// Convert exit status to failure since it may mean some test
	// results were not properly parsed.
	if pkg.result == failed && exitStatus == 0 {
		exitStatus = 1
	}

	report := openFile(pkgFileName)
	defer report.Close()

	writeReport(report, pkg)

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

func forceVerboseFlag(args []string) []string {
	for i, arg := range args {
		if arg == "-test.v=true" {
			return args
		} else if arg == "-test.v=false" {
			args = append(args[0:i], args[i+1:]...)
			break
		}
	}

	return append(args, "-test.v=true")
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
