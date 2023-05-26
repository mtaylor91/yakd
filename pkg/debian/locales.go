package debian

import (
	"context"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const localeGen = `
en_CA.UTF-8 UTF-8
en_US.UTF-8 UTF-8
`

// configureLocales configures the locales
func configureLocales(ctx context.Context, exec executor.Executor, root string) error {
	// Install locales
	if err := installPackages(ctx, exec, "locales"); err != nil {
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
	if err := exec.RunCmd(ctx, "locale-gen"); err != nil {
		return err
	}

	return nil
}
