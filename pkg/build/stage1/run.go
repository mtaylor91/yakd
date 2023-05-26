package stage1

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/util/bootstrap"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// BuildStage1 builds a stage1 tarball
func BuildStage1(
	ctx context.Context,
	force bool,
	target, suite, mirror, mountpoint string,
	tmpfsSize int,
) error {
	debian := debian.DebianDefault
	debian.Suite = suite
	debian.Mirror = mirror

	tmpfs := &bootstrap.TmpFS{
		Path:   mountpoint,
		SizeMB: tmpfsSize,
	}

	err := tmpfs.Bootstrap(ctx, debian)
	if err != nil {
		return err
	}

	// Create archive
	log.Infof("Creating stage1 archive at %s", target)
	err = executor.RunCmd(ctx, "tar", "-C", mountpoint, "-caf", target, ".")
	if err != nil {
		return err
	}

	return nil
}
