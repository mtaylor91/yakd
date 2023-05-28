package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

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

	// Infer if we're building a raw image
	raw := false
	if filepath.Ext(target) == ".raw" {
		raw = true
	}

	rawName := target

	if !raw {
		rawName = target + ".raw"
	}

	r := util.NewRawImage(rawName, c.SizeMB, true, true)

	// Check if target exists
	if _, err := os.Stat(target); err == nil {
		if c.Force {
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
		if c.Force {
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

	// Sleep to allow kernel to update partition table
	log.Infof("Sleeping for 5 seconds to allow kernel to update partition table")
	time.Sleep(5 * time.Second)

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

	if !raw {
		// Convert image to qcow2
		log.Infof("Converting image %s to %s", r.ImagePath, target)
		if err := r.Convert(ctx, target); err != nil {
			return err
		}
	}

	return nil
}
