package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/AkihiroSuda/multidocker/commands"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "multidocker"
	app.Usage = "Multiple Docker daemons on a single machine"
	app.Commands = []cli.Command{
		commands.Create,
		commands.LS,
		commands.Start,
		commands.Stop,
		commands.RM,
		commands.Env,
	}
	return app
}

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
