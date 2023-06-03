package os

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/system"
)

type OS interface {
	// BootstrapInstaller bootstraps the stage1 filesystem
	BootstrapInstaller(target string) OSBootstrapInstaller

	// DiskInstaller installs the bootloader to disk
	DiskInstaller(
		device, target string, sys system.System,
	) OSBootloaderInstaller

	// HybridISOSourceBuilder builds hybrid ISO source files
	HybridISOSourceBuilder(fsDir, isoDir string) HybridISOSourceBuilder

	// HybridISOBuilder builds hybrid ISO from source files
	HybridISOBuilder(isoDir, target string) HybridISOBuilder
}

type OSBootstrapInstaller interface {
	Bootstrap(ctx context.Context) error
	PostBootstrap(ctx context.Context, chroot system.System) error
}

type OSBootloaderInstaller interface {
	Install(ctx context.Context) error
}

type HybridISOSourceBuilder interface {
	BuildISOFS(ctx context.Context, chroot system.System) error
	BuildISOSources(ctx context.Context) error
}

type HybridISOBuilder interface {
	BuildISO(ctx context.Context) error
}
