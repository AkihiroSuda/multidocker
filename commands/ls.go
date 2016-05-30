package commands

import (
	"fmt"
	"log"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

// LS is the "ls" command.
var LS = cli.Command{
	Name:   "ls",
	Usage:  "List hosts",
	Action: ls,
}

func ls(c *cli.Context) error {
	w := c.App.Writer
	store := md.DefaultStore()
	hosts, err := store.Hosts()
	if err != nil {
		if hosts == nil {
			return err
		}
		// just warn, not an error
		log.Printf("error while scanning hosts: %v", err)
	}
	fmt.Fprintf(w, "NAME\tSTATE\tURL\n")
	for _, host := range hosts {
		name := host.Config().Name
		state, err := host.State()
		if err != nil {
			log.Printf("error while fetching state of host %s: %v",
				name, err)
			continue
		}
		url, err := host.URL()
		if err != nil {
			log.Printf("error while fetching url of host %s: %v",
				name, err)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			name, state, url)
	}
	return nil
}
