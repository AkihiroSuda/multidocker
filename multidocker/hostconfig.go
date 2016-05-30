package multidocker

const (
	// HostConfigJSONFile is the basename for the HostConfig JSON.
	HostConfigJSONFile = "config.v0.json"
)

// HostConfig is a config JSON for a host.
type HostConfig struct {
	// name of the host. e.g. "foobar"
	Name string `json:"name"`

	// storage driver
	EngineStorageDriver string `json:"engine_storage_driver"`
}
