package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// Bootstrap uses debootstrap to bootstrap a Debian system
func (c *BootstrapConfig) Bootstrap(ctx context.Context) error {
	log.Infof("Bootstrapping Debian %s at %s", c.Suite, c.Target)
	debootstrap := DefaultDebootstrap
	if c.Debootstrap != "" {
		debootstrap = c.Debootstrap
	}

	err := executor.Default.RunCmd(ctx, debootstrap, c.Suite, c.Target, c.Mirror)
	if err != nil {
		return err
	}

	return nil
}

// PostBootstrap runs post-bootstrap steps
func (c *BootstrapConfig) PostBootstrap(
	ctx context.Context, chroot executor.Executor) error {

	// Configure locales
	if err := configureLocales(ctx, chroot, c.Target); err != nil {
		return err
	}

	// Install base packages
	if err := installBasePackages(ctx, chroot); err != nil {
		return err
	}

	// Configure repositories
	if err := c.configureRepositories(ctx); err != nil {
		return err
	}

	// Install Kubernetes packages
	if err := installKubePackages(ctx, chroot); err != nil {
		return err
	}

	// Configure system for kubernetes
	if err := configureKubernetes(ctx, chroot, c.Target); err != nil {
		return err
	}

	// Configure networking
	if err := configureNetworking(ctx, chroot, c.Target); err != nil {
		return err
	}

	// Configure the admin user
	if err := configureAdminUser(ctx, chroot); err != nil {
		return err
	}

	// Install kernel
	if err := c.installKernel(ctx, chroot); err != nil {
		return err
	}

	return nil
}
