package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/build/release/common"
	"github.com/mtaylor91/yakd/pkg/system"
)

type BootstrapInstaller struct {
	Suite  string
	Mirror string
	Target string
}

// Bootstrap uses debootstrap to bootstrap a Debian system
func (b *BootstrapInstaller) Bootstrap(ctx context.Context) error {
	log.Infof("Bootstrapping Debian %s at %s", b.Suite, b.Target)
	sys := system.Local.WithContext(ctx)
	err := sys.RunCommand("debootstrap", b.Suite, b.Target, b.Mirror)
	if err != nil {
		return err
	}

	return nil
}

// PostBootstrap runs post-bootstrap steps
func (b *BootstrapInstaller) Install(ctx context.Context, chroot system.System) error {
	// Configure locales
	if err := configureLocales(chroot, b.Target); err != nil {
		return err
	}

	// Install base packages
	if err := installBasePackages(chroot); err != nil {
		return err
	}

	// Configure repositories
	if err := b.configureRepositories(ctx); err != nil {
		return err
	}

	// Install Kubernetes packages
	if err := installKubePackages(chroot); err != nil {
		return err
	}

	// Configure system for kubernetes
	if err := common.ConfigureKubernetes(chroot, b.Target); err != nil {
		return err
	}

	// Configure networking
	if err := common.ConfigureNetwork(chroot, b.Target); err != nil {
		return err
	}

	// Configure the admin user
	if err := configureAdminUser(chroot); err != nil {
		return err
	}

	// Install kernel
	if err := installKernel(chroot); err != nil {
		return err
	}

	return nil
}
