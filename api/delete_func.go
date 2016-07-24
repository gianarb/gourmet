package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/gourmet/runner"
	"github.com/gorilla/mux"
)

func DeleteFuncHandler(runner runner.Runner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		vars := mux.Vars(r)
		imageName, ok := vars["id"]
		if !ok {
			errorRender(404, 4313, errors.New("id is required"), w)
			return
		}
		if os.Getenv("GOURMET_REGISTRY_URL") != "" {
			imageName = fmt.Sprintf("%s/%s", os.Getenv("GOURMET_REGISTRY_URL"), imageName)
		}
		logrus.Infof("Delete function %s", imageName)
		err := runner.DeleteFunc(imageName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warnf("Function not fund")
			errorRender(404, 4317, err, w)
			return
		}
		w.WriteHeader(204)
	}
}
