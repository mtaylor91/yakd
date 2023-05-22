package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const localeGen = `
en_CA.UTF-8 UTF-8
en_US.UTF-8 UTF-8
`

// configureLocales configures the locales
func configureLocales(target string) error {
	// Look for chroot
	chroot, err := exec.LookPath("chroot")
	if err != nil {
		return err
	}

	// Install locales
	cmd := exec.Command(chroot, target, "apt-get", "install", "-y", "locales")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Write locale.gen
	log.Infof("Writing locale.gen")
	localeGenPath := target + "/etc/locale.gen"
	if err := os.WriteFile(localeGenPath, []byte(localeGen), 0644); err != nil {
		return err
	}

	// Configure locales
	log.Infof("Configuring locales")
	cmd = exec.Command(chroot, target, "locale-gen")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
