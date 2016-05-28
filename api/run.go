package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gianarb/gourmet/runner"
	"github.com/gorilla/mux"
)

type RunResponse struct {
	Logs string
}

func RunHandler(runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseStruct := RunResponse{}
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id, ok := vars["id"]
		env := []string{}
		if !ok {
			w.WriteHeader(404)
			errStruct := Error{errors.New("id is required"), 4321}
			w.Write(errStruct.ToJson())
			return
		}
		cId, err := runner.BuildContainer(id, env, []string{"sleep", "1000"})
		if err != nil {
			w.WriteHeader(500)
			errStruct := Error{errors.New("id is required"), 4321}
			w.Write(errStruct.ToJson())
			return
		}
		runner.Exec(cId, []string{"bin/console"})
		runner.RemoveContainer(cId)
		responseStruct.Logs = runner.GetStream().String()

		json, _ := json.Marshal(responseStruct)
		w.WriteHeader(200)
		w.Write(json)
	}
}
