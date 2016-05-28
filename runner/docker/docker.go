package docker

import (
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
	container, err := dr.Docker.InspectContainer(id)
	container.Config.Cmd = []string{"bin/console"}
	container.Config.Entrypoint = []string{"bin/console"}
	if err != nil {
		return "", err
	}
	opt := docker.CommitContainerOptions{
		Container: container.ID,
		Run:       container.Config,
	}
	img, err := dr.Docker.CommitContainer(opt)
	if err != nil {
		return "", err
	}
	return img.ID, nil
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
