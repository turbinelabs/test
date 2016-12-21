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

/*
  Package runner executes a test executable in a subprocess, parses
  its go-style output to calculate test results, and generates a
  junit-style test result XML file.

  testrunner takes no command line arguments directly. Its first
  argument must be the path of a test executable. Subsequent arguments
  are passed to that executable (in a subprocess).  Specific test
  parsers may modify the arguments passed to the subprocess.


  Environment Variables

  testrunner is configured using environment variables (to avoid
  potential collisions with test flag names). The available settings
  are:

  TEST_RUNNER_ROOT_PACKAGE (default: "github.com/turbinelabs") -
  Sets the root package name. By default, the package name in the
  JUnit formatted output is the test executable's path.  If the root
  package name is found in the test executable path, then the portion
  of the executable path up to the root package name is removed. For
  example, if the root package is set to "a/b/c" and the executable is
  "/home/foo/bar/a/b/c/blah", the package will be "a/b/c/blah".

  TEST_RUNNER_OUTPUT (default: "testresults") - Sets the path to the
  directory where formatted test results are written.


  Parsers

  Currently only a single parser exists. It parses go-style test
  verbose output. It modifies test command line arguments to include
  the verbose flag ("-test.v=true").
*/
package main
