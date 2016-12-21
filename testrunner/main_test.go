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
