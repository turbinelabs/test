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

package matcher

import (
	"fmt"
)

// And creates a Matcher from two undelrying Matchers. Both underlying
// Matchers must match for this Matcher to match. The second matcher
// is not checked if the first does not match.
func And(m1, m2 Matcher) Matcher {
	return &and{m1, m2}
}

type and struct {
	m1, m2 Matcher
}

func (a *and) Matches(i interface{}) bool {
	return a.m1.Matches(i) && a.m2.Matches(i)
}

func (a *and) String() string {
	return fmt.Sprintf("(%s) and (%s)", a.m1.String(), a.m2.String())
}

// Or creates a Matcher from two undelrying Matchers. One of the
// underlying Matchers must match for this Matcher to match. The
// second Matcher is not checked if the first matches.
func Or(m1, m2 Matcher) Matcher {
	return &or{m1, m2}
}

type or struct {
	m1, m2 Matcher
}

func (o *or) Matches(i interface{}) bool {
	return o.m1.Matches(i) || o.m2.Matches(i)
}

func (o *or) String() string {
	return fmt.Sprintf("(%s) or (%s)", o.m1.String(), o.m2.String())
}
