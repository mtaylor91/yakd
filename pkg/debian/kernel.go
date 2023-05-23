package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// installKernel installs the kernel specified in the kernel string
func (c *BootstrapConfig) installKernel() error {
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
