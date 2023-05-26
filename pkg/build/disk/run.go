package disk

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/debian"
	"github.com/mtaylor91/yakd/pkg/util"
)

// BuildDisk builds a disk from a stage1 tarball
func BuildDisk(ctx context.Context, target, stage1, mountpoint string) error {
	log.Infof("Building disk %s from %s", target, stage1)

	debian := debian.DebianDefault

	// Check if target exists
	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("target %s: %s", target, err)
	}

	// Check if stage1 exists
	if _, err := os.Stat(stage1); err != nil {
		return fmt.Errorf("stage1 tarball %s: %s", stage1, err)
	}

	// Partition disk
	log.Infof("Partitioning %s", target)
	if err := util.PartitionDisk(ctx, target); err != nil {
		return fmt.Errorf("partitioning %s: %s", target, err)
	}

	// Initialize disk
	d, err := util.NewDisk(target, mountpoint, true)
	if err != nil {
		return fmt.Errorf("initializing disk: %s", err)
	}

	// Format disk partitions
	log.Infof("Formatting %s", target)
	if err := d.Format(ctx); err != nil {
		return err
	}

	// Populate disk
	log.Infof("Populating %s", target)
	if err := d.Populate(ctx, stage1, debian); err != nil {
		return fmt.Errorf("populating %s: %s", target, err)
	}

	return nil
}
