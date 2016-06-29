package tempfile

// Package tempfile provides wrappers around `ioutil.TempFile` to
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
// the cleanup function is passed to `defer`. Failure to create the
// file causes a fatal error via the testing context.
func Make(t *testing.T, prefix ...string) (string, func()) {
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

// Writes the given data to a newly create temporary file. Uses `Make`
// to create the file.
func Write(t *testing.T, data string, prefix ...string) (string, func()) {
	filename, cleanup := Make(t, prefix...)
	err := ioutil.WriteFile(filename, []byte(data), defaultPermissions)
	if err != nil {
		t.Fatalf("failed to write temp file for test: %v", err)
		return "", func() {}
	}
	return filename, cleanup
}

// Appends the given data to a previously created file.
func Append(t *testing.T, file, data string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, defaultPermissions)
	if err != nil {
		t.Fatalf("failed to append to temp file '%s' for test: %v", file, err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return err
}
