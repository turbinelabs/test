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

// Allows for the categorization and selective execution of tests.
package category

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type TestCategory string

const (
	IntegrationTest TestCategory = "integration"
)

func (c TestCategory) EnvName() string {
	return fmt.Sprintf("%s_TEST", strings.ToUpper(string(c)))
}

func SkipUnless(t *testing.T, category TestCategory) {
	if os.Getenv(category.EnvName()) == "" {
		t.Skipf("Skipping test: set %s to enable this test", category.EnvName())
	}
}
