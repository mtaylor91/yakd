package image

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
)

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
