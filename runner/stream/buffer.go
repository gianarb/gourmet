package stream

import (
	"bytes"
)

type BufferStream struct {
	Buffer *bytes.Buffer
}

func (r BufferStream) Write(data []byte) (n int, err error) {
	return r.Buffer.Write(data)
}

func (r BufferStream) String() (string) {
	return r.Buffer.String()
}
