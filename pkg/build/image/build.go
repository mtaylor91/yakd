package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/gentoo"
	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildImage builds a yakd image from a stage1 tarball
func (c *Config) BuildImage(
	ctx context.Context,
) error {
	// Construct stage1 path
	stage1, err := util.TemplateString(c.Stage1Template, map[string]string{
		"OS":   c.OS,
		"Arch": runtime.GOARCH,
	})
	if err != nil {
		return err
	}

	// Construct target path
	target, err := util.TemplateString(c.TargetTemplate, map[string]string{
		"OS":   c.OS,
		"Arch": runtime.GOARCH,
	})
	if err != nil {
		return err
	}

	// Detect image type
	switch filepath.Ext(target) {
	case ".img":
		return c.buildIMG(ctx, stage1, target)
	case ".iso":
		return c.buildISO(ctx, stage1, target)
	case ".qcow2":
		return c.buildQcow2(ctx, stage1, target)
	default:
		return fmt.Errorf("unknown image type: %s", filepath.Ext(target))
	}
}

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

	return fmt.Errorf("ISO image creation not yet implemented")
}

// buildQcow2 builds a qcow2 image
func (c *Config) buildQcow2(
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

	// Build raw image
	rawImagePath := target + ".img"
	rawImage := util.NewRawImage(rawImagePath, c.SizeMB, true, true)
	if err := c.buildIMG(ctx, stage1, rawImagePath); err != nil {
		return err
	}

	// Clean up the raw image when we're done
	defer rawImage.Free()

	// Convert image to qcow2
	log.Infof("Converting image %s to %s", rawImagePath, target)
	if err := rawImage.Convert(ctx, target); err != nil {
		return err
	}

	return nil
}
