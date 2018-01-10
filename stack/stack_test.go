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

package stack

import (
	"os"
	"path"
	"strings"
	"testing"
)

func withMinDepth(n int) Stack {
	if n == 0 {
		return New()
	}

	fs := make([]func() Stack, n)

	fs[0] = func() Stack { return New() }
	for i := 1; i < n; i++ {
		j := i - 1
		fs[i] = func() Stack { return fs[j]() }
	}

	return fs[n-1]()
}

func TestStack(t *testing.T) {
	stack := New()

	if len(stack) == 0 {
		t.Errorf("expected non-empty stack")
	}

	for i, frame := range stack {
		if frame.function == "New" && strings.Contains(frame.filepath, "/test/stack/stack.go") {
			t.Errorf(
				"expected New to omit itself from the stack, but found stack[%d]: %v",
				i,
				frame,
			)
		}
	}
}

func TestStackPop(t *testing.T) {
	stack := withMinDepth(5)

	n := len(stack)
	if n < 5 {
		t.Errorf("expected stack with at least 5 frames, got %d", n)
	}

	originalStack := make(Stack, n)
	copy(originalStack, stack)

	if err := stack.Pop(n + 100); err == nil {
		t.Errorf("expected an error popping %d frames, but did not", n+100)
	}

	if err := stack.Pop(3); err != nil {
		t.Errorf("expected to successfully pop 3 frames, but got %v", err)
	}

	if len(stack) != n-3 {
		t.Errorf("expected stack to shrink to %d frames, got %d", n-3, len(stack))
	}

	for i := range stack {
		if stack[i] != originalStack[i+3] {
			t.Errorf("stack[%d]: expected %v, got %v", i, originalStack[i+3], stack[i])
		}
	}
}

func TestStackPopFrames(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to find working dir: %s", err.Error())
	}

	// $GOPATH/src/github.com/turbinelabs/test/stack -> $GOPATH/src
	srcDir := path.Dir(path.Dir(path.Dir(path.Dir(wd))))

	stack := withMinDepth(5)
	n := len(stack)

	stack.PopFrames(srcDir)

	if n-len(stack) < 7 {
		// 5 frames + withMinDepth itself + this function
		t.Errorf("expected at least 7 frames to be popped, but only popped %d", n-len(stack))
	}

	for _, frame := range stack {
		if strings.Contains(frame.filepath, "github.com/turbinelabs") {
			t.Errorf(
				"expected no frames containing github.com/turbinelabs, got\n%s",
				stack.Format(false),
			)
		}
	}
}

func TestTrimPaths(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to find working dir: %s", err.Error())
	}

	stack := withMinDepth(5)

	numMatching := 0
	for _, frame := range stack {
		if strings.HasPrefix(frame.filepath, wd) {
			numMatching++
		} else {
			break
		}
	}

	if numMatching == 0 {
		t.Errorf("expected some frames with file path matching prefix %s, but found none", wd)
	}

	// $GOPATH/src/github.com/turbinelabs/test/stack -> $GOPATH/src
	srcDir := path.Dir(path.Dir(path.Dir(path.Dir(wd)))) + "/"
	pkg := wd[len(srcDir):]

	stack.TrimPaths(srcDir)

	for i, frame := range stack {
		if i >= numMatching {
			break
		}

		if strings.HasPrefix(frame.filepath, wd) {
			t.Errorf("expected stack[%d] filepath to be trimmed, but got %q", i, frame.filepath)
		}

		if !strings.HasPrefix(frame.filepath, pkg) {
			t.Errorf(
				"expected stack[%d] filepath to start with %q, but got %q",
				i,
				pkg,
				frame.filepath,
			)
		}
	}

}

func TestFormat(t *testing.T) {
	stack := Stack{
		{"/foo/bar/code.go", "xyz", 100},
		{"/foo/bar/code.go", "Pdq", 10},
		{"code.go", "boingboing", 999},
	}

	header := `
function   file:line
`[1:]

	expected := `
xyz        /foo/bar/code.go:100
Pdq        /foo/bar/code.go:10
boingboing code.go:999
`[1:]

	formatted := stack.Format(false)
	if formatted != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, formatted)
	}

	formatted = stack.Format(true)
	if formatted != header+expected {
		t.Errorf("expected:\n%s\ngot:\n%s", header+expected, formatted)
	}
}

func TestFrames(t *testing.T) {
	stack := Stack{
		{"/foo/bar/code.go", "xyz", 100},
		{"/foo/bar/code.go", "Pdq", 10},
		{"code.go", "boingboing", 999},
	}
	frames := stack.Frames()

	expected := []string{
		"xyz (/foo/bar/code.go:100)",
		"Pdq (/foo/bar/code.go:10)",
		"boingboing (code.go:999)",
	}

	if len(frames) != len(expected) {
		t.Errorf("expected %d frames, got %d: %v", len(expected), len(frames), frames)
	}

	for i := range expected {
		if frames[i] != expected[i] {
			t.Errorf("expected frame[%d] to be %q, got %q", i, expected[i], frames[i])
		}
	}
}
