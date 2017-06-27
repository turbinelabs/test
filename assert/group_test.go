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

import (
	"strings"
	"testing"
)

func TestUngrouped(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	True(mockT, false)

	if len(mockT.log) != 1 || mockT.log[0].op != "Error" {
		tr.Errorf("got %+v, want single Error op", mockT.log)
	}

	expectedPrefix := "got (bool) false, want (bool) true"
	if !strings.HasPrefix(mockT.log[0].args, expectedPrefix) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestGroupPassing(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	Group("name", mockT, func(g *G) {
		True(g, true)
		g.Group("sub-group", func(g *G) {
			False(g, false)
		})
	})

	if len(mockT.log) != 0 {
		tr.Errorf("Expected no testing.T operations, got: %v", mockT.log)
	}
}

func TestGroupFailing(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	Group("name", mockT, func(g *G) {
		True(g, false)
		g.Group("sub-group", func(g *G) {
			False(g, false)
		})
	})

	if len(mockT.log) != 1 || mockT.log[0].op != "Error" {
		tr.Errorf("got %+v, want single Error op", mockT.log)
	}

	expectedPrefix := "name: "
	if !strings.HasPrefix(mockT.log[0].args, expectedPrefix) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestNestedGroupFailing(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	Group("main-group", mockT, func(g *G) {
		True(g, true)
		g.Group("sub-group", func(g *G) {
			True(g, false)
		})
	})

	if len(mockT.log) != 1 || mockT.log[0].op != "Error" {
		tr.Errorf("got %+v, want single Error op", mockT.log)
	}

	expectedPrefix := "main-group sub-group: "
	if !strings.HasPrefix(mockT.log[0].args, expectedPrefix) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestGroupErrof(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	Group("main-group", mockT, func(g *G) {
		g.Errorf("failed %s", "here")
	})

	if len(mockT.log) != 1 || mockT.log[0].op != "Errorf" {
		tr.Errorf("got %+v, want single Errorf op", mockT.log)
	}

	expectedPrefix := "main-group: failed here in "
	if !strings.HasPrefix(mockT.log[0].args, expectedPrefix) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[0].args, expectedPrefix)
	}
}

func TestGroupFatal(t *testing.T) {
	tr := Tracing(t)
	mockT := &mockT{}

	Group("main-group", mockT, func(g *G) {
		g.Fatal("boom")
		g.Fatalf("boom %s", "two")
	})

	if len(mockT.log) != 2 {
		tr.Errorf("expected 2 ops, got %d", len(mockT.log))
	}

	if mockT.log[0].op != "Fatal" {
		tr.Errorf("got %+v, want Fatal op", mockT.log[0])
	}
	if mockT.log[1].op != "Fatalf" {
		tr.Errorf("got %+v, want Fatalf op", mockT.log[1])
	}

	expectedPrefix1 := "main-group: boom in "
	if !strings.HasPrefix(mockT.log[0].args, expectedPrefix1) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[0].args, expectedPrefix1)
	}

	expectedPrefix2 := "main-group: boom two in "
	if !strings.HasPrefix(mockT.log[1].args, expectedPrefix2) {
		tr.Errorf("got '%s', want prefix '%s'", mockT.log[1].args, expectedPrefix2)
	}
}
