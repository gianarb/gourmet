package runner

type Runner interface {
	CreateFunc(img string, envVars []string, source string) (string, error)
	RunFunc(funcId string, envVars []string) (string, error)
	CommitContainer(containerId string) (string, error)
	PullImage(repository string) error
	DeleteFunc(funcId string) error
}
