package stage1

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/build/release"
	"github.com/mtaylor91/yakd/pkg/build/release/debian"
	"github.com/mtaylor91/yakd/pkg/build/release/gentoo"
	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/tmpfs"
)

// BuildStage1 builds a stage1 tarball
func (stage1 *Stage1) Build(ctx context.Context) error {
	sys := system.Local.WithContext(ctx)

	// Construct target path
	target, err := util.TemplateString(stage1.TargetTemplate, map[string]string{
		"OS":   stage1.OS,
		"Arch": runtime.GOARCH,
	})
	if err != nil {
		return err
	}

	// Check if target exists
	if _, err := os.Stat(target); err == nil {
		if stage1.Force {
			// Remove target
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove target: %s", err)
			}
		} else {
			return fmt.Errorf("target already exists: %s", target)
		}
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	log.Infof("Building %s", target)

	tmpfs := &tmpfs.TmpFS{
		Path:   stage1.Mountpoint,
		SizeMB: stage1.TmpFSSize,
	}

	// Create mountpoint
	if err := tmpfs.Allocate(ctx); err != nil {
		return err
	}

	defer tmpfs.Destroy()

	// Select base OS
	var rel release.OS
	switch stage1.OS {
	case "debian":
		debian := debian.Default
		debian.Suite = stage1.DebianSuite
		debian.Mirror = stage1.DebianMirror
		rel = debian
	case "gentoo":
		gentoo := gentoo.Default
		gentoo.BinPkgsCache = stage1.GentooBinPkgsCache
		gentoo.Stage3 = stage1.GentooStage3
		rel = gentoo
	default:
		return fmt.Errorf("unknown operating system: %s", stage1.OS)
	}

	// Bootstrap OS
	installer := rel.BootstrapInstaller(tmpfs.Path)
	if err := installer.Bootstrap(ctx); err != nil {
		return err
	}

	// PostBootstrap via chroot
	if err := install(ctx, tmpfs.Path, installer, sys); err != nil {
		return err
	}

	// Create archive
	log.Infof("Creating stage1 archive at %s", target)
	err = sys.RunCommand("tar", "-C", stage1.Mountpoint, "-caf", target, ".")
	if err != nil {
		return err
	}

	return nil
}

func install(
	ctx context.Context,
	path string,
	installer release.BootstrapInstaller,
	localSystem system.System,
) error {
	log.Infof("Setting up chroot at %s", path)
	chrootSystem := system.Chroot(localSystem, path)

	if err := chrootSystem.Setup(); err != nil {
		return err
	}

	defer chrootSystem.Teardown()

	// Run post-bootstrap step
	log.Infof("Running install step")
	if err := installer.Install(ctx, chrootSystem); err != nil {
		return err
	}

	return nil
}
