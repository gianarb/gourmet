package docker

import (
	"io"
	"github.com/fsouza/go-dockerclient"
)

type DockerRunner struct {
	Docker *docker.Client
	Stream io.Writer
}

func (dr *DockerRunner) BuildContainer(img string, envVars []string) (string, error) {

	container, err := dr.Docker.CreateContainer(docker.CreateContainerOptions{
		"",
		&docker.Config{
			Image:        img,
			Cmd:          []string{"sleep", "1000"},
			WorkingDir:   "/tmp",
			AttachStdout: false,
			AttachStderr: false,
			Env:          envVars,
		},
		nil,
	})

	if err != nil {
		return "", err;
	}

	err = dr.Docker.StartContainer(
		container.ID,
		&docker.HostConfig{
			DNS:   []string{"8.8.8.8", "8.8.4.4"},
		},
	)
	if err != nil {
		return "", err;
	}
	return container.ID, nil;
}

func (dr *DockerRunner) Exec(containerId string, command []string) (error) {

	exec, err := dr.Docker.CreateExec(docker.CreateExecOptions{
		Container:    containerId,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          command,
	})

	if err != nil {
		return err;
	}

	err = dr.Docker.StartExec(exec.ID, docker.StartExecOptions{
		Detach:      false,
		Tty:         false,
		RawTerminal: true,
		OutputStream: dr.Stream,
		ErrorStream:  dr.Stream,
	})

	//inspect, err := dr.Docker.InspectExec(exec.ID)

	if err != nil {
		return err;
	}

	if err != nil {
		return err;
	}

	return nil
}

func (dr *DockerRunner) RemoveContainer(containerId string) error {
	err := dr.Docker.KillContainer(docker.KillContainerOptions{ID: containerId})
	err = dr.Docker.RemoveContainer(docker.RemoveContainerOptions{ID: containerId, RemoveVolumes: true})
	if(err != nil) {
		return err;
	}

	return nil
}
