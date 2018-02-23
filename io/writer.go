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

package io

import (
	"errors"
	"fmt"
	"io"
)

type noopWriter struct{}

// NewNoopWriter produces a writer that does nothing.
func NewNoopWriter() io.WriteCloser {
	return &noopWriter{}
}

func (_ *noopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (_ *noopWriter) Close() error {
	return nil
}

type channelWriter struct {
	ch chan<- string
}

type failingWriter struct{}

// NewFailingWriter produces a Writer that always fails.
func NewFailingWriter() io.WriteCloser {
	return &failingWriter{}
}

func (_ *failingWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("failed to write: %s", string(p))
}

func (_ *failingWriter) Close() error {
	return errors.New("failed to close")
}

// NewChannelWriter produces a writer that publishes to a string
// channel. Writes that occur when the channel is full will block.
func NewChannelWriter(ch chan<- string) io.WriteCloser {
	return &channelWriter{ch}
}

func (cw *channelWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("channel writer error: %+v", x)
		}
	}()

	if p == nil {
		cw.ch <- ""
		return 0, nil
	}

	cw.ch <- string(p)
	return len(p), nil
}

func (cw *channelWriter) Close() (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("channel writer close error: %+v", x)
		}
	}()

	close(cw.ch)
	return nil
}
