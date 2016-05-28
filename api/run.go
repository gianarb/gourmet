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
type RunRequest struct {
	Env []string
}

func RunHandler(runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseStruct := RunResponse{}
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var t RunRequest
		decoder.Decode(&t)
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(404)
			errStruct := Error{errors.New("id is required"), 4321}
			w.Write(errStruct.ToJson())
			return
		}
		cId, err := runner.BuildContainer(id, t.Env)
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
