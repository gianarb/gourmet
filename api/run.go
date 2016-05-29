package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

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
		w.Header().Set("Content-Type", "application/text")
		decoder := json.NewDecoder(r.Body)
		var t RunRequest
		decoder.Decode(&t)
		vars := mux.Vars(r)
		imageName, ok := vars["id"]
		if !ok {
			errorRender(404, 4313, errors.New("id is required"), w)
			return
		}
		if os.Getenv("GOURMET_REGISTRY_URL") != "" {
			imageName = fmt.Sprintf("%s/%s", os.Getenv("GOURMET_REGISTRY_URL"), imageName)
		}
		cId, err := runner.BuildContainer(imageName, t.Env)
		if err != nil {
			err := runner.PullImage(imageName)
			if err != nil {
				errorRender(500, 4312, err, w)
				return
			}
			cId, err = runner.BuildContainer(imageName, t.Env)
			if err != nil {
				errorRender(500, 4317, err, w)
				return
			}
		}
		logger.Printf("Container %s created", cId)
		g, b, err := runner.Exec(cId, []string{"bin/console"})
		if err != nil {
			logger.Printf("Container %s :: \n %s :: \n", cId, b.String())
			logger.Printf("Container %s :: \n %s :: \n", cId, err)
		}
		logger.Printf("Container %s :: \n %s :: \n", cId, g.String())
		runner.RemoveContainer(cId)
		w.WriteHeader(200)
		logger.Printf("Container %s deleted", cId)
		w.Write([]byte(g.String()))
	}
}
