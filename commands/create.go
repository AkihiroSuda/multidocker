package commands

import (
	"log"

	"gopkg.in/urfave/cli.v1"

	md "github.com/AkihiroSuda/multidocker/multidocker"
)

const (
	fEngineInstallBundle = "engine-install-bundle"
	fEngineStorageDriver = "engine-storage-driver"
)

// Create is the "create" command.
var Create = cli.Command{
	Name:  "create",
	Usage: "Create a host",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  fEngineInstallBundle,
			Value: "/go/src/github.com/docker/docker/bundles/latest",
			Usage: "Custom bundle path to use for engine installation",
		},
		cli.StringFlag{
			Name:  fEngineStorageDriver,
			Value: "overlay",
			Usage: "Custom bundle path to use for engine installation",
		},
	},
	Action: create,
}

func create(c *cli.Context) error {
	name, err := arg0(c)
	if err != nil {
		return err
	}
	hostConfig := md.HostConfig{
		Name:                name,
		EngineStorageDriver: c.String(fEngineStorageDriver),
	}
	store := md.DefaultStore()
	hostDir, err := store.AllocateHostDir(name)
	if err != nil {
		return err
	}
	log.Printf("Allocated directory %s for %s ",
		hostDir, name)
	host, err := md.NewHost(hostConfig, hostDir)
	if err != nil {
		return err
	}

	bundle := c.String(fEngineInstallBundle)
	log.Printf("Provisioning the bundle %s", bundle)
	prov, err := md.NewBundleProvisioner(bundle)
	if err != nil {
		return err
	}
	if err = prov.Install(host); err != nil {
		return err
	}

	log.Printf("Registering the host %s", name)
	if err = store.RegisterHost(host); err != nil {
		return err
	}

	log.Printf("Creating the host")
	if err = host.Create(); host != nil {
		return err
	}

	log.Printf("Starting the host")
	if err = host.Start(); host != nil {
		return err
	}
	return nil
}
