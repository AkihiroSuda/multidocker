package multidocker

// Provisioner installs Docker Engine to a host.
type Provisioner interface {
	Install(Host) error
}
