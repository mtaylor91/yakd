package util

// CreateMountpoint creates the mountpoint for the bootstrap
func CreateMountpointAt(path string) error {
	// Create mountpoint if it doesn't exist
	if err := RunCmd("mkdir", "-p", path); err != nil {
		return err
	}

	return nil
}

// MountPartitionAt mounts the specified disk at the specified location
func MountPartitionAt(partition, location string) error {
	// Mount filesystem
	if err := RunCmd("mount", partition, location); err != nil {
		return err
	}

	return nil
}
