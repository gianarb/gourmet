package logger

import "fmt"

type Console struct {
}

func (r Console) Write(data []byte) (n int, err error) {
	fmt.Print(string(data))
	return len(data), nil
}
