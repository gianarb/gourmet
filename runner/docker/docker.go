package docker

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/gianarb/gourmet/runner/stream"
)

type DockerRunner struct {
	Docker *docker.Client
	Stream stream.BufferStream
}

func (dr *DockerRunner) GetStream() stream.BufferStream {
	return dr.Stream
}

func (dr *DockerRunner) CommitContainer(id string) (string, error) {
	id, err := dr.commit(id)
	if err != nil {
		return "", err
	}
	if os.Getenv("GOURMET_REGISTRY_URL") != "" {
		u, _ := url.Parse(os.Getenv("GOURMET_REGISTRY_URL"))
		name := fmt.Sprintf("%s/%s", u, id)
		err = dr.push(name)
		if err != nil {
			return "", err
		}
	}
	return id, nil
}

func (dr *DockerRunner) BuildContainer(img string, envVars []string) (string, error) {
	container, err := dr.Docker.CreateContainer(docker.CreateContainerOptions{
		Name: "",
		Config: &docker.Config{
			Cmd:          []string{"sleep", "1000"},
			Image:        img,
			WorkingDir:   "/root",
			AttachStdout: false,
			AttachStderr: false,
			Env:          envVars,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
	})
	if err != nil {
		return "", err
	}
	err = dr.Docker.StartContainer(
		container.ID,
		&docker.HostConfig{
			DNS: []string{"8.8.8.8", "8.8.4.4"},
		},
	)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func (dr *DockerRunner) Exec(containerId string, command []string) error {
	exec, err := dr.Docker.CreateExec(docker.CreateExecOptions{
		Container:    containerId,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          command,
	})
	if err != nil {
		return err
	}
	err = dr.Docker.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		Tty:          false,
		RawTerminal:  true,
		OutputStream: dr.Stream,
		ErrorStream:  dr.Stream,
	})
	if err != nil {
		return err
	}
	return nil
}

func (dr *DockerRunner) RemoveContainer(containerId string) error {
	err := dr.Docker.KillContainer(docker.KillContainerOptions{ID: containerId})
	err = dr.Docker.RemoveContainer(docker.RemoveContainerOptions{ID: containerId, RemoveVolumes: true})
	if err != nil {
		return err
	}

	return nil
}

func (dr *DockerRunner) push(name string) error {
	opt := docker.PushImageOptions{
		Name: name,
		Tag:  "latest",
	}
	auth := docker.AuthConfiguration{}
	err := dr.Docker.PushImage(opt, auth)
	return err
}

func (dr *DockerRunner) commit(id string) (string, error) {
	container, err := dr.Docker.InspectContainer(id)
	container.Config.Cmd = []string{"bin/console"}
	container.Config.Entrypoint = []string{"bin/console"}
	name := fmt.Sprintf("gourmet/%s", randStringRunes(10))
	if err != nil {
		return "", err
	}
	opt := docker.CommitContainerOptions{
		Container:  container.ID,
		Repository: name,
		Run:        container.Config,
	}
	_, err = dr.Docker.CommitContainer(opt)
	if err != nil {
		return "", err
	}
	return name, nil
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
