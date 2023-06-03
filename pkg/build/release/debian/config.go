package debian

import (
	"github.com/mtaylor91/yakd/pkg/build/release"
	"github.com/mtaylor91/yakd/pkg/system"
)

type Debian struct {
	Suite  string
	Mirror string
}

func (d *Debian) BootloaderInstaller(
	device, target string, sys system.System,
) release.BootloaderInstaller {
	return &BootloaderInstaller{device, target, sys}
}

func (d *Debian) BootstrapInstaller(target string) release.BootstrapInstaller {
	return &BootstrapInstaller{Suite: d.Suite, Mirror: d.Mirror, Target: target}
}

func (d *Debian) HybridISOSourceBuilder(
	fsDir, isoDir string) release.HybridISOSourceBuilder {
	return &HybridISOSourceBuilder{FSDir: fsDir, ISODir: isoDir}
}
