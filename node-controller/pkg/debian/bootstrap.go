package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Bootstrap uses debootstrap to bootstrap a Debian system
func (c *BootstrapConfig) Bootstrap() error {
	log.Infof("Bootstrapping Debian %s at %s", c.Suite, c.Target)
	debootstrap := DefaultDebootstrap
	if c.Debootstrap != "" {
		debootstrap = c.Debootstrap
	}

	debootstrap, err := exec.LookPath(debootstrap)
	if err != nil {
		return err
	}

	cmd := exec.Command(debootstrap, c.Suite, c.Target, c.Mirror)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// PostBootstrap runs post-bootstrap steps
func (c *BootstrapConfig) PostBootstrap() error {
	// Install ancillary packages
	if err := c.InstallBasePackages(); err != nil {
		return err
	}

	// Configure repositories
	if err := c.ConfigureRepositories(); err != nil {
		return err
	}

	// Install kernel
	if err := c.InstallKernel(); err != nil {
		return err
	}

	// Install bootloader
	if err := c.InstallBootloader(); err != nil {
		return err
	}

	return nil
}
