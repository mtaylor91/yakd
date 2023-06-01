package gentoo

import (
	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type Gentoo struct {
	BinPkgsCache string
	Stage3       string
}

func (g *Gentoo) BootstrapInstaller(
	target string,
) os.OSBootstrapInstaller {
	return &GentooBootstrapInstaller{g.BinPkgsCache, g.Stage3, target}
}

func (g *Gentoo) DiskInstaller(
	device, target string, exec executor.Executor,
) os.OSBootloaderInstaller {
	return &GentooBootloaderInstaller{g.BinPkgsCache, device, target, exec}
}

func (g *Gentoo) HybridISOSourceBuilder(fsDir, isoDir string) os.HybridISOSourceBuilder {
	return &HybridISOSourceBuilder{g.BinPkgsCache, fsDir, isoDir}
}

func (g *Gentoo) HybridISOBuilder(isoDir, target string) os.HybridISOBuilder {
	return &HybridISOBuilder{isoDir, target}
}
