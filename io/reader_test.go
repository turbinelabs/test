package io

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestFailingReaderAlwaysFails(t *testing.T) {
	f := NewFailingReader()
	for i := 0; i < 10; i++ {
		n, err := f.Read(make([]byte, 1))
		assert.Equal(t, n, 0)
		assert.NonNil(t, err)
		assert.Equal(t, err.Error(), FailingReaderMessage)
	}
}

func TestFailingReaderCloses(t *testing.T) {
	f := NewFailingReader()
	assert.Nil(t, f.Close())
}
