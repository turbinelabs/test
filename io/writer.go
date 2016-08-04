package io

import (
	"io"
)

type noopWriter struct{}

func NewNoopWriter() io.WriteCloser {
	return &noopWriter{}
}

func (_ *noopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (_ *noopWriter) Close() error {
	return nil
}
