// +build linux

package multidocker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type host struct {
	config HostConfig
	dir    string
}

// NewHost create a host with config. dir must be an absolute path.
func NewHost(config HostConfig, dir string) (Host, error) {
	if !filepath.IsAbs(dir) {
		return nil, fmt.Errorf("%s it not absolute", dir)
	}
	return &host{config, dir}, nil
}

// storeBinPath is equivalent to /usr/bin
func (h *host) storeBinPath() string {
	return filepath.Join(h.dir, "bin")
}

// storeSocketPath is equivalent to /var/run/docker.sock
func (h *host) storeSocketPath() string {
	return filepath.Join(h.dir, "socket")
}

// storePIDPath is equivalent to /var/run/docker.pid
func (h *host) storePIDPath() string {
	return filepath.Join(h.dir, "pid")
}

// storeGraphPath is equivalent to /var/lib/docker
func (h *host) storeGraphPath() string {
	return filepath.Join(h.dir, "graph")
}

// storeRunPath is equivalent to /var/run/docker
func (h *host) storeRunPath() string {
	return filepath.Join(h.dir, "run")
}

// storeStdoutPath is stdout log for dockerd
func (h *host) storeStdoutPath() string {
	return filepath.Join(h.dir, "stdout")
}

// storeStdoutPath is stdout log for dockerd
func (h *host) storeStderrPath() string {
	return filepath.Join(h.dir, "stderr")
}

// pid is -1 if the daemon is not started or unknown
func (h *host) pid() int {
	pidPath := h.storePIDPath()
	pidStrBytes, err := ioutil.ReadFile(pidPath)
	if err != nil {
		// this error is expected
		// because the file does not show up until Start()
		return -1
	}
	pid, err := strconv.Atoi(string(pidStrBytes))
	if err != nil {
		err = fmt.Errorf("contaminated file? %s: %v (content=%s)",
			pidPath, err, pidStrBytes)
		panic(err)
	}
	return pid
}

func (h *host) Config() HostConfig {
	return h.config
}

func (h *host) Dir() (string, error) {
	return h.dir, nil
}

func (h *host) Create() error {
	// NOTE: bin should have been already provisioned by the caller.
	log.Printf("Creating a directory %s", h.storeGraphPath())
	if err := os.MkdirAll(h.storeGraphPath(), 0755); err != nil {
		return err
	}

	log.Printf("Creating a directory %s", h.storeRunPath())
	if err := os.MkdirAll(h.storeRunPath(), 0755); err != nil {
		return err
	}
	return h.Start()
}

func (h *host) URL() (string, error) {
	return "unix://" + h.storeSocketPath(), nil
}

func (h *host) State() (State, error) {
	pid := h.pid()
	if pid < 0 {
		return Stopped, nil
	}
	_, err := findProcess(pid)
	if err != nil {
		log.Printf("cannot find process %d: %v", pid, err)
		return Stopped, nil
	}
	return Running, nil
}

func (h *host) Remove() error {
	state, err := h.State()
	if err != nil {
		return err
	}
	if state != Stopped {
		return fmt.Errorf("the host %s is not stopped (%s)",
			h.Config().Name, state)
	}
	dir, err := h.Dir()
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func (h *host) dockerd() (string, error) {
	// TODO: support old `docker daemon` (docker <= 1.11)
	return filepath.Join(h.storeBinPath(), "dockerd"), nil
}

func (h *host) dockerdArgs() ([]string, error) {
	url, err := h.URL()
	if err != nil {
		return nil, err
	}
	// TODO: support other options
	args := []string{
		"-D",
		"-H", url,
		"-p", h.storePIDPath(),
		"--iptables=false",
		"--ip-masq=false",
		"--bridge=none",
		"--graph=" + h.storeGraphPath(),
		"--exec-root=" + h.storeRunPath(),
		"-s", h.config.EngineStorageDriver,
	}
	return args, nil
}

func (h *host) dockerdEnv() ([]string, error) {
	envPATH := fmt.Sprintf("PATH=%s:%s",
		h.storeBinPath(), os.Getenv("PATH"))
	env := append([]string{envPATH}, os.Environ()...)
	return env, nil
}

func (h *host) dockerdCmd() (*exec.Cmd, error) {
	dockerd, err := h.dockerd()
	if err != nil {
		return nil, err
	}
	args, err := h.dockerdArgs()
	if err != nil {
		return nil, err
	}
	env, err := h.dockerdEnv()
	if err != nil {
		return nil, err
	}

	cmd := &exec.Cmd{
		Path:        dockerd,
		Args:        args,
		Env:         env,
		SysProcAttr: &syscall.SysProcAttr{Setpgid: true},
	}
	return cmd, nil
}

func (h *host) Start() error {
	stdout, err := os.Create(h.storeStdoutPath())
	if err != nil {
		return err
	}
	defer stdout.Close()
	stderr, err := os.Create(h.storeStderrPath())
	if err != nil {
		return err
	}
	defer stderr.Close()

	cmd, err := h.dockerdCmd()
	if err != nil {
		return err
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	log.Printf("Starting process %s (%s)", cmd.Path, cmd.Args)
	return cmd.Start()
}

func (h *host) Stop() error {
	pid := h.pid()
	if pid < 0 {
		return errors.New("the machine seems not started")
	}
	return syscall.Kill(pid, syscall.SIGTERM)
}
