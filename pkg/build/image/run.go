package image

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildImage builds a yakd image from a stage1 tarball
func BuildImage(
	ctx context.Context,
	force, raw bool, sizeMB int,
	stage1, target, mountpoint string,
) error {
	debian := debian.DebianDefault
	rawName := target

	if !raw {
		rawName = target + ".raw"
	}

	r := util.NewRawImage(rawName, sizeMB, true, true)

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

	// Check if raw image exists
	if _, err := os.Stat(r.ImagePath); err == nil {
		if force {
			// Remove raw image
			if err := os.Remove(r.ImagePath); err != nil {
				return fmt.Errorf("failed to remove raw image: %s", err)
			}
		} else {
			return fmt.Errorf("raw image already exists: %s", r.ImagePath)
		}
	}

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

	// Allocate raw image file
	log.Infof("Creating raw image at %s", r.ImagePath)
	if err := r.Alloc(ctx); err != nil {
		return err
	}

	if !raw {
		defer r.Free()
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

	// Format image
	log.Infof("Formatting image %s on %s", r.ImagePath, loop.DevicePath)
	if err := loop.Format(ctx); err != nil {
		return err
	}

	// Initialize image loop device disk
	d := util.NewDisk(loop.DevicePath, mountpoint, true)

	// Populate disk image
	log.Infof("Populating disk image mounted at %s", mountpoint)
	if err := d.Populate(ctx, stage1, debian); err != nil {
		return err
	}

	if !raw {
		// Convert image to qcow2
		log.Infof("Converting image %s to %s", r.ImagePath, target)
		if err := r.Convert(ctx, target); err != nil {
			return err
		}
	}

	return nil
}
