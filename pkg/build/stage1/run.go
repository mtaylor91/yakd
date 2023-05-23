package stage1

import (
	"github.com/mtaylor91/yakd/pkg/bootstrap"
	"github.com/mtaylor91/yakd/pkg/debian"
)

// BuildStage1 builds a stage1 tarball
func BuildStage1(
	force bool, target, suite, mirror, mountpoint string,
	tmpfsSize int, cleanup bool,
) error {
	debian := debian.DebianDefault
	debian.Suite = suite
	debian.Mirror = mirror

	stage1 := &Stage1{
		Source: mountpoint,
		Target: target,
	}

	tmpfs := &bootstrap.TmpFS{
		Path:   mountpoint,
		SizeMB: tmpfsSize,
	}

	err := tmpfs.Bootstrap(debian)
	if err != nil {
		return err
	}

	if cleanup {
		defer tmpfs.Destroy()
	}

	err = stage1.BuildArchive()
	if err != nil {
		return err
	}

	return nil
}
