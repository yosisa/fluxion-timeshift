package timeshift

import (
	"bytes"
	"reflect"
	"sync"

	"github.com/ugorji/go/codec"
)

type MsgpackBuffer struct {
	buf *bytes.Buffer
	enc *codec.Encoder
	m   sync.RWMutex
}

func NewMsgpackBuffer() *MsgpackBuffer {
	b := &bytes.Buffer{}
	return &MsgpackBuffer{
		buf: b,
		enc: codec.NewEncoder(b, mh),
	}
}

func (b *MsgpackBuffer) Add(t int64, v interface{}) error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.enc.Encode(v)
}

func (b *MsgpackBuffer) Bytes() []byte {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.buf.Bytes()
}

var mh = &codec.MsgpackHandle{RawToString: true, WriteExt: true}

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}
