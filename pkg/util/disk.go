package util

import (
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/os"
)

// Disk represents the bootstrap configuration for a disk
type Disk struct {
	DevicePath    string
	espPartition  string
	rootPartition string
	mountpoint    string
	cleanup       bool
}

// NewDisk initializes a new Disk struct
func NewDisk(devicePath, mountpoint string, cleanup bool) *Disk {
	return &Disk{
		DevicePath:    devicePath,
		espPartition:  devicePath + "p1",
		rootPartition: devicePath + "p2",
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

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems")
	if err := MountMetadataFilesystems(d.mountpoint); err != nil {
		return err
	}

	defer UnmountMetadataFilesystems(d.mountpoint)

	// Configure filesystems
	log.Infof("Configuring filesystems")
	err := ConfigureFilesystems(d.mountpoint, d.rootPartition, d.espPartition)
	if err != nil {
		return err
	}

	// Install bootloader
	log.Infof("Installing bootloader")
	bootloader := os.Bootloader(d.mountpoint)
	if err := bootloader.Install(d.DevicePath); err != nil {
		return err
	}

	return nil
}
