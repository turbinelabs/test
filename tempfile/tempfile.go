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

package tempfile

// Package tempfile provides wrappers around ioutil.TempFile to
// easily create temporary files or file names.

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const defaultPermissions = 0666

// Generates an empty temporary file and returns its name and a
// cleanup function that removes the file. The prefix is optional and
// may be used to control the name of the temporary file. Typically,
// the cleanup function is passed to defer. Failure to create the
// file causes a fatal error via the testing context.
func Make(t testing.TB, prefix ...string) (string, func()) {
	p := "test-tmp."
	if len(prefix) > 0 {
		p = strings.Join(prefix, "-")
		if !strings.HasSuffix(p, ".") {
			p = p + "."
		}
	}

	f, err := ioutil.TempFile("", p)
	if err != nil {
		t.Fatalf("failed to create temp file for test: %v", err)
		return "", func() {}
	}
	defer f.Close()

	name := f.Name()
	return name, func() { os.Remove(name) }
}

// Writes the given data to a newly create temporary file. Uses Make
// to create the file.
func Write(t testing.TB, data string, prefix ...string) (string, func()) {
	filename, cleanup := Make(t, prefix...)
	err := ioutil.WriteFile(filename, []byte(data), defaultPermissions)
	if err != nil {
		t.Fatalf("failed to write temp file for test: %v", err)
		return "", func() {}
	}
	return filename, cleanup
}

// Appends the given data to a previously created file.
func Append(t testing.TB, file, data string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, defaultPermissions)
	if err != nil {
		t.Fatalf("failed to append to temp file '%s' for test: %v", file, err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return err
}
