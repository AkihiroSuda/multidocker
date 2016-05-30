package commands

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

const (
	fSetClientPath = "set-client-path"
	fShell         = "shell"
)

// Env is the "env" command.
var Env = cli.Command{
	Name:  "env",
	Usage: "Display the commands to set up the environment for the Docker client",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  fSetClientPath,
			Usage: "Set PATH for the Docker client",
		},
		cli.StringFlag{
			Name:  fShell,
			Value: "auto-detect",
			Usage: "Force environment to be configured for a specified shell: [fish, sh], default is auto-detect",
		},
	},
	Action: env,
}

func env(c *cli.Context) error {
	name, err := arg0(c)
	if err != nil {
		return err
	}
	shell := c.String(fShell)
	if shell == "auto-detect" {
		shell, err = detectShell()
		if err != nil {
			return err
		}
	}
	store := md.DefaultStore()
	host, err := store.Host(name)
	if err != nil {
		return err
	}
	s, err := envString(host, shell, c.Bool(fSetClientPath))
	if err != nil {
		return err
	}
	w := c.App.Writer
	fmt.Fprintf(w, "%s", s)
	return nil
}

func envString(host md.Host, shell string, setClientPath bool) (string, error) {
	url, err := host.URL()
	if err != nil {
		return "", err
	}
	path := ""
	if setClientPath {
		dir, err := host.Dir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, "bin")
	}

	switch shell {
	case "sh":
		return envStringBsh(url, path)
	case "bash":
		return envStringBsh(url, path)
	case "fish":
		return envStringFish(url, path)
	default:
		return "", fmt.Errorf("unsupported shell %s", shell)
	}
}

func envStringBsh(url, path string) (string, error) {
	s := fmt.Sprintf("DOCKER_HOST=%s\n", url)
	s += "export DOCKER_HOST\n"
	if path != "" {
		s += fmt.Sprintf("PATH=%s:$PATH\n", path)
		s += "export PATH\n"
	}
	return s, nil
}

func envStringFish(url, path string) (string, error) {
	s := fmt.Sprintf("set -gx DOCKER_HOST %s;\n", url)
	if path != "" {
		s += fmt.Sprintf("set -gx PATH %s $PATH;\n", path)
	}
	return s, nil
}
