package os

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type OS interface {
	BootstrapInstaller(target string) OSBootstrapInstaller
	DiskInstaller(
		device, target string, exec executor.Executor,
	) OSBootloaderInstaller
}

type OSBootstrapInstaller interface {
	Bootstrap(ctx context.Context) error
	PostBootstrap(ctx context.Context, chroot executor.Executor) error
}

type OSBootloaderInstaller interface {
	Install(ctx context.Context) error
}
