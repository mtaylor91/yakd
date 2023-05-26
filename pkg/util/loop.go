package util

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// LoopDevice represents a loop device
type LoopDevice struct {
	DevicePath string
}

// Detach detaches the loop device
func (l *LoopDevice) Detach() {
	// Detach loop device
	ctx := context.Background()
	if err := executor.RunCmd(ctx, "losetup", "-d", l.DevicePath); err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}

// Format formats the image partitions via the loop device
func (l *LoopDevice) Format(ctx context.Context) error {
	// Create FAT32 filesystem on EFI partition
	err := executor.RunCmd(ctx, "mkfs.vfat", "-F", "32", l.DevicePath+"p2")
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	if err := executor.RunCmd(ctx, "mkfs.ext4", l.DevicePath+"p3"); err != nil {
		return err
	}

	return nil
}
