package timeshift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgpackBuffer(t *testing.T) {
	b := NewMsgpackBuffer()
	b.Add(0, "a")
	assert.Equal(t, []byte{0xa1, 0x61}, b.Bytes())
	b.Add(0, "b")
	assert.Equal(t, []byte{0xa1, 0x61, 0xa1, 0x62}, b.Bytes())
}
