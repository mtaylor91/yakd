package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type GrubInstaller struct {
	Device   string
	Target   string
	Executor executor.Executor
}

func NewGrubInstaller(device, target string, exec executor.Executor) *GrubInstaller {
	return &GrubInstaller{
		Device:   device,
		Target:   target,
		Executor: exec,
	}
}

func (g *GrubInstaller) Install(ctx context.Context) error {
	// Install grub-efi
	log.Infof("Installing grub")
	if err := installPackages(ctx, g.Executor, "grub-efi"); err != nil {
		return err
	}

	// Run grub-install
	log.Infof("Running grub-install")
	err := g.Executor.RunCmd(
		ctx, "grub-install", "--force-extra-removable", g.Device)
	if err != nil {
		return err
	}

	// Run grub-mkconfig
	log.Infof("Running grub-mkconfig")
	err = g.Executor.RunCmd(ctx, "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	if err != nil {
		return err
	}

	return nil
}
