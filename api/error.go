package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	code    int    `json:"code"`
}

func (er *Error) ToJson() []byte {
	json, _ := json.Marshal(er)
	return json
}

func errorRender(statusCode int, intCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	errStruct := Error{
		Message: err.Error(),
		code:    intCode,
	}
	w.Write(errStruct.ToJson())
}
