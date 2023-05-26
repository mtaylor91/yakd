package stage1

import (
	"context"
	"fmt"
	"os"

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
	// Check if target exists
	if _, err := os.Stat(target); err == nil {
		if force {
			// Remove target
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove target: %s", err)
			}
		} else {
			return fmt.Errorf("target already exists: %s", target)
		}
	}

	debian := debian.DebianDefault
	debian.Suite = suite
	debian.Mirror = mirror

	tmpfs := &bootstrap.TmpFS{
		Path:   mountpoint,
		SizeMB: tmpfsSize,
	}

	// Create mountpoint
	if err := tmpfs.Allocate(ctx); err != nil {
		return err
	}

	defer tmpfs.Destroy()

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
