package util

import log "github.com/sirupsen/logrus"

// CreateMountpoint creates the mountpoint for the bootstrap
func CreateMountpointAt(path string) error {
	// Create mountpoint if it doesn't exist
	if err := RunCmd("mkdir", "-p", path); err != nil {
		return err
	}

	return nil
}

// RemoveMountpointAt removes the specified mountpoint
func RemoveMountpointAt(p string) {
	// Remove mountpoint
	if err := RunCmd("rmdir", p); err != nil {
		log.Errorf("Remove %s failed: %s", p, err)
	}
}

// MountPartitionAt mounts the specified disk at the specified location
func MountPartitionAt(partition, location string) error {
	// Mount filesystem
	if err := RunCmd("mount", partition, location); err != nil {
		return err
	}

	return nil
}

// UnmountFilesystems recursively unmounts the specified filesystem(s)
func UnmountFilesystems(p string) {
	// Unmount filesystem
	if err := RunCmd("umount", "-R", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}
