// +build !linux

package multidocker

// NewBundleProvisioner instantiates a provisioner for the bundle.
func NewBundleProvisioner(bundlePath string) (Provisioner, error) {
	return nil, unsupportedErr
}
