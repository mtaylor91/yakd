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

func (g *Gentoo) BootloaderInstaller(
	device, target string, exec executor.Executor,
) os.OSBootloaderInstaller {
	return &GentooBootloaderInstaller{device, target, exec}
}
