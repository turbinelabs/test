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

package tempfile

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestTempDir(t *testing.T) {
	mockT := new(testing.T)

	dir := TempDir(mockT)
	assert.False(t, mockT.Failed())
	assert.StringContains(t, dir.Path(), "test-tmp.")

	contents, err := ioutil.ReadDir(dir.Path())
	assert.Nil(t, err)
	assert.Equal(t, len(contents), 0)

	dir.Cleanup()
	_, err = ioutil.ReadDir(dir.Path())
	assert.NonNil(t, err)
}

func TestTempDirMake(t *testing.T) {
	mockT := new(testing.T)

	dir := TempDir(mockT, "testdir")
	assert.False(t, mockT.Failed())
	assert.StringContains(t, dir.Path(), "testdir.")

	filename1 := dir.Make(mockT)
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename1, "test-tmp."))
	assert.Equal(t, path.Dir(filename1), dir.Path())

	filename2 := dir.Make(mockT, "xyz")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename2, "xyz."))
	assert.Equal(t, path.Dir(filename2), dir.Path())

	contents, err := ioutil.ReadDir(dir.Path())
	assert.Nil(t, err)
	assert.Equal(t, len(contents), 2)
	assert.Equal(t, path.Base(filename1), contents[0].Name())
	assert.Equal(t, path.Base(filename2), contents[1].Name())

	dir.Cleanup()
	_, err = ioutil.ReadDir(dir.Path())
	assert.NonNil(t, err)
}

func TestTempDirWrite(t *testing.T) {
	mockT := new(testing.T)

	dir := TempDir(mockT, "testdir")
	assert.False(t, mockT.Failed())
	assert.StringContains(t, dir.Path(), "testdir.")

	filename1 := dir.Write(mockT, "data1")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename1, "test-tmp."))
	assert.Equal(t, path.Dir(filename1), dir.Path())

	filename2 := dir.Write(mockT, "data2", "xyz")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename2, "xyz."))
	assert.Equal(t, path.Dir(filename2), dir.Path())

	contents, err := ioutil.ReadDir(dir.Path())
	assert.Nil(t, err)
	assert.Equal(t, len(contents), 2)
	assert.Equal(t, path.Base(filename1), contents[0].Name())
	assert.Equal(t, path.Base(filename2), contents[1].Name())

	data1, err := ioutil.ReadFile(filename1)
	assert.Nil(t, err)
	assert.Equal(t, string(data1), "data1")

	data2, err := ioutil.ReadFile(filename2)
	assert.Nil(t, err)
	assert.Equal(t, string(data2), "data2")

	dir.Cleanup()
	_, err = ioutil.ReadDir(dir.Path())
	assert.NonNil(t, err)
}
