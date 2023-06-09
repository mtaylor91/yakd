package disk

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/mtaylor91/yakd/pkg/build/release"
	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/log"
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
	sys := system.Local.WithContext(ctx)
	err := sys.RunCommand("mkfs.vfat", "-F", "32", d.espPartition)
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	if err := sys.RunCommand("mkfs.ext4", d.rootPartition); err != nil {
		return err
	}

	return nil
}

// Populate the disk from the specified source
func (d *Disk) Populate(ctx context.Context, source string, os release.OS) error {
	log := log.FromContext(ctx)

	// Create mountpoint
	log.Infof("Creating mountpoint %s", d.mountpoint)
	if err := util.CreateMountpointAt(ctx, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer util.RemoveMountpointAt(d.mountpoint)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", d.rootPartition, d.mountpoint)
	if err := util.Mount(ctx, d.rootPartition, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer util.UnmountRecursive(d.mountpoint)
	}

	// Create ESP mountpoint
	esp := path.Join(d.mountpoint, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := util.CreateMountpointAt(ctx, esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", d.espPartition, esp)
	if err := util.Mount(ctx, d.espPartition, esp); err != nil {
		return err
	}

	if d.cleanup {
		defer util.UnmountRecursive(esp)
	}

	// Copy source to root
	log.Infof("Copying source %s to %s", source, d.mountpoint)
	if err := util.UnpackTarball(ctx, source, d.mountpoint); err != nil {
		return err
	}

	// Setup chroot executor
	log.Infof("Setting up chroot")
	sys := system.Local.WithContext(ctx)
	chrootSys := system.Chroot(sys, d.mountpoint)
	if err := chrootSys.Setup(); err != nil {
		return err
	}

	defer chrootSys.Teardown()

	// Configure filesystems
	log.Infof("Configuring filesystems")
	err := util.ConfigureFilesystems(
		ctx, d.mountpoint, d.rootPartition, d.espPartition)
	if err != nil {
		return err
	}

	// Install bootloader
	log.Infof("Installing bootloader")
	bootloader := os.BootloaderInstaller(d.DevicePath, d.mountpoint, chrootSys)
	if err := bootloader.Install(ctx); err != nil {
		return err
	}

	return nil
}

func identifyPartition(devicePath string, number int) (string, error) {
	v1 := fmt.Sprintf("%sp%d", devicePath, number)
	v2 := fmt.Sprintf("%s%d", devicePath, number)

	if _, err := os.Stat(v1); err == nil {
		return v1, nil
	}

	if _, err := os.Stat(v2); err == nil {
		return v2, nil
	}

	return "", fmt.Errorf(
		"failed to identify partition %d on %s", number, devicePath)
}
