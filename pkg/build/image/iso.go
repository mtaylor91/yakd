package image

import (
	"context"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/gentoo"
	yakdOS "github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/tmpfs"
)

// buildISO builds an ISO image
func (c *Config) buildISO(
	ctx context.Context,
	stage1 string,
	target string,
) error {
	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

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

	// Allocate tmpfs for ISO filesystem
	tmpfs := &tmpfs.TmpFS{Path: c.Mountpoint, SizeMB: c.SizeMB}
	if err := tmpfs.Allocate(ctx); err != nil {
		return err
	}

	defer tmpfs.Destroy()

	// Construct tmpfs subpaths
	fsDir := path.Join(tmpfs.Path, "fs")
	isoDir := path.Join(tmpfs.Path, "iso")

	// Create tmpfs subpaths
	if err := os.MkdirAll(fsDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(isoDir, 0755); err != nil {
		return err
	}

	// Populate fsDir from stage1
	log.Infof("Unpacking %s to %s", stage1, fsDir)
	if err := util.UnpackTarball(ctx, stage1, fsDir); err != nil {
		return err
	}

	// Select base OS
	var yakdOS yakdOS.OS
	switch c.OS {
	case "debian":
		debian := debian.DebianDefault
		yakdOS = debian
	case "gentoo":
		gentoo := gentoo.DefaultGentoo
		gentoo.BinPkgsCache = c.GentooBinPkgsCache
		yakdOS = gentoo
	default:
		return fmt.Errorf("unknown operating system: %s", c.OS)
	}

	sourceBuilder := yakdOS.HybridISOSourceBuilder(fsDir, isoDir)
	err := c.buildISOChroot(ctx, fsDir, isoDir, sourceBuilder)
	if err != nil {
		return err
	}

	isoBuilder := yakdOS.HybridISOBuilder(isoDir, target)
	err = c.buildISOHybrid(ctx, isoDir, target, sourceBuilder, isoBuilder)
	return err
}

// buildISOChroot builds the ISO filesystem and bootloader images in a chroot
func (c *Config) buildISOChroot(
	ctx context.Context,
	fsDir, isoDir string,
	sourceBuilder yakdOS.HybridISOSourceBuilder,
) error {
	// Setup chroot
	log.Infof("Setting up chroot")
	localSystem := system.Local.WithContext(ctx)
	chrootSystem := system.Chroot(localSystem, fsDir)
	if err := chrootSystem.Setup(); err != nil {
		return err
	}

	defer chrootSystem.Teardown()

	// Build ISO filesystem
	log.Infof("Building source(s) for %s hybrid ISO", c.OS)
	return sourceBuilder.BuildISOFS(ctx, chrootSystem)
}

// buildISOHybrid builds the ISO image from the ISO sources
func (c *Config) buildISOHybrid(
	ctx context.Context,
	isoDir, target string,
	sourceBuilder yakdOS.HybridISOSourceBuilder,
	isoBuilder yakdOS.HybridISOBuilder,
) error {
	// Build ISO sources
	log.Infof("Building source(s) for %s hybrid ISO", c.OS)
	if err := sourceBuilder.BuildISOSources(ctx); err != nil {
		return err
	}

	// Build ISO
	log.Infof("Building %s hybrid ISO", c.OS)
	if err := isoBuilder.BuildISO(ctx); err != nil {
		return err
	}

	return nil
}
