package os

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type OSBootstrapInstaller interface {
	Bootstrap(ctx context.Context) error
	PostBootstrap(ctx context.Context, chroot executor.Executor) error
}
