package runner

import "github.com/gianarb/gourmet/runner/stream"

type Runner interface {
	BuildContainer(img string, envVars []string) (string, error)
	Exec(containerId string, command []string) (*stream.BufferStream, *stream.BufferStream, error)
	CommitContainer(containerId string) (string, error)
	PullImage(repository string) error
	RemoveContainer(containerId string) error
	DeleteImage(containerId string) error
}
