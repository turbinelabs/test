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

	result := extractPackageFromTestExecutable(
		"/var/something/tmp-XYZZYZ/github.com/foo/bar/package/sub/_test/sub.test")

	assert.Equal(t, result, "github.com/foo/bar/package/sub")
}
