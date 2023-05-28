package gentoo

import (
	"context"
	"os"
	"path"

	"github.com/mtaylor91/yakd/pkg/util/executor"
	log "github.com/sirupsen/logrus"
)

type GentooBootloaderInstaller struct {
	binPkgsCache string
	device       string
	target       string
	exec         executor.Executor
}

// Install installs the bootloader.
func (g *GentooBootloaderInstaller) Install(ctx context.Context) error {
	// Ensure binPkgsCache exists
	err := os.MkdirAll(g.binPkgsCache, 0755)
	if err != nil {
		return err
	}

	// Bind binPkgsCache to /var/cache/binpkgs
	if err = executor.Default.RunCmd(
		ctx, "mount", "--bind",
		g.binPkgsCache,
		path.Join(g.target, "var/cache/binpkgs"),
	); err != nil {
		return err
	}

	// Unmount /var/cache/binpkgs on exit
	defer func() {
		if err := executor.Default.RunCmd(
			ctx, "umount", path.Join(g.target, "var/cache/binpkgs"),
		); err != nil {
			log.Warnf("Failed to unmount /var/cache/binpkgs: %s", err)
		}
	}()

	err = installPackages(ctx, g.exec, "sys-boot/grub")
	if err != nil {
		return err
	}

	err = g.exec.RunCmd(ctx,
		"grub-install", "--removable",
		"--efi-directory", "/boot/efi",
		g.device)
	if err != nil {
		return err
	}

	err = g.exec.RunCmd(ctx, "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	if err != nil {
		return err
	}

	return nil
}
