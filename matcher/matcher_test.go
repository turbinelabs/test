package matcher

import (
	"bytes"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestAnyWriter(t *testing.T) {
	aw := AnyWriter{[]byte("yep")}
	buf := &bytes.Buffer{}

	assert.False(t, aw.Matches("nope"))

	assert.True(t, aw.Matches(buf))
	assert.Equal(t, buf.String(), "yep")

	assert.Equal(t, aw.String(), `AnyWriter("yep")`)
}
