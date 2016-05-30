package multidocker

// Host is similar to Docker Machine's Driver.
type Host interface {
	Config() HostConfig
	Dir() (string, error)
	Create() error
	URL() (string, error)
	State() (State, error)
	Remove() error
	Start() error
	Stop() error
}
