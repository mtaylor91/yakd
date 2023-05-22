package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

var basePackages = []string{
	"apt-transport-https",
	"ca-certificates",
	"curl",
	"gnupg2",
	"locales",
	"lvm2",
	"sudo",
}

// InstallBasePackages installs the base packages
func (c *BootstrapConfig) InstallBasePackages() error {
	// Look for chroot
	chroot, err := exec.LookPath("chroot")
	if err != nil {
		return err
	}

	// Install packages
	log.Infof("Installing base packages")
	args := []string{c.Target, "apt-get", "install", "-y"}
	args = append(args, basePackages...)
	cmd := exec.Command(chroot, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
