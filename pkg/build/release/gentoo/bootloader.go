package gentoo

import (
	"context"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
)

const grubDefault = `GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="YAKD"
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX="console=tty1 console=ttyS0,115200n8"
GRUB_TERMINAL="console serial"
GRUB_SERIAL_COMMAND="serial --speed=115200 --unit=0 --word=8 --parity=no --stop=1"
`

type GentooBootloaderInstaller struct {
	binPkgsCache string
	device       string
	target       string
	system       system.System
}

// Install installs the bootloader.
func (g *GentooBootloaderInstaller) Install(ctx context.Context) error {
	sys := system.Local.WithContext(ctx)

	// Ensure binPkgsCache exists
	err := os.MkdirAll(g.binPkgsCache, 0755)
	if err != nil {
		return err
	}

	// Bind binPkgsCache to /var/cache/binpkgs
	if err = sys.RunCommand(
		"mount", "--bind",
		g.binPkgsCache,
		path.Join(g.target, "var/cache/binpkgs"),
	); err != nil {
		return err
	}

	// Unmount /var/cache/binpkgs on exit
	defer func() {
		if err := sys.RunCommand(
			"umount", path.Join(g.target, "var/cache/binpkgs"),
		); err != nil {
			log.Warnf("Failed to unmount /var/cache/binpkgs: %s", err)
		}
	}()

	err = installPackages(g.system, "sys-boot/grub")
	if err != nil {
		return err
	}

	err = util.WriteFile(path.Join(g.target, "etc/default/grub"), grubDefault)
	if err != nil {
		return err
	}

	err = g.system.RunCommand(
		"grub-install", "--removable", "--efi-directory", "/boot/efi", g.device)
	if err != nil {
		return err
	}

	err = g.system.RunCommand("grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	if err != nil {
		return err
	}

	return nil
}
