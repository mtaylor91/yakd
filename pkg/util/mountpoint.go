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

// Mount mounts the specified disk at the specified location
func Mount(ctx context.Context, device, location string) error {
	// Mount filesystem
	if err := executor.RunCmd(ctx, "mount", device, location); err != nil {
		return err
	}

	return nil
}

// Unmount unmounts the filesystem(s) as the specified location
func Unmount(ctx context.Context, p string) {
	// Unmount filesystem
	if err := executor.RunCmd(ctx, "umount", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}

// UnmountRecursive recursively unmounts the filesystem(s) as the specified location
func UnmountRecursive(p string) {
	// Unmount filesystem
	ctx := context.Background()
	if err := executor.RunCmd(ctx, "umount", "-R", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}
