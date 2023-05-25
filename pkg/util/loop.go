package util

import (
	log "github.com/sirupsen/logrus"
)

// LoopDevice represents a loop device
type LoopDevice struct {
	DevicePath string
}

// Detach detaches the loop device
func (l *LoopDevice) Detach() {
	// Detach loop device
	if err := RunCmd("losetup", "-d", l.DevicePath); err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}

// Format formats the image partitions via the loop device
func (l *LoopDevice) Format() error {
	// Create FAT32 filesystem on EFI partition
	if err := RunCmd("mkfs.vfat", "-F", "32", l.DevicePath+"p2"); err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	if err := RunCmd("mkfs.ext4", l.DevicePath+"p3"); err != nil {
		return err
	}

	return nil
}
