package util

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
	log "github.com/sirupsen/logrus"
)

// CreateMountpoint creates the mountpoint for the bootstrap
func CreateMountpointAt(ctx context.Context, path string) error {
	// Create mountpoint if it doesn't exist
	if err := executor.RunCmd(ctx, "mkdir", "-p", path); err != nil {
		return err
	}

	return nil
}

// RemoveMountpointAt removes the specified mountpoint
func RemoveMountpointAt(p string) {
	// Remove mountpoint
	ctx := context.Background()
	if err := executor.RunCmd(ctx, "rmdir", p); err != nil {
		log.Errorf("Remove %s failed: %s", p, err)
	}
}

// MountPartitionAt mounts the specified disk at the specified location
func MountPartitionAt(ctx context.Context, partition, location string) error {
	// Mount filesystem
	if err := executor.RunCmd(ctx, "mount", partition, location); err != nil {
		return err
	}

	return nil
}

// UnmountFilesystems recursively unmounts the specified filesystem(s)
func UnmountFilesystems(p string) {
	// Unmount filesystem
	ctx := context.Background()
	if err := executor.RunCmd(ctx, "umount", "-R", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}
