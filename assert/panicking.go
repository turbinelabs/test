package assert

import (
	"fmt"
	"testing"
)

// A PanickingTB embeds a testing.TB, overriding Fatalf and Fatal to panic.
type PanickingTB struct {
	testing.TB
}

// Panicking wraps a testing.T or testing.TB so Fatalf calls panic. This is
// useful when Fatalf is called from within goroutines.
func Panicking(t testing.TB) testing.TB {
	switch obj := t.(type) {
	case *G:
		return obj
	case *PanickingTB:
		return obj
	default:
		return &PanickingTB{t}
	}
}

// Fatalf invokes the underlying testing.TB's Fatalf function with the
// given error message and arguments.
func (*PanickingTB) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// Fatal invokes the underlying testing.TB's Fatal function with the
// given arguments.
func (*PanickingTB) Fatal(args ...interface{}) {
	panic(fmt.Sprint(args...))
}
