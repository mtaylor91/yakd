package util

import (
	"context"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	yakdOS "github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util/chroot"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// Disk represents the bootstrap configuration for a disk
type Disk struct {
	DevicePath    string
	biosPartition string
	espPartition  string
	rootPartition string
	mountpoint    string
	cleanup       bool
}

// NewDisk initializes a new Disk struct
func NewDisk(devicePath, mountpoint string, cleanup bool) (*Disk, error) {
	p1, err := identifyPartition(devicePath, 1)
	if err != nil {
		return nil, err
	}

	p2, err := identifyPartition(devicePath, 2)
	if err != nil {
		return nil, err
	}

	p3, err := identifyPartition(devicePath, 3)
	if err != nil {
		return nil, err
	}

	return &Disk{
		DevicePath:    devicePath,
		biosPartition: p1,
		espPartition:  p2,
		rootPartition: p3,
		mountpoint:    mountpoint,
		cleanup:       cleanup,
	}, nil
}

// Format the disk partitions
func (d *Disk) Format(ctx context.Context) error {
	// Create FAT32 filesystem on EFI partition
	err := executor.RunCmd(ctx, "mkfs.vfat", "-F", "32", d.espPartition)
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	if err := executor.RunCmd(ctx, "mkfs.ext4", d.rootPartition); err != nil {
		return err
	}

	return nil
}

// Populate the disk from the specified source
func (d *Disk) Populate(ctx context.Context, source string, yakdOS yakdOS.OS) error {
	// Create mountpoint
	log.Infof("Creating mountpoint %s", d.mountpoint)
	if err := CreateMountpointAt(ctx, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer RemoveMountpointAt(d.mountpoint)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", d.rootPartition, d.mountpoint)
	if err := MountPartitionAt(ctx, d.rootPartition, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer UnmountFilesystems(d.mountpoint)
	}

	// Create ESP mountpoint
	esp := path.Join(d.mountpoint, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := CreateMountpointAt(ctx, esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", d.espPartition, esp)
	if err := MountPartitionAt(ctx, d.espPartition, esp); err != nil {
		return err
	}

	if d.cleanup {
		defer UnmountFilesystems(esp)
	}

	// Copy source to root
	log.Infof("Copying source %s to %s", source, d.mountpoint)
	if err := UnpackTarball(ctx, source, d.mountpoint); err != nil {
		return err
	}

	// Setup chroot executor
	log.Infof("Setting up chroot")
	chrootExecutor := chroot.NewExecutor(ctx, d.mountpoint)
	defer chrootExecutor.Teardown()

	// Configure filesystems
	log.Infof("Configuring filesystems")
	err := ConfigureFilesystems(ctx, d.mountpoint, d.rootPartition, d.espPartition)
	if err != nil {
		return err
	}

	// Install bootloader
	log.Infof("Installing bootloader")
	bootloader := yakdOS.BootloaderInstaller(
		d.DevicePath, d.mountpoint, chrootExecutor)
	if err := bootloader.Install(ctx); err != nil {
		return err
	}

	return nil
}

func identifyPartition(devicePath string, number int) (string, error) {
	v1 := fmt.Sprintf("%sp%d", devicePath, number)
	v2 := fmt.Sprintf("%s%d", devicePath, number)

	log.Debugf("looking for partition %d on %s (trying %s, %s)",
		number, devicePath, v1, v2)

	if _, err := os.Stat(v1); err == nil {
		return v1, nil
	}

	if _, err := os.Stat(v2); err == nil {
		return v2, nil
	}

	return "", fmt.Errorf(
		"failed to identify partition %d on %s", number, devicePath)
}
