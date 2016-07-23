package command

import (
	"flag"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/gourmet/api"
	"github.com/gianarb/gourmet/runner"
	"github.com/gorilla/mux"
)

type ApiCommand struct {
	Runner runner.Runner
}

func (c *ApiCommand) Run(args []string) int {
	var port string
	cmdFlags := flag.NewFlagSet("event", flag.ContinueOnError)
	cmdFlags.StringVar(&port, "port", ":8000", "port")
	if err := cmdFlags.Parse(args); err != nil {
		logrus.WithField("error", err).Warn("Problem to parse arguments")
	}
	logrus.Infof("API Server run on port %s", port)
	r := mux.NewRouter()
	r.HandleFunc("/func", api.CreateFuncHandler(c.Runner)).Methods("POST")
	r.HandleFunc("/func/{id}", api.RunFuncHandler(c.Runner)).Methods("POST")
	r.HandleFunc("/func/{id}", api.DeleteFuncHandler(c.Runner)).Methods("DELETE")
	r.HandleFunc("/ping", api.PingHandler()).Methods("GET")
	http.ListenAndServe(port, r)
	return 0
}

func (c *ApiCommand) Help() string {
	helpText := `
Usage: start gourmet API handler.
Options:
	-port=:8000			Servert port
`
	return strings.TrimSpace(helpText)
}

func (r *ApiCommand) Synopsis() string {
	return "prepare docker container for testing purposes"
}
