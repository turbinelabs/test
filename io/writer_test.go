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

func TestFailingWriterAlwaysFails(t *testing.T) {
	w := NewFailingWriter()
	for i := 0; i < 10; i++ {
		n, err := w.Write([]byte("something"))
		assert.Equal(t, n, 0)
		assert.NonNil(t, err)
	}
}

func TestFailingWriterFailesToClose(t *testing.T) {
	w := NewFailingWriter()
	assert.NonNil(t, w.Close())
}

func readChan(c <-chan string) (string, bool) {
	select {
	case s := <-c:
		return s, true
	default:
		return "", false
	}
}

func TestChannelWriter(t *testing.T) {
	c := make(chan string, 3)
	w := NewChannelWriter(c)

	n, err := w.Write([]byte("abc"))
	assert.Equal(t, n, 3)
	assert.Nil(t, err)

	n, err = w.Write([]byte("defghi"))
	assert.Equal(t, n, 6)
	assert.Nil(t, err)

	r, ok := readChan(c)
	assert.True(t, ok)
	assert.Equal(t, r, "abc")

	r, ok = readChan(c)
	assert.True(t, ok)
	assert.Equal(t, r, "defghi")

	_, ok = readChan(c)
	assert.False(t, ok)

	n, err = w.Write([]byte{})
	assert.Equal(t, n, 0)
	assert.Nil(t, err)

	n, err = w.Write(nil)
	assert.Equal(t, n, 0)
	assert.Nil(t, err)

	r, ok = readChan(c)
	assert.True(t, ok)
	assert.Equal(t, r, "")

	r, ok = readChan(c)
	assert.True(t, ok)
	assert.Equal(t, r, "")

	assert.Nil(t, w.Close())
}

func TestChannelWriterConvertsPanicsToErrors(t *testing.T) {
	c := make(chan string, 3)
	w := NewChannelWriter(c)
	close(c)

	n, err := w.Write([]byte("abc"))
	assert.Equal(t, n, 0)
	assert.NonNil(t, err)

	err = w.Close()
	assert.NonNil(t, err)
}
