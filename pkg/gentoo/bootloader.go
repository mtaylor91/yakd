package gentoo

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type GentooBootloaderInstaller struct {
	device string
	target string
	exec   executor.Executor
}

// Install installs the bootloader.
func (g *GentooBootloaderInstaller) Install(ctx context.Context) error {
	return nil
}
