package util

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// UnpackTarball unpacks a tarball to the specified target
func UnpackTarball(ctx context.Context, source, target string) error {
	// Unpack via tar
	return executor.RunCmd(ctx, "tar", "-xf", source, "-C", target)
}
