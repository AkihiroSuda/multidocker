// +build linux

package multidocker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type bundleProvisioner struct {
	bundlePath string
}

func (p *bundleProvisioner) Install(host Host) error {
	hostDir, err := host.Dir()
	if err != nil {
		return err
	}
	dst := filepath.Join(hostDir, "bin")
	if err = os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	// TODO: support dynbinary-daemon
	srcs := []string{
		// NOTE: filepath.Join() does not work well with "/."
		// for "/.", please refer to the manpage of cp(1).
		filepath.Join(p.bundlePath, "binary-daemon") + "/.",
		filepath.Join(p.bundlePath, "binary-client") + "/.",
	}
	for _, src := range srcs {
		log.Printf("Copying %s to %s", src, dst)
		// FIXME: poor man's file copier
		out, err := exec.Command("cp", "-r", src, dst).CombinedOutput()
		outstr := string(out)
		if outstr != "" {
			log.Printf("%s", outstr)
		}
		if err != nil {
			return fmt.Errorf("error while copying %s to %s: %v(%s)",
				src, dst, err, outstr)
		}

	}

	return nil
}

// NewBundleProvisioner instantiates a provisioner for the bundle.
func NewBundleProvisioner(bundlePath string) (Provisioner, error) {
	return &bundleProvisioner{bundlePath}, nil
}
