package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/gianarb/gourmet/command"
	dockerRun "github.com/gianarb/gourmet/runner/docker"
	"github.com/mitchellh/cli"
)

func main() {
	logrus.Info("Gourmet started")
	c := cli.NewCLI("gourmet", "0.0.0")
	c.Args = os.Args[1:]

	client, err := docker.NewClientFromEnv()

	if err != nil {
		logrus.WithField("error", err).Warn("Problem to communicate with docker")
	}

	dockerRunner := dockerRun.DockerRunner{client}

	c.Commands = map[string]cli.CommandFactory{
		"api": func() (cli.Command, error) {
			return &command.ApiCommand{&dockerRunner}, nil
		},
	}

	exitStatus, _ := c.Run()

	os.Exit(exitStatus)
}
