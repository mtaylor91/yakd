package util

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// LoopDevice represents a loop device
type LoopDevice struct {
	DevicePath string
}

// Detach detaches the loop device
func (l *LoopDevice) Detach() {
	// Get losetup path
	losetup, err := exec.LookPath("losetup")
	if err != nil {
		log.Errorf("Failed to get losetup path: %s", err)
	}

	// Detach loop device
	cmd := exec.Command(losetup, "-d", l.DevicePath)
	err = cmd.Run()
	if err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}

// Format formats the image partitions via the loop device
func (l *LoopDevice) Format() error {
	// Get mkfs.vfat path
	mkfsVfat, err := exec.LookPath("mkfs.vfat")
	if err != nil {
		return err
	}

	// Get mkfs.ext4 path
	mkfsExt4, err := exec.LookPath("mkfs.ext4")
	if err != nil {
		return err
	}

	// Create FAT32 filesystem on EFI partition
	cmd := exec.Command(mkfsVfat, l.DevicePath+"p1")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	cmd = exec.Command(mkfsExt4, l.DevicePath+"p2")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
