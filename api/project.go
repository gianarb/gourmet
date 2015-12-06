package api

import (
    "net/http"
	"github.com/gianarb/gourmet/runner"
	"encoding/json"
)

type StartBuildRequest struct {
	Source string
}

func ProjectHandler (runner runner.Runner) func(w http.ResponseWriter, r *http.Request) {
	return  func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var t StartBuildRequest
		decoder.Decode(&t)

		containerId, _ := runner.BuildContainer()

		runner.Exec(containerId, []string{"wget", t.Source})
		runner.Exec(containerId, []string{"unzip", "gourmet.zip", "-d", "."})
		runner.Exec(containerId, []string{"bin/console"})

		runner.RemoveContainer(containerId)
	}
}

