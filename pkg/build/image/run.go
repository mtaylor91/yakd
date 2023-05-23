package image

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildImage builds a yakd image from a stage1 tarball
func BuildImage(
	force bool, sizeMB int,
	stage1, target, mountpoint string,
) error {
	debian := debian.DebianDefault
	raw := util.NewRawImage(target+".raw", sizeMB, true, true)

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
	if _, err := os.Stat(raw.ImagePath); err == nil {
		if force {
			// Remove raw image
			if err := os.Remove(raw.ImagePath); err != nil {
				return fmt.Errorf("failed to remove raw image: %s", err)
			}
		} else {
			return fmt.Errorf("raw image already exists: %s", raw.ImagePath)
		}
	}

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

	// Allocate raw image file
	log.Infof("Creating raw image at %s", raw.ImagePath)
	if err := raw.Alloc(); err != nil {
		return err
	}

	defer raw.Free()

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

	// Initialize image loop device disk
	d := util.NewDisk(loop.DevicePath, esp, root, mountpoint, true)

	// Populate disk image
	log.Infof("Populating disk image mounted at %s", mountpoint)
	if err := d.Populate(stage1, debian); err != nil {
		return err
	}

	// Convert image to qcow2
	log.Infof("Converting image %s to %s", raw.ImagePath, target)
	if err := raw.Convert(target); err != nil {
		return err
	}

	return nil
}
