package debian

import (
	"os"
	"path"

	"github.com/mtaylor91/yakd/pkg/system"
	log "github.com/sirupsen/logrus"
)

const localeGen = `
en_CA.UTF-8 UTF-8
en_US.UTF-8 UTF-8
`

// configureLocales configures the locales
func configureLocales(sys system.System, root string) error {
	// Install locales
	if err := installPackages(sys, "locales"); err != nil {
		return err
	}

	// Write locale.gen
	log.Infof("Writing locale.gen")
	localeGenPath := path.Join(root, "etc", "locale.gen")
	if err := os.WriteFile(localeGenPath, []byte(localeGen), 0644); err != nil {
		return err
	}

	// Configure locales
	log.Infof("Configuring locales")
	if err := sys.RunCommand("locale-gen"); err != nil {
		return err
	}

	return nil
}
