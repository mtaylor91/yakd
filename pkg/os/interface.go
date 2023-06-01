package os

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type OS interface {
	// BootstrapInstaller bootstraps the stage1 filesystem
	BootstrapInstaller(target string) OSBootstrapInstaller

	// DiskInstaller installs the bootloader to disk
	DiskInstaller(
		device, target string, exec executor.Executor,
	) OSBootloaderInstaller

	// HybridISOSourceBuilder builds hybrid ISO source files
	HybridISOSourceBuilder(fsDir, isoDir string) HybridISOSourceBuilder

	// HybridISOBuilder builds hybrid ISO from source files
	HybridISOBuilder(isoDir, target string) HybridISOBuilder
}

type OSBootstrapInstaller interface {
	Bootstrap(ctx context.Context) error
	PostBootstrap(ctx context.Context, chroot executor.Executor) error
}

type OSBootloaderInstaller interface {
	Install(ctx context.Context) error
}

type HybridISOSourceBuilder interface {
	BuildISOFS(ctx context.Context, chroot executor.Executor) error
	BuildISOSources(ctx context.Context) error
}

type HybridISOBuilder interface {
	BuildISO(ctx context.Context) error
}
