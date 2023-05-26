package os

import "context"

type OSBootloaderInstaller interface {
	Install(ctx context.Context) error
}
