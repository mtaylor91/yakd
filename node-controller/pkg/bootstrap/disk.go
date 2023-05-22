package bootstrap

import (
	"path"

	log "github.com/sirupsen/logrus"
)

// Disk represents the bootstrap configuration for a disk
type Disk struct {
	Path          string
	ESPPartition  string
	RootPartition string
	Mount         string
	Cleanup       bool
	OS            OS
}

// Bootstrap runs filesystem bootstrapping
func (d *Disk) Bootstrap() error {
	// Create mountpoint
	log.Infof("Creating mountpoint %s", d.Mount)
	if err := CreateMountpointAt(d.Mount); err != nil {
		return err
	}

	if d.Cleanup {
		defer RemoveMountpointAt(d.Mount)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", d.RootPartition, d.Mount)
	if err := MountPartitionAt(d.RootPartition, d.Mount); err != nil {
		return err
	}

	if d.Cleanup {
		defer UnmountFilesystems(d.Mount)
	}

	// Bootstrap OS root filesystem
	if err := d.OS.Bootstrap(); err != nil {
		return err
	}

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems on %s", d.Mount)
	if err := MountMetadataFilesystems(d.Mount); err != nil {
		return err
	}

	if d.Cleanup {
		defer UnmountMetadataFilesystems(d.Mount)
	}

	// Create ESP mountpoint
	esp := path.Join(d.Mount, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := CreateMountpointAt(esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", d.ESPPartition, esp)
	if err := MountPartitionAt(d.ESPPartition, esp); err != nil {
		return err
	}

	if d.Cleanup {
		defer UnmountFilesystems(esp)
	}

	// Run post-bootstrap step
	log.Infof("Running post-bootstrap step")
	if err := d.OS.PostBootstrap(); err != nil {
		return err
	}

	return nil
}
