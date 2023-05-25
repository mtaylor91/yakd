package debian

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// Bootstrap uses debootstrap to bootstrap a Debian system
func (c *BootstrapConfig) Bootstrap() error {
	log.Infof("Bootstrapping Debian %s at %s", c.Suite, c.Target)
	debootstrap := DefaultDebootstrap
	if c.Debootstrap != "" {
		debootstrap = c.Debootstrap
	}

	err := executor.Default.RunCmd(debootstrap, c.Suite, c.Target, c.Mirror)
	if err != nil {
		return err
	}

	return nil
}

// PostBootstrap runs post-bootstrap steps
func (c *BootstrapConfig) PostBootstrap(chroot executor.Executor) error {
	// Configure locales
	if err := configureLocales(chroot, c.Target); err != nil {
		return err
	}

	// Install base packages
	if err := installBasePackages(chroot); err != nil {
		return err
	}

	// Configure repositories
	if err := c.configureRepositories(); err != nil {
		return err
	}

	// Install Kubernetes packages
	if err := installKubePackages(chroot); err != nil {
		return err
	}

	// Configure system for kubernetes
	if err := configureKubernetes(c.Target); err != nil {
		return err
	}

	// Configure the admin user
	if err := configureAdminUser(chroot); err != nil {
		return err
	}

	// Install kernel
	if err := c.installKernel(chroot); err != nil {
		return err
	}

	return nil
}
