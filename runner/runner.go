package runner

import (
	"github.com/gianarb/gourmet/runner/stream"
)

type Runner interface {
	BuildContainer(img string, envVars []string, cmd []string) (string, error)
	Exec(containerId string, command []string) error
	CommitContainer(containerId string) (string, error)
	RemoveContainer(containerId string) error
	GetStream() stream.BufferStream
}
