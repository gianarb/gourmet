package main

import (
	dockerRun "github.com/gianarb/gourmet/runner/docker"
	"github.com/mitchellh/cli"
	"github.com/gianarb/gourmet/command"
	"github.com/gianarb/gourmet/logger"
	"github.com/gianarb/gourmet/runner/stream"
	"github.com/fsouza/go-dockerclient"
	"log"
	"bytes"
	"os"
)

func main() {
	logger := log.New(&logger.Console{}, "", 1)
	c := cli.NewCLI("gourmet", "0.0.0")
    c.Args = os.Args[1:]

	client, err := docker.NewClientFromEnv()

	if err != nil {
		logger.Fatal(err)
	}

	b := new(bytes.Buffer)
	s := stream.BufferStream{b}
	dockerRunner := dockerRun.DockerRunner{client, s}

    c.Commands = map[string]cli.CommandFactory{
        "api": func() (cli.Command, error) {
			return &command.ApiCommand{&dockerRunner, logger}, nil;
		},
    }

    exitStatus, _ := c.Run()

    os.Exit(exitStatus)
}
