package runner

import(
)

type Runner interface {
	BuildContainer() (string, error)
	Exec(containerId string, command []string) (error)
	RemoveContainer(containerId string) (error)
}
