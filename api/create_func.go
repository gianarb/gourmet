package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/gourmet/runner"
)

type StartBuildRequest struct {
	Source string
	Img    string
}

type ProjectResponse struct {
	FuncId string `json:"funcId"`
}

func CreateFuncHandler(runner runner.Runner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseStruct := ProjectResponse{}
		w.Header().Set("Content-Type", "application/json")
		registry := os.Getenv("GOURMET_REGISTRY_URL")
		if registry == "" {
			registry = "hub.docker.com"
		}

		decoder := json.NewDecoder(r.Body)
		var t StartBuildRequest
		decoder.Decode(&t)

		logrus.Infof("Started new build %s", t.Img)

		containerId, err := runner.CreateFunc(t.Img, []string{}, t.Source)
		if err != nil {
			err := runner.PullImage(t.Img)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Warnf("We can not download %s from the registry %s", t.Img, registry)
				errorRender(500, 4312, err, w)
				return
			}
			containerId, err = runner.CreateFunc(fmt.Sprintf("%s/%s", registry, t.Img), []string{}, t.Source)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Warnf("We can not create a new container from this image %s", t.Img)
				errorRender(500, 4317, err, w)
				return
			}
		}

		logrus.WithFields(logrus.Fields{
			"container": containerId,
		}).Info("Build completed")

		image, err := runner.CommitContainer(containerId)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"container": containerId,
				"error":     err,
			}).Warnf("We can not push this container into the registry %s", registry)
			errorRender(500, 4310, err, w)
			return
		}
		logrus.WithFields(logrus.Fields{
			"container": containerId,
		}).Info("Container removed")
		responseStruct.FuncId = image
		json, _ := json.Marshal(responseStruct)
		w.WriteHeader(200)
		w.Write(json)
	}
}
