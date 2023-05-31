package image

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

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
