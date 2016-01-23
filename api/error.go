package api

import (
	"encoding/json"
)

type Error struct {
	Message error
	code int
}

func (er *Error) ToJson() []byte {
	json, _ := json.Marshal(er)
	return json
}
