package main

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestExtractPackageFromTestExecutable(t *testing.T) {
	saved := RootPackage
	defer func() {
		RootPackage = saved
	}()
	RootPackage = "github.com/foo/bar"

	result := extractPackageFromTestExecutable("/var/something/tmp-XYZZYZ/github.com/foo/bar/package/sub/_test/sub.test")

	assert.Equal(t, result, "github.com/foo/bar/package/sub")
}

func TestForceVerboseFlag(t *testing.T) {
	nonVerboseArgs := []string{"-test.timeout=4s"}
	result := forceVerboseFlag(nonVerboseArgs)
	assert.Equal(t, result[len(result)-1], "-test.v=true")

	verboseArgs := []string{"-test.v=true", "-test.timeout=4s"}
	result = forceVerboseFlag(verboseArgs)
	assert.DeepEqual(t, result, verboseArgs)

	verboseOffArgs := []string{"-test.v=false"}
	result = forceVerboseFlag(verboseOffArgs)
	assert.DeepEqual(t, result, []string{"-test.v=true"})
}
