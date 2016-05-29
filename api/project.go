package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gianarb/gourmet/runner"
)

type StartBuildRequest struct {
	Source string
	Img    string
}

type ProjectResponse struct {
	RunId string
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
			err := runner.PullImage(t.Img)
			if err != nil {
				errorRender(500, 4312, err, w)
				return
			}
			containerId, err = runner.BuildContainer(fmt.Sprintf("%s/%s", os.Getenv("GOURMET_REGISTRY_URL"), t.Img), []string{})
			if err != nil {
				errorRender(500, 4317, err, w)
				runner.RemoveContainer(containerId)
				return
			}
		}
		logger.Printf("Build %s started - source %s", containerId, t.Img)
		g, b, err := runner.Exec(containerId, []string{"wget", t.Source})
		logger.Printf("Container %s :: \n %s :: \n", containerId, g.String())
		if err != nil {
			logger.Printf("Container %s :: \n %s :: \n", containerId, b.String())
			logger.Printf("Container %s :: \n %s :: \n", containerId, err)
		}
		g, b, err = runner.Exec(containerId, []string{"unzip", "gourmet.zip", "-d", "."})
		if err != nil {
			logger.Printf("Container %s :: \n %s :: \n", containerId, b.String())
			logger.Printf("Container %s :: \n %s :: \n", containerId, err)
		}
		logger.Printf("Container %s :: \n %s :: \n", containerId, g.String())
		image, err := runner.CommitContainer(containerId)
		if err != nil {
			logger.Printf("%s", err)
			errorRender(500, 4310, err, w)
			runner.RemoveContainer(containerId)
			return
		}
		runner.RemoveContainer(containerId)
		logger.Printf("Build %s removed", containerId)
		responseStruct.RunId = image
		json, _ := json.Marshal(responseStruct)
		w.WriteHeader(200)
		w.Write(json)
	}
}
