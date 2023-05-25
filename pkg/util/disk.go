package util

import (
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util/chroot"
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
func NewDisk(devicePath, mountpoint string, cleanup bool) *Disk {
	return &Disk{
		DevicePath:    devicePath,
		biosPartition: devicePath + "p1",
		espPartition:  devicePath + "p2",
		rootPartition: devicePath + "p3",
		mountpoint:    mountpoint,
		cleanup:       cleanup,
	}
}

// Populate populates the disk from the specified source
func (d *Disk) Populate(source string, os os.OS) error {
	// Create mountpoint
	log.Infof("Creating mountpoint %s", d.mountpoint)
	if err := CreateMountpointAt(d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer RemoveMountpointAt(d.mountpoint)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", d.rootPartition, d.mountpoint)
	if err := MountPartitionAt(d.rootPartition, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer UnmountFilesystems(d.mountpoint)
	}

	// Create ESP mountpoint
	esp := path.Join(d.mountpoint, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := CreateMountpointAt(esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", d.espPartition, esp)
	if err := MountPartitionAt(d.espPartition, esp); err != nil {
		return err
	}

	if d.cleanup {
		defer UnmountFilesystems(esp)
	}

	// Copy source to root
	log.Infof("Copying source %s to %s", source, d.mountpoint)
	if err := UnpackTarball(source, d.mountpoint); err != nil {
		return err
	}

	// Setup chroot executor
	log.Infof("Setting up chroot")
	chrootExecutor := chroot.NewExecutor(d.mountpoint)
	defer chrootExecutor.Teardown()

	// Configure filesystems
	log.Infof("Configuring filesystems")
	err := ConfigureFilesystems(d.mountpoint, d.rootPartition, d.espPartition)
	if err != nil {
		return err
	}

	// Install bootloader
	log.Infof("Installing bootloader")
	bootloader := os.BootloaderInstaller(d.DevicePath, d.mountpoint, chrootExecutor)
	if err := bootloader.Install(); err != nil {
		return err
	}

	return nil
}
