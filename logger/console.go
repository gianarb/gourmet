package logger

import "log"

type Console struct {
}

func (r Console) Write(data []byte) (n int, err error) {
	log.Printf("%s", string(data))
	return len(data), nil
}
