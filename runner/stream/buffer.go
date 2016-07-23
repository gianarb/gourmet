package stream

import (
	"bytes"

	"github.com/mitchellh/colorstring"
)

type BufferStream struct {
	Buffer *bytes.Buffer
}

func (r BufferStream) Write(data []byte) (n int, err error) {
	colorstring.Printf("[blue]%s \n", string(data))
	return r.Buffer.Write(data)
}

func (r BufferStream) String() string {
	return r.Buffer.String()
}
