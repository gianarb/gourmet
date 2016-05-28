package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message error
	code    int
}

func (er *Error) ToJson() []byte {
	json, _ := json.Marshal(er)
	return json
}

func errorRender(statusCode int, intCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	errStruct := Error{err, intCode}
	w.Write(errStruct.ToJson())
}
