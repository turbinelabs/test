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
	"testing"
)

func TestUngrouped(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	True(mockT, false)

	mockT.CheckPredicates(
		tr,
		Match(
			ErrorOp(),
			PrefixedArgs("got (bool) false, want (bool) true"),
		),
	)
}

func TestGroupPassing(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	Group("name", mockT, func(g *G) {
		True(g, true)
		g.Group("sub-group", func(g *G) {
			False(g, false)
		})
	})

	mockT.CheckSuccess(tr)
}

func TestGroupFailing(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	Group("name", mockT, func(g *G) {
		True(g, false)
		g.Group("sub-group", func(g *G) {
			False(g, false)
		})
	})

	mockT.CheckPredicates(
		tr,
		Match(
			ErrorOp(),
			PrefixedArgs("name: "),
		),
	)
}

func TestNestedGroupFailing(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	Group("main-group", mockT, func(g *G) {
		True(g, true)
		g.Group("sub-group", func(g *G) {
			True(g, false)
		})
	})

	mockT.CheckPredicates(
		tr,
		Match(
			ErrorOp(),
			PrefixedArgs("main-group sub-group: "),
		),
	)
}

func TestGroupErrof(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	Group("main-group", mockT, func(g *G) {
		g.Errorf("failed %s", "here")
	})

	mockT.CheckPredicates(
		tr,
		Match(
			ErrorOp(),
			PrefixedArgs("main-group: failed here in "),
		),
	)
}

func TestGroupFatal(t *testing.T) {
	tr := Tracing(t)
	mockT := &MockT{}

	Group("main-group", mockT, func(g *G) {
		g.Fatal("boom")
		g.Fatalf("boom %s", "two")
	})

	mockT.CheckPredicates(
		tr,
		Match(ExactOp("Fatal"), PrefixedArgs("main-group: boom in ")),
		Match(ExactOp("Fatalf"), PrefixedArgs("main-group: boom two in ")),
	)
}
