package os

import "github.com/mtaylor91/yakd/pkg/util/executor"

type OS interface {
	BootstrapInstaller(target string) OSBootstrapInstaller
	BootloaderInstaller(
		device, target string, exec executor.Executor,
	) OSBootloaderInstaller
}
