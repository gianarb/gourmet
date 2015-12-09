package main

import (
	dockerRun "github.com/gianarb/gourmet/runner/docker"
	"github.com/mitchellh/cli"
	"github.com/gianarb/gourmet/command"
	"github.com/gianarb/gourmet/runner/stream"
	"github.com/fsouza/go-dockerclient"
	"os"
)

func main() {
	c := cli.NewCLI("gourmet", "0.0.0")
    c.Args = os.Args[1:]

	client, _ := docker.NewClientFromEnv()

	dockerRunner := dockerRun.DockerRunner{client, stream.ConsoleStream{}}

    c.Commands = map[string]cli.CommandFactory{
        "api": func() (cli.Command, error) {
			return &command.ApiCommand{&dockerRunner}, nil;
		},
    }

    exitStatus, _ := c.Run()

    os.Exit(exitStatus)
}
