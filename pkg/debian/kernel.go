package debian

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// installKernel installs the kernel specified in the kernel string
func (c *BootstrapConfig) installKernel(exec executor.Executor) error {
	// Install kernel
	log.Infof("Installing kernel")
	if err := installPackages(exec, "linux-image-amd64"); err != nil {
		return err
	}

	return nil
}
