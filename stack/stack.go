// Package to produce user friendly stack traces.
//
// Usage:
//   s := stack.New()
//   stackStr := s.Format(true)
//   fmt.Println(stackStr)
package stack

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"
)

type frame struct {
	filepath string
	function string
	line     int
}

// Series of function calls with most recent in smaller indexes.
type Stack []frame

// Generate a new stack trace.
func New() Stack {
	var results Stack

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			name := fn.Name()

			s := strings.Split(name, "/")
			s = strings.Split(s[len(s)-1], ".")

			i := 0
			if len(s) > 1 {
				i = 1
			}

			name = strings.Join(s[i:], ".")

			results = append(results, frame{strings.TrimSpace(file), name, line})
		}
	}

	return results
}

// Remove frames from the top of the stack.
func (s *Stack) Pop(n int) error {
	if l := len(*s); l < n {
		return errors.New(fmt.Sprintf("Attempting to pop too many frames from stack; %d deep", l))
	}

	*s = (*s)[n:]

	return nil
}

// Because the depth in a certain library may not be consistent use this to
// remove the top of the stack that matches some prefix. If the stack contains
// an entry that doesn't match the prefix it assumes the remaining entries
// should stay.
//
// NB: Because it's prefix matching you should be cognizant of its interplay
// with TrimPaths.
func (s *Stack) PopFrames(prefixes ...string) {
	newStack := *s

	for _, f := range *s {
		prefixHit := false

		for _, p := range prefixes {
			if strings.HasPrefix(f.filepath, p) {
				prefixHit = true
				break
			}
		}

		if prefixHit {
			if len(newStack) > 0 {
				newStack = newStack[1:]
			}
		} else {
			break
		}
	}

	*s = newStack
}

// Examine filepaths and remove common path prefixes.
func (s Stack) TrimPaths(pathList ...string) {
	for i := range s {
		for _, pfx := range pathList {
			if strings.HasPrefix(s[i].filepath, pfx) {
				s[i].filepath = s[i].filepath[len(pfx):]
			}
		}
	}
}

// Produce a string containing the stack trace formated for consumption. If
// includeHeader is set will include a column header for function, line, and
// file.
//
// Example output:
//   --- FAIL: TestClusterConstraintsEqualsSuccess (0.00s)
//   	assert.go:60: got: (bool) true, want (bool) false in
//   		function                            file:line
//   		TestClusterConstraintsEqualsSuccess api/cluster_constraint_test.go:16
//   		tRunner                             go/src/testing/testing.go:456
//   		goexit                              go/src/runtime/asm_amd64.s:1696
func (s Stack) Format(includeHeader bool) string {
	buf := new(bytes.Buffer)
	w := new(tabwriter.Writer)

	w.Init(buf, 0, 8, 1, ' ', 0)
	if includeHeader {
		fmt.Fprintln(w, "function\tfile:line")
	}
	for _, f := range s {
		line := fmt.Sprintf("%s\t%s:%d", f.function, f.filepath, f.line)
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return buf.String()
}
