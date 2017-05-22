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

// Package log provides functions for producing useful Loggers for testing
package log

import (
	"bytes"
	"log"

	"github.com/turbinelabs/test/io"
)

// NewNoopLogger produces a Logger that does nothing.
func NewNoopLogger() *log.Logger {
	return log.New(io.NewNoopWriter(), "", log.LstdFlags)
}

// NewBufferLogger produces a Logger that writes to a buffer, along with a
// pointer to the buffer to which it writes. This is useful in verifying that
// a Logger is used as expected. By default all Logger flags are disabled.
func NewBufferLogger() (*log.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	return log.New(buf, "", 0), buf
}

// NewChannelLogger produces a Logger that writes to the returned
// string channel. The channel is constructed with the given
// capacity. Logging will block if the channel is full. By default all
// Logger flags are disabled.
func NewChannelLogger(capacity int) (*log.Logger, <-chan string) {
	ch := make(chan string, capacity)
	w := io.NewChannelWriter(ch)

	return log.New(w, "", 0), ch
}
