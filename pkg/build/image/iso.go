package image

import (
	"context"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/build/release"
	"github.com/mtaylor91/yakd/pkg/build/release/debian"
	"github.com/mtaylor91/yakd/pkg/build/release/gentoo"
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
	var release release.OS
	switch c.OS {
	case "debian":
		debian := debian.Default
		release = debian
	case "gentoo":
		gentoo := gentoo.Default
		gentoo.BinPkgsCache = c.GentooBinPkgsCache
		release = gentoo
	default:
		return fmt.Errorf("unknown operating system: %s", c.OS)
	}

	sourceBuilder := release.HybridISOSourceBuilder(fsDir, isoDir)
	err := c.buildISOChroot(ctx, fsDir, isoDir, sourceBuilder)
	if err != nil {
		return err
	}

	err = c.buildISOHybrid(ctx, isoDir, target, sourceBuilder)
	return err
}

// buildISOChroot builds the ISO filesystem and bootloader images in a chroot
func (c *Config) buildISOChroot(
	ctx context.Context,
	fsDir, isoDir string,
	sourceBuilder release.HybridISOSourceBuilder,
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
	sourceBuilder release.HybridISOSourceBuilder,
) error {
	// Build ISO sources
	log.Infof("Building source(s) for %s hybrid ISO", c.OS)
	if err := sourceBuilder.BuildISOSources(ctx); err != nil {
		return err
	}

	// Build ISO
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand(
		"xorrisofs",
		"-iso-level", "3",
		"-full-iso9660-filenames",
		"-volid", "YAKD",
		"-eltorito-boot", "bios.img",
		"-no-emul-boot", "-boot-load-size", "4", "-boot-info-table",
		"-isohybrid-mbr",
		path.Join(isoDir, "isohdpfx.bin"),
		"--efi-boot", "efi.img",
		"-efi-boot-part",
		"--efi-boot-image",
		"--protective-msdos-label",
		"-output", target,
		isoDir,
	); err != nil {
		return err
	}

	return nil
}
