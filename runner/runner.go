package runner

type Runner interface {
	CreateFunc(img string, envVars []string, source string) (string, error)
	RunFunc(id string, envVars []string) error
	CommitContainer(containerId string) (string, error)
	PullImage(repository string) error
	RemoveContainer(containerId string) error
	DeleteImage(containerId string) error
}
