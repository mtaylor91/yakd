package os

import "github.com/mtaylor91/yakd/pkg/util/executor"

type OS interface {
	Installer(target string) OSInstaller
	BootloaderInstaller(
		device, target string, exec executor.Executor,
	) OSBootloaderInstaller
}
