package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/os/common"
	"github.com/mtaylor91/yakd/pkg/system"
)

type BootstrapConfig struct {
	Suite       string
	Mirror      string
	Target      string
	Debootstrap string
}

// Bootstrap uses debootstrap to bootstrap a Debian system
func (c *BootstrapConfig) Bootstrap(ctx context.Context) error {
	log.Infof("Bootstrapping Debian %s at %s", c.Suite, c.Target)
	debootstrap := DefaultDebootstrap
	if c.Debootstrap != "" {
		debootstrap = c.Debootstrap
	}

	sys := system.Local.WithContext(ctx)
	err := sys.RunCommand(debootstrap, c.Suite, c.Target, c.Mirror)
	if err != nil {
		return err
	}

	return nil
}

// PostBootstrap runs post-bootstrap steps
func (c *BootstrapConfig) PostBootstrap(
	ctx context.Context, chroot system.System) error {

	// Configure locales
	if err := configureLocales(chroot, c.Target); err != nil {
		return err
	}

	// Install base packages
	if err := installBasePackages(chroot); err != nil {
		return err
	}

	// Configure repositories
	if err := c.configureRepositories(ctx); err != nil {
		return err
	}

	// Install Kubernetes packages
	if err := installKubePackages(chroot); err != nil {
		return err
	}

	// Configure system for kubernetes
	if err := common.ConfigureKubernetes(chroot, c.Target); err != nil {
		return err
	}

	// Configure networking
	if err := common.ConfigureNetwork(chroot, c.Target); err != nil {
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
