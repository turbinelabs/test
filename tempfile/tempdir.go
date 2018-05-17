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

package tempfile

import (
	"io/ioutil"
	"os"
	"testing"
)

// Dir represents a temporary test directory.
type Dir interface {
	// Path returns the full path of the temporary directory.
	Path() string

	// Cleanup removes all files within the temporary directory and the directory
	// itself.
	Cleanup()

	// Make creates a temporary file within the temporary directory and returns its
	// path. See the package-level Make function.
	Make(t testing.TB, prefix ...string) string

	// Write creates a temporary file with the given contents within the temporary
	// directory and returns its path. See the package-level Write function.
	Write(t testing.TB, data string, prefix ...string) string
}

type tempDir struct {
	path string
}

// TempDir creates a temporary directory and returns a Dir object representing
// it. The prefix is optional and may be used to control the name of the temporary
// file. Failure to create the directory causes a fatal error via the testing
// context.
func TempDir(t testing.TB, prefix ...string) Dir {
	p := mkPrefix(prefix)

	tdir, err := ioutil.TempDir("", p)
	if err != nil {
		t.Fatalf("failed to create temp dir for test %v", err)
		return nil
	}

	return &tempDir{path: tdir}
}

func (d *tempDir) Path() string { return d.path }
func (d *tempDir) Cleanup()     { os.RemoveAll(d.path) }

func (d *tempDir) Make(t testing.TB, prefix ...string) string {
	return makeFile(t, d.path, prefix)
}

func (d *tempDir) Write(t testing.TB, data string, prefix ...string) string {
	return writeFile(t, d.path, data, prefix)
}
