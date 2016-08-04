package io

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNoopWriterAlwaysSucceeds(t *testing.T) {
	w := NewNoopWriter()
	for i := 0; i < 10; i++ {
		n, err := w.Write([]byte("something"))
		assert.Equal(t, n, 9)
		assert.Nil(t, err)
	}
}

func TestNoopWriterCloses(t *testing.T) {
	w := NewNoopWriter()
	assert.Nil(t, w.Close())
}
