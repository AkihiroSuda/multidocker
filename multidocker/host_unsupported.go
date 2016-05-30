// +build !linux

package multidocker

// NewHost create a host with config. dir must be an absolute path.
func NewHost(config HostConfig, dir string) (Host, error) {
	return nil, unsupportedErr
}
