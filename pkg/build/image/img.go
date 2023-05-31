package image

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/gentoo"
	"github.com/mtaylor91/yakd/pkg/util"
)

// buildIMG builds a raw image
func (c *Config) buildIMG(
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

	// Allocate raw image file
	r := util.NewRawImage(target, c.SizeMB, true, true)
	log.Infof("Creating raw image at %s", r.ImagePath)
	if err := r.Alloc(ctx); err != nil {
		return err
	}

	// Create partition table
	log.Infof("Creating partition table on %s", r.ImagePath)
	if err := util.PartitionDisk(ctx, r.ImagePath); err != nil {
		return err
	}

	// Attach image
	log.Infof("Attaching image %s", r.ImagePath)
	loop, err := r.Attach(ctx)
	if err != nil {
		return err
	}

	defer loop.Detach()

	// Initialize image loop device disk
	d, err := util.NewDisk(loop.DevicePath, c.Mountpoint, true)
	if err != nil {
		return err
	}

	// Format disk image partitions
	log.Infof("Formatting image %s on %s", r.ImagePath, loop.DevicePath)
	if err := d.Format(ctx); err != nil {
		return err
	}

	// Populate disk image
	log.Infof("Populating disk image mounted at %s", c.Mountpoint)
	switch c.OS {
	case "debian":
		debian := debian.DebianDefault
		err = d.Populate(ctx, stage1, debian)
	case "gentoo":
		gentoo := gentoo.DefaultGentoo
		err = d.Populate(ctx, stage1, gentoo)
	default:
		err = fmt.Errorf("unsupported OS: %s", c.OS)
	}
	if err != nil {
		return err
	}

	return nil
}
