package command

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gianarb/gourmet/api"
	"github.com/gianarb/gourmet/runner"
	"github.com/gorilla/mux"
)

type ApiCommand struct {
	Runner runner.Runner
	Logger *log.Logger
}

func (c *ApiCommand) Run(args []string) int {
	var port string
	cmdFlags := flag.NewFlagSet("event", flag.ContinueOnError)
	cmdFlags.StringVar(&port, "port", ":8000", "port")
	if err := cmdFlags.Parse(args); err != nil {
		c.Logger.Fatal(err)
	}
	c.Logger.Print("API Server run on port ", port)
	r := mux.NewRouter()
	r.HandleFunc("/project", api.ProjectHandler(c.Runner, c.Logger)).Methods("POST")
	r.HandleFunc("/run/{id}", api.RunHandler(c.Runner, c.Logger)).Methods("POST")
	r.HandleFunc("/ping", api.PingHandler(c.Logger)).Methods("GET")
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
