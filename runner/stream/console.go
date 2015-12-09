package stream

import "fmt"

type ConsoleStream struct {
}

func (r ConsoleStream) Write(data []byte) (n int, err error) {
	fmt.Printf("%s \n", string(data))
	return len(data), nil
}
