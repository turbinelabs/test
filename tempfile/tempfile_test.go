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
