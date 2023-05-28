package gentoo

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type GentooBootloaderInstaller struct {
	device string
	target string
	exec   executor.Executor
}

// Install installs the bootloader.
func (g *GentooBootloaderInstaller) Install(ctx context.Context) error {
	err := installPackages(ctx, g.exec, "sys-boot/grub")
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
