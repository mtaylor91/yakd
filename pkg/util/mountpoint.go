package util

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/system"

	log "github.com/sirupsen/logrus"
)

// CreateMountpoint creates the mountpoint for the bootstrap
func CreateMountpointAt(ctx context.Context, path string) error {
	// Create mountpoint if it doesn't exist
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("mkdir", "-p", path); err != nil {
		return err
	}

	return nil
}

// RemoveMountpointAt removes the specified mountpoint
func RemoveMountpointAt(p string) {
	// Remove mountpoint
	if err := system.Local.RunCommand("rmdir", p); err != nil {
		log.Errorf("Remove %s failed: %s", p, err)
	}
}

// Mount mounts the specified disk at the specified location
func Mount(ctx context.Context, device, location string) error {
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("mount", device, location); err != nil {
		return err
	}

	return nil
}

// Unmount unmounts the filesystem(s) as the specified location
func Unmount(ctx context.Context, p string) {
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("umount", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}

// UnmountRecursive recursively unmounts the filesystem(s) as the specified location
func UnmountRecursive(p string) {
	// Unmount filesystem
	if err := system.Local.RunCommand("umount", "-R", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}
