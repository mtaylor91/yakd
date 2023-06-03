package disk

import (
	"context"
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/build/release/debian"
	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildDisk builds a disk from a stage1 tarball
func (c *Config) BuildDisk(ctx context.Context) error {
	// Construct stage1 path
	stage1, err := util.TemplateString(c.Stage1Template, map[string]string{
		"OS":   c.OS,
		"Arch": runtime.GOARCH,
	})
	if err != nil {
		return err
	}

	log.Infof("Building disk %s from %s", c.Target, stage1)

	debian := debian.Default

	// Check if target exists
	if _, err := os.Stat(c.Target); err != nil {
		return fmt.Errorf("target %s: %s", c.Target, err)
	}

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball %s: %s", stage1, err)
	}

	// Partition disk
	log.Infof("Partitioning %s", c.Target)
	if err := util.PartitionDisk(ctx, c.Target); err != nil {
		return fmt.Errorf("partitioning %s: %s", c.Target, err)
	}

	// Initialize disk
	d, err := NewDisk(c.Target, c.Mountpoint, true)
	if err != nil {
		return fmt.Errorf("initializing disk: %s", err)
	}

	// Format disk partitions
	log.Infof("Formatting %s", c.Target)
	if err := d.Format(ctx); err != nil {
		return err
	}

	// Populate disk
	log.Infof("Populating %s", c.Target)
	if err := d.Populate(ctx, stage1, debian); err != nil {
		return fmt.Errorf("populating %s: %s", c.Target, err)
	}

	return nil
}
