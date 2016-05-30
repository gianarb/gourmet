package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gianarb/gourmet/runner"
	"github.com/gianarb/gourmet/runner/stream"
	"github.com/gorilla/mux"
)

type ContextMapper struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}
type RunRequest struct {
	Env []string
}

func RunHandler(runner runner.Runner, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
		logger.Printf("Container %s deleted", cId)

		writeResponseFromStdOutput(g, w)
	}
}

func writeResponseFromStdOutput(g *stream.BufferStream, w http.ResponseWriter) {
	context := ContextMapper{}
	err := json.Unmarshal([]byte(g.String()), &context)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%t", context)
	if context.StatusCode == 0 {
		context.StatusCode = 200
	}
	for h, v := range context.Headers {
		w.Header().Set(h, v)
	}
	w.WriteHeader(context.StatusCode)
	w.Write(context.Body)
}
