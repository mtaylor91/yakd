package debian

import (
	"github.com/mtaylor91/yakd/pkg/system"
)

// installKernel installs the kernel specified in the kernel string
func installKernel(sys system.System) error {
	// Install kernel
	sys.Logger().Infof("Installing kernel")
	if err := installPackages(sys, "linux-image-amd64"); err != nil {
		return err
	}

	return nil
}
