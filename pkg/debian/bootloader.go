package debian

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
)

type GrubEFI struct {
	Target string
	Device string
}

func NewGrubEFI(target string) *GrubEFI {
	return &GrubEFI{
		Target: target,
	}
}

func (g *GrubEFI) Install(device string) error {
	// Install grub-efi
	log.Infof("Installing grub-efi")
	if err := installPackages(g.Target, "grub-efi"); err != nil {
		return err
	}

	// Run grub-install
	log.Infof("Running grub-install")
	if err := util.RunCmd("chroot", g.Target, "grub-install", "--removable", device); err != nil {
		return err
	}

	// Run grub-mkconfig
	log.Infof("Running grub-mkconfig")
	if err := util.RunCmd("chroot", g.Target, "grub-mkconfig", "-o", "/boot/grub/grub.cfg"); err != nil {
		return err
	}

	return nil
}
