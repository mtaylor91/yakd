package stage1

import (
	"context"
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/gentoo"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/bootstrap"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// BuildStage1 builds a stage1 tarball
func (stage1 *Stage1) Build(ctx context.Context) error {
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

	log.Infof("Building %s", target)

	tmpfs := &bootstrap.TmpFS{
		Path:   stage1.Mountpoint,
		SizeMB: stage1.TmpFSSize,
	}

	// Create mountpoint
	if err := tmpfs.Allocate(ctx); err != nil {
		return err
	}

	defer tmpfs.Destroy()

	// Bootstrap operating system
	switch stage1.OS {
	case "debian":
		debian := debian.DebianDefault
		debian.Suite = stage1.DebianSuite
		debian.Mirror = stage1.DebianMirror
		err = tmpfs.Bootstrap(ctx, debian)
	case "gentoo":
		gentoo := gentoo.DefaultGentoo
		err = tmpfs.Bootstrap(ctx, gentoo)
	default:
		return fmt.Errorf("unknown operating system: %s", stage1.OS)
	}
	if err != nil {
		return err
	}

	// Create archive
	log.Infof("Creating stage1 archive at %s", target)
	err = executor.RunCmd(ctx, "tar", "-C", stage1.Mountpoint, "-caf", target, ".")
	if err != nil {
		return err
	}

	return nil
}
