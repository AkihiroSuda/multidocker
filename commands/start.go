package commands

import (
	"log"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

// Start is the "start" command.
var Start = cli.Command{
	Name:   "start",
	Usage:  "Start a host",
	Action: start,
}

func start(c *cli.Context) error {
	name, err := arg0(c)
	if err != nil {
		return err
	}
	store := md.DefaultStore()
	host, err := store.Host(name)
	if err != nil {
		return err
	}

	log.Printf("Starting the host %s", name)
	if err = host.Start(); host != nil {
		return err
	}
	return nil
}
