package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// InstallKernel installs the kernel specified in the kernel string
func (c *BootstrapConfig) InstallKernel() error {
	// Look for chroot
	chroot, err := exec.LookPath("chroot")
	if err != nil {
		return err
	}

	// Install kernel
	log.Infof("Installing kernel")
	cmd := exec.Command(chroot, c.Target, "apt-get", "install", "-y", "linux-image-amd64")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// InstallBootloader installs the bootloader
func (c *BootstrapConfig) InstallBootloader() error {
	// Look for chroot
	chroot, err := exec.LookPath("chroot")
	if err != nil {
		return err
	}

	// Install grub-efi
	log.Infof("Installing grub-efi")
	cmd := exec.Command(chroot, c.Target, "apt-get", "install", "-y", "grub-efi")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
