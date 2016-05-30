// +build linux

package multidocker

import (
	"fmt"
	"os"
)

func findProcess(pid int) (*os.Process, error) {
	// On Unix systems, os.FindProcess always succeeds
	// and returns a Process for the given pid,
	// regardless of whether the process exists.
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return os.FindProcess(pid)
}
