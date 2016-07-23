package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/gourmet/runner"
	"github.com/gorilla/mux"
)

type RunResponse struct {
	Logs string
}
type RunRequest struct {
	Env []string
}

func RunFuncHandler(runner runner.Runner) func(w http.ResponseWriter, r *http.Request) {
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
		logrus.Infof("Run function from image %s", imageName)
		cId, err := runner.BuildContainer(imageName, t.Env)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warnf("Function not fund")
			errorRender(500, 4317, err, w)
			return
		}
		logrus.WithFields(logrus.Fields{
			"container": cId,
		}).Infof("Running command %s", []string{"bin/console"})
		g, b, err := runner.Exec(cId, []string{"bin/console"})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"container": cId,
				"error":     err,
			}).Warn("Function failed")
			logrus.WithFields(logrus.Fields{
				"container": cId,
				"error":     err,
			}).Debug("Failed stream \n %s", b.String())
		}
		runner.RemoveContainer(cId)
		w.WriteHeader(200)
		logrus.WithFields(logrus.Fields{
			"container": cId,
		}).Info("Function completed")
		w.Write([]byte(g.String()))
	}
}
