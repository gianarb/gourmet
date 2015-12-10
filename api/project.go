package api

import (
    "net/http"
	"github.com/gianarb/gourmet/runner"
	"encoding/json"
	"log"
)

type StartBuildRequest struct {
	Source string
}

func ProjectHandler (runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return  func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t StartBuildRequest
		decoder.Decode(&t)

		containerId, err := runner.BuildContainer()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Printf("Build %s started - source %s", containerId[0:12])

		runner.Exec(containerId, []string{"wget", t.Source})
		runner.Exec(containerId, []string{"unzip", "gourmet.zip", "-d", "."})
		runner.Exec(containerId, []string{"bin/console"})

		runner.RemoveContainer(containerId)
		logger.Printf("Build %s finished and removed", containerId[0:12])
	}
}

