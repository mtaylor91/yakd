package image

import (
	"fmt"
	"os"
)

// BuildImage builds a yakd image from a stage1 tarball
func BuildImage(force bool, stage1, target, mountpoint string, noCleanup bool) error {
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

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball not found: %s", stage1)
	}

	return fmt.Errorf("Not implemented")
}
