package util

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/system"
)

// UnpackTarball unpacks a tarball to the specified target
func UnpackTarball(ctx context.Context, source, target string) error {
	// Unpack via tar
	sys := system.Local.WithContext(ctx)
	return sys.RunCommand("tar", "-xpf", source, "-C", target,
		"--xattrs-include='*.*'", "--numeric-owner")
}
