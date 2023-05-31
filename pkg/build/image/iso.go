package image

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/chroot"
	"github.com/mtaylor91/yakd/pkg/util/tmpfs"
)

// buildISO builds an ISO image
func (c *Config) buildISO(
	ctx context.Context,
	stage1 string,
	target string,
) error {
	// Check if target exists
	if _, err := os.Stat(target); err == nil {
		if c.Force {
			// Remove target
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove %s: %s",
					target, err)
			}
		} else {
			return fmt.Errorf("%s already exists", target)
		}
	}

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

	tmpfs := &tmpfs.TmpFS{
		Path:   c.Mountpoint,
		SizeMB: c.SizeMB,
	}

	// Allocate tmpfs for ISO filesystem
	if err := tmpfs.Allocate(ctx); err != nil {
		return err
	}

	defer tmpfs.Destroy()

	// Populate tmpfs
	log.Infof("Copying source %s to %s", stage1, tmpfs.Path)
	if err := util.UnpackTarball(ctx, stage1, tmpfs.Path); err != nil {
		return err
	}

	if err := c.buildISOChroot(ctx, tmpfs); err != nil {
		return err
	}

	return fmt.Errorf("ISO image creation not yet implemented")
}

func (c *Config) buildISOChroot(
	ctx context.Context,
	tmpfs *tmpfs.TmpFS,
) error {
	log.Infof("Setting up chroot")
	chrootExecutor := chroot.NewExecutor(ctx, tmpfs.Path)
	defer chrootExecutor.Teardown()
	defer chrootExecutor.RunCmdWithStdin(ctx, "/bin/bash", os.Stdin)
	return nil
}
