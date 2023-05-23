package util

import (
	"path"

	log "github.com/sirupsen/logrus"
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
func NewDisk(devicePath, espPartition, rootPartition, mountpoint string, cleanup bool) *Disk {
	return &Disk{
		DevicePath:    devicePath,
		espPartition:  espPartition,
		rootPartition: rootPartition,
		mountpoint:    mountpoint,
		cleanup:       cleanup,
	}
}

// Populate populates the disk from the specified source
func (d *Disk) Populate(source string) error {
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

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems on %s", d.mountpoint)
	if err := MountMetadataFilesystems(d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer UnmountMetadataFilesystems(d.mountpoint)
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

	// TODO: install bootloader
	log.Infof("Installing bootloader")
	log.Errorf("Bootloader installation not implemented")

	return nil
}
