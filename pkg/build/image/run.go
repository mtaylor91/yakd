package image

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildImage builds a yakd image from a stage1 tarball
func BuildImage(
	force bool, sizeMB int,
	stage1, target, mountpoint string,
	noCleanup bool,
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

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

	// Create image
	raw := util.NewRawImage(target+".raw", sizeMB, true, true)
	log.Infof("Creating raw image at %s", raw.ImagePath)
	if err := raw.Alloc(); err != nil {
		return err
	}

	// Create partition table
	log.Infof("Creating partition table on %s", raw.ImagePath)
	if err := raw.Partition(); err != nil {
		return err
	}

	// Attach image
	log.Infof("Attaching image %s", raw.ImagePath)
	loop, err := raw.Attach()
	if err != nil {
		return err
	}

	defer loop.Detach()

	// Format image
	log.Infof("Formatting image %s on %s", raw.ImagePath, loop.DevicePath)
	if err := loop.Format(); err != nil {
		return err
	}

	// Identify partitions
	esp := loop.DevicePath + "p1"
	root := loop.DevicePath + "p2"

	// Initialize image disk
	d := util.NewDisk(loop.DevicePath, esp, root, mountpoint, true)

	// Populate image disk
	if err := d.Populate(stage1); err != nil {
		return err
	}

	return fmt.Errorf("Not implemented")
}
