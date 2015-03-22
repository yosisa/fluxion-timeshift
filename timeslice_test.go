package timeshift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type String string

func (s *String) Add(t int64, v interface{}) error {
	*s += String(v.(string))
	return nil
}

func TestTimeSlice(t *testing.T) {
	ts := NewTimeSlice(1, func() interface{} {
		var s String
		return &s
	})
	s := int64(1426942000)
	ts.Add(s, "a")
	ts.Add(s, "a")
	ts.Add(s+10, "b")
	ts.Add(s+20, "c")

	vals := ts.Range(s-10, s+20)
	assert.Equal(t, 2, len(vals))
	assert.Equal(t, "aa", *vals[0].(*String))
	assert.Equal(t, "b", *vals[1].(*String))

	vals = ts.Range(s+5, s+30)
	assert.Equal(t, 2, len(vals))
	assert.Equal(t, "b", *vals[0].(*String))
	assert.Equal(t, "c", *vals[1].(*String))
}

func TestNestedTimeSlice(t *testing.T) {
	ts := NewTimeSlice(10, func() interface{} {
		return NewTimeSlice(1, func() interface{} {
			var s String
			return &s
		})
	})
	s := int64(1426942000)
	ts.Add(s, "a")
	ts.Add(s, "a")
	ts.Add(s+5, "b")
	ts.Add(s+10, "c")
	ts.Add(s+15, "d")
	ts.Add(s+20, "e")
	ts.Add(s+25, "f")

	vals := ts.Range(s-5, s+20)
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, "aa", *vals[0].(*String))
	assert.Equal(t, "b", *vals[1].(*String))
	assert.Equal(t, "c", *vals[2].(*String))
	assert.Equal(t, "d", *vals[3].(*String))

	vals = ts.Range(s+5, s+25)
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, "b", *vals[0].(*String))
	assert.Equal(t, "c", *vals[1].(*String))
	assert.Equal(t, "d", *vals[2].(*String))
	assert.Equal(t, "e", *vals[3].(*String))
}
