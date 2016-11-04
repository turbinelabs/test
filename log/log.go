package log

import (
	"bytes"
	"log"

	"github.com/turbinelabs/test/io"
)

func NewNoopLogger() *log.Logger {
	return log.New(io.NewNoopWriter(), "", log.LstdFlags)
}

func NewBufferLogger() (*log.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	return log.New(buf, "", log.LstdFlags), buf
}
