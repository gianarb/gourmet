package command

import (
	"net/http"
	"fmt"
    "github.com/gorilla/mux"
	"strings"
	"flag"
	"github.com/gianarb/gourmet/runner"
	"github.com/gianarb/gourmet/api"
)

type ApiCommand struct {
	Runner runner.Runner
}

func (c *ApiCommand) Run(args []string) int {
	var port string

	cmdFlags := flag.NewFlagSet("event", flag.ContinueOnError)
	cmdFlags.StringVar(&port, "port", ":8080", "port")
	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()
    r.HandleFunc("/project", api.ProjectHandler(c.Runner)).Methods("POST")
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
