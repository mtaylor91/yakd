package gentoo

import (
	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type Gentoo struct {
	Stage3 string
}

type GentooBootloaderInstaller struct {
	device string
	target string
	exec   executor.Executor
}

func (g *Gentoo) BootstrapInstaller(
	target string,
) os.OSBootstrapInstaller {
	return &GentooBootstrapInstaller{g.Stage3, target}
}

func (g *Gentoo) BootloaderInstaller(
	device, target string, exec executor.Executor,
) os.OSBootloaderInstaller {
	return &GentooBootloaderInstaller{}
}
