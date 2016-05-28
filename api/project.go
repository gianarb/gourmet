package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gianarb/gourmet/runner"
)

type StartBuildRequest struct {
	Source string
	Img    string
}

type ProjectResponse struct {
	RunId string
	Logs  string
}

func ProjectHandler(runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseStruct := ProjectResponse{}
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		var t StartBuildRequest
		decoder.Decode(&t)

		containerId, err := runner.BuildContainer(t.Img, []string{})
		if err != nil {
			errorRender(500, 4311, err, w)
			return
		}
		logger.Printf("Build %s started - source %s", containerId[0:12], t.Img)
		runner.Exec(containerId, []string{"wget", t.Source})
		runner.Exec(containerId, []string{"unzip", "gourmet.zip", "-d", "."})
		image, err := runner.CommitContainer(containerId)
		if err != nil {
			logger.Printf("%s", err)
			errorRender(500, 4310, err, w)
			return
		}
		runner.RemoveContainer(containerId)

		logger.Printf("Container %s :: \n %s :: \n", containerId[0:12], runner.GetStream().String())

		responseStruct.Logs = runner.GetStream().String()
		logger.Printf("Build %s removed", containerId[0:12])
		responseStruct.RunId = image
		json, _ := json.Marshal(responseStruct)
		w.WriteHeader(200)
		w.Write(json)
	}
}
