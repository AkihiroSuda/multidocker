package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

// arg0 is an ugly named utility function.
func arg0(c *cli.Context) (string, error) {
	name := "default"
	if c.NArg() == 0 {
		return name, nil
	} else if c.NArg() == 1 {
		return c.Args()[0], nil
	} else {
		return "", fmt.Errorf("too many arguments: %s", c.Args())
	}
}

func detectShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if os.Getenv("__fish_bin_dir") != "" {
			return "fish", nil
		}
		return "", errors.New("SHELL is not set")
	}
	return filepath.Base(shell), nil
}
