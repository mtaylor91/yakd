package gentoo

import (
	"github.com/mtaylor91/yakd/pkg/build/release"
	"github.com/mtaylor91/yakd/pkg/system"
)

type Gentoo struct {
	BinPkgsCache string
	Stage3       string
}

func (g *Gentoo) BootstrapInstaller(
	target string,
) release.BootstrapInstaller {
	return &GentooBootstrapInstaller{g.BinPkgsCache, g.Stage3, target}
}

func (g *Gentoo) BootloaderInstaller(
	device, target string, sys system.System,
) release.BootloaderInstaller {
	return &GentooBootloaderInstaller{g.BinPkgsCache, device, target, sys}
}

func (g *Gentoo) HybridISOSourceBuilder(
	fsDir, isoDir string) release.HybridISOSourceBuilder {
	return &HybridISOSourceBuilder{g.BinPkgsCache, fsDir, isoDir}
}
