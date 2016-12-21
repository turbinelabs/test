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

package tempfile

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestMake(t *testing.T) {
	mockT := new(testing.T)

	filename, cleanup := Make(mockT)
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename, "test-tmp."))
	assert.NonNil(t, cleanup)

	data, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	assert.Equal(t, string(data), "")

	cleanup()
	_, err = os.Stat(filename)
	assert.NonNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestMakePrefix(t *testing.T) {
	mockT := new(testing.T)

	filename, cleanup := Make(mockT, "foo", "bar")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename, "foo-bar."))
	assert.NonNil(t, cleanup)

	cleanup()
	_, err := os.Stat(filename)
	assert.NonNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestWrite(t *testing.T) {
	mockT := new(testing.T)

	filename, cleanup := Write(mockT, "stuff")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename, "test-tmp."))
	assert.NonNil(t, cleanup)

	data, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	assert.Equal(t, string(data), "stuff")

	cleanup()
	_, err = os.Stat(filename)
	assert.NonNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestWritePrefix(t *testing.T) {
	mockT := new(testing.T)

	filename, cleanup := Write(mockT, "stuff", "things")
	assert.False(t, mockT.Failed())
	assert.True(t, strings.Contains(filename, "things."))
	assert.NonNil(t, cleanup)

	data, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	assert.Equal(t, string(data), "stuff")

	cleanup()
	_, err = os.Stat(filename)
	assert.NonNil(t, err)
	assert.True(t, os.IsNotExist(err))
}
