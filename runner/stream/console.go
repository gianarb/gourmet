package stream

import (
	"github.com/mitchellh/colorstring"
)

type ConsoleStream struct {
}

func (r ConsoleStream) Write(data []byte) (n int, err error) {
	colorstring.Printf("[blue]%s \n", string(data))
	return len(data), nil
}
