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
