package runner

import(
)

type Runner interface {
	BuildContainer(img string, envVars []string) (string, error)
	Exec(containerId string, command []string) (error)
	RemoveContainer(containerId string) (error)
}
