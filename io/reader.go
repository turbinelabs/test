package io

import (
	"errors"
	"io"
)

const FailingReaderMessage = "failingReader error"

type failingReader struct{}

func NewFailingReader() io.ReadCloser {
	return &failingReader{}
}

func (r *failingReader) Read(p []byte) (int, error) {
	return 0, errors.New(FailingReaderMessage)
}

func (r *failingReader) Close() error {
	return nil
}
