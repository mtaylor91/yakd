package bootstrap

import (
	"path"

	"github.com/mtaylor91/yakd/pkg/util"
	log "github.com/sirupsen/logrus"
)

// Disk represents the bootstrap configuration for a disk
type Disk struct {
	devicePath    string
	espPartition  string
	rootPartition string
	mountpoint    string
	cleanup       bool
}

// NewDisk initializes a new Disk struct
func NewDisk(devicePath, espPartition, rootPartition, mountpoint string, cleanup bool) *Disk {
	return &Disk{
		devicePath:    devicePath,
		espPartition:  espPartition,
		rootPartition: rootPartition,
		mountpoint:    mountpoint,
		cleanup:       cleanup,
	}
}

// Bootstrap runs filesystem bootstrapping
func (d *Disk) Bootstrap() error {
	// Create mountpoint
	log.Infof("Creating mountpoint %s", d.mountpoint)
	if err := util.CreateMountpointAt(d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer util.RemoveMountpointAt(d.mountpoint)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", d.rootPartition, d.mountpoint)
	if err := util.MountPartitionAt(d.rootPartition, d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer util.UnmountFilesystems(d.mountpoint)
	}

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems on %s", d.mountpoint)
	if err := util.MountMetadataFilesystems(d.mountpoint); err != nil {
		return err
	}

	if d.cleanup {
		defer util.UnmountMetadataFilesystems(d.mountpoint)
	}

	// Create ESP mountpoint
	esp := path.Join(d.mountpoint, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := util.CreateMountpointAt(esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", d.espPartition, esp)
	if err := util.MountPartitionAt(d.espPartition, esp); err != nil {
		return err
	}

	if d.cleanup {
		defer util.UnmountFilesystems(esp)
	}

	return nil
}
