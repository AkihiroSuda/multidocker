package multidocker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Store is the root directory for a multidocker environment.
type Store interface {
	StorePath() string
	Host(string) (Host, error)
	Hosts() ([]Host, error)
	AllocateHostDir(string) (string, error)
	RegisterHost(Host) error
}

type store struct {
	storePath string
}

// StorePath returns the directory for the store.
func (s *store) StorePath() string {
	return s.storePath
}

// Host returns a host named name.
func (s *store) Host(name string) (Host, error) {
	var config HostConfig
	dir := filepath.Join(s.StorePath(), name)
	jsonPath := filepath.Join(dir, HostConfigJSONFile)
	bytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return NewHost(config, dir)
}

// Hosts collects hosts.
// Even when error is non-nil, it may return some non-nil []Host.
func (s *store) Hosts() ([]Host, error) {
	dirs, err := ioutil.ReadDir(s.StorePath())
	if err != nil {
		return nil, err
	}
	var hosts []Host
	var errors []error
	for _, dir := range dirs {
		host, err := s.Host(dir.Name())
		if host != nil {
			hosts = append(hosts, host)
		}
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return hosts, fmt.Errorf("%s", errors)
	}
	return hosts, nil
}

// AllocateHostDir allocates a directory for the host named name.
func (s *store) AllocateHostDir(name string) (string, error) {
	path := filepath.Join(s.StorePath(), name)
	err := os.MkdirAll(path, 0755)
	return path, err
}

// RegisterHost registers host to store.
func (s *store) RegisterHost(host Host) error {
	config := host.Config()
	dir, err := host.Dir()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	jsonPath := filepath.Join(dir, HostConfigJSONFile)
	return ioutil.WriteFile(jsonPath, bytes, 0644)
}

// DefaultStore uses /var/lib/multidocker as the store.
func DefaultStore() Store {
	return &store{"/var/lib/multidocker"}
}
