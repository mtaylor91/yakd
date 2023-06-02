package debian

import (
	"context"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const grubDefault = `GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="YAKD"
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX="console=tty1 console=ttyS0,115200n8"
GRUB_TERMINAL="console serial"
GRUB_SERIAL_COMMAND="serial --speed=115200 --unit=0 --word=8 --parity=no --stop=1"
`

type GrubDiskInstaller struct {
	Device   string
	Target   string
	Executor executor.Executor
}

func NewGrubDiskInstaller(
	device, target string, exec executor.Executor,
) *GrubDiskInstaller {
	return &GrubDiskInstaller{
		Device:   device,
		Target:   target,
		Executor: exec,
	}
}

func (g *GrubDiskInstaller) Install(ctx context.Context) error {
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

	// Write grub default
	log.Infof("Writing grub default")
	err = util.WriteFile(
		path.Join(g.Target, "etc/default/grub"), grubDefault)
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
