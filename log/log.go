package log

import (
	"log"

	"github.com/turbinelabs/test/io"
)

func NewNoopLogger() *log.Logger {
	return log.New(io.NewNoopWriter(), "", log.LstdFlags)
}
