// Package matcher provides helpful implementations of the gomock.Matcher interface.
package matcher

import (
	"fmt"
	"io"
)

// AnyWriter implements gomock.Matcher, and can be used
// to mock a call that passes an io.Writer, filling that
// writer with the provided byte array
type AnyWriter struct {
	Data []byte
}

func (a AnyWriter) Matches(x interface{}) bool {
	if writer, ok := x.(io.Writer); ok {
		writer.Write(a.Data)
		return true
	}
	return false
}

func (a AnyWriter) String() string {
	return fmt.Sprintf("AnyWriter(%q)", a.Data)
}
