package api

import (
    "net/http"
	"encoding/json"
	"github.com/gianarb/gourmet/runner"
	"log"
)

type StartBuildRequest struct {
	Source string
	Img string
	Env []string
}

type ProjectResponse struct {
	ContainerId string
	Logs	string
}

func ProjectHandler (runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return  func(w http.ResponseWriter, r *http.Request) {
		responseStruct := ProjectResponse{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		decoder := json.NewDecoder(r.Body)
		var t StartBuildRequest
		decoder.Decode(&t)

		containerId, err := runner.BuildContainer(t.Img, t.Env)
		responseStruct.ContainerId = containerId
		if err != nil {
			errStruct := Error{err, 4321}
			logger.Fatal(err)
			w.Write(errStruct.ToJson())
			http.Error(w, http.StatusText(500), 500)
		}
		logger.Printf("Build %s started - source %s", containerId[0:12], t.Img)

		runner.Exec(containerId, []string{"wget", t.Source})
		runner.Exec(containerId, []string{"unzip", "gourmet.zip", "-d", "."})
		runner.Exec(containerId, []string{"bin/console"})

		runner.RemoveContainer(containerId)

		logger.Printf("Container %s :: \n %s :: \n", containerId[0:12], runner.GetStream().String())
		responseStruct.Logs = runner.GetStream().String()

		logger.Printf("Build %s removed", containerId[0:12])
		json, _ := json.Marshal(responseStruct)
		w.Write(json)
	}
}

