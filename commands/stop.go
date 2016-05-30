package commands

import (
	"log"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

// Stop is the "stop" command.
var Stop = cli.Command{
	Name:   "stop",
	Usage:  "Stop a host",
	Action: stop,
}

func stop(c *cli.Context) error {
	name, err := arg0(c)
	if err != nil {
		return err
	}
	store := md.DefaultStore()
	host, err := store.Host(name)
	if err != nil {
		return err
	}

	log.Printf("Stopping the host %s", name)
	if err = host.Stop(); host != nil {
		return err
	}
	return nil
}
