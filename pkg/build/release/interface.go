package release

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/system"
)

type OS interface {
	BootstrapInstaller(target string) BootstrapInstaller
	BootloaderInstaller(device, target string, sys system.System) BootloaderInstaller
	HybridISOSourceBuilder(fsDir, isoDir string) HybridISOSourceBuilder
}

type BootstrapInstaller interface {
	Bootstrap(ctx context.Context) error
	Install(ctx context.Context, chroot system.System) error
}

type BootloaderInstaller interface {
	Install(ctx context.Context) error
}

type HybridISOSourceBuilder interface {
	BuildISOFS(ctx context.Context, chroot system.System) error
	BuildISOSources(ctx context.Context) error
}
