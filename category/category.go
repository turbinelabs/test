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
