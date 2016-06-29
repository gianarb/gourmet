package docker

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/gianarb/gourmet/runner/stream"
)

type DockerRunner struct {
	Docker *docker.Client
}

func (dr *DockerRunner) CommitContainer(id string) (string, error) {
	imageId, err := dr.commit(id)
	if err != nil {
		return "", err
	}
	if os.Getenv("GOURMET_REGISTRY_URL") != "" {
		u, _ := url.Parse(os.Getenv("GOURMET_REGISTRY_URL"))
		err = dr.tag(imageId, u)
		if err != nil {
			return "", err
		}
		err = dr.push(fmt.Sprintf("%s/%s", u, imageId))
		if err != nil {
			return "", err
		}
	}
	return imageId, nil
}

func (dr *DockerRunner) PullImage(id string) error {
	err := dr.pull(id)
	if err != nil {
		return err
	}
	return nil
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

func (dr *DockerRunner) Exec(containerId string, command []string) (*stream.BufferStream, *stream.BufferStream, error) {
	oStream := stream.BufferStream{&bytes.Buffer{}}
	eStream := stream.BufferStream{&bytes.Buffer{}}
	exec, err := dr.Docker.CreateExec(docker.CreateExecOptions{
		Container:    containerId,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          command,
	})
	if err != nil {
		return nil, nil, err
	}
	err = dr.Docker.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		Tty:          false,
		RawTerminal:  true,
		OutputStream: oStream,
		ErrorStream:  eStream,
	})
	if err != nil {
		return nil, &eStream, err
	}
	return &oStream, nil, nil
}

func (dr *DockerRunner) DeleteImage(name string) error {
	err := dr.Docker.RemoveImageExtended(name, docker.RemoveImageOptions{
		Force: true,
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

func (dr *DockerRunner) pull(name string) error {
	opt := docker.PullImageOptions{
		Repository: name,
		Tag:        "latest",
	}
	auth := docker.AuthConfiguration{}
	err := dr.Docker.PullImage(opt, auth)
	return err
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
	name := fmt.Sprintf("%s", randStringRunes(10))
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

func (dr *DockerRunner) tag(image string, u *url.URL) error {
	opt := docker.TagImageOptions{
		Repo: fmt.Sprintf("%s/%s", u, image),
		Tag:  "latest",
	}
	return dr.Docker.TagImage(image, opt)
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
