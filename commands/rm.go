package commands

import (
	"log"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

// RM is the "rm" command.
var RM = cli.Command{
	Name:   "rm",
	Usage:  "Remove a host",
	Action: rm,
}

func rm(c *cli.Context) error {
	name, err := arg0(c)
	if err != nil {
		return err
	}
	store := md.DefaultStore()
	host, err := store.Host(name)
	if err != nil {
		return err
	}

	log.Printf("Removing the host %s", name)
	if err = host.Remove(); host != nil {
		return err
	}
	return nil
}
