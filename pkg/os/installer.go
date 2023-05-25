package os

import "github.com/mtaylor91/yakd/pkg/util/executor"

type OSInstaller interface {
	Bootstrap() error
	PostBootstrap(chroot executor.Executor) error
}
