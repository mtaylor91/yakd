package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// installKernel installs the kernel specified in the kernel string
func (c *BootstrapConfig) installKernel(ctx context.Context, exec executor.Executor) error {
	// Install kernel
	log.Infof("Installing kernel")
	if err := installPackages(ctx, exec, "linux-image-amd64"); err != nil {
		return err
	}

	return nil
}
