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

package assert

import "testing"

func TestTracing(t *testing.T) {
	switch obj := Tracing(t).(type) {
	case *TracingTB:
		Equal(obj, obj.TB, t)
	default:
		obj.Errorf("got *TracingTB, want %T", obj)
	}
}

func TestTracingNoWrap(t *testing.T) {
	tr := Tracing(t)
	obj := Tracing(tr)
	Equal(tr, tr, obj)
}

func TestTracingNoWrapG(t *testing.T) {
	Group("Foo", t, func(g *G) {
		obj := Tracing(g)
		Equal(g, g, obj)
	})
}
