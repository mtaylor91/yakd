package util

import (
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

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

// MountMetadataFilesystems creates the mountpoints for the bootstrap
func MountMetadataFilesystems(root string) error {
	mount, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	commands := []*exec.Cmd{
		exec.Command(mount, "--rbind", "/dev", path.Join(root, "dev")),
		exec.Command(mount, "--make-rslave", path.Join(root, "dev")),
		exec.Command(mount, "-t", "proc", "/proc", path.Join(root, "proc")),
		exec.Command(mount, "--rbind", "/sys", path.Join(root, "sys")),
		exec.Command(mount, "--make-rslave", path.Join(root, "sys")),
		exec.Command(mount, "--bind", "/run", path.Join(root, "run")),
		exec.Command(mount, "--make-slave", path.Join(root, "run")),
	}

	for _, cmd := range commands {
		if err = cmd.Run(); err != nil {
			return err
		}
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

// UnmountFilesystems recursively unmounts the specified filesystem(s)
func UnmountFilesystems(p string) {
	// Unmount filesystem
	if err := RunCmd("umount", "-R", p); err != nil {
		log.Errorf("Unmount %s failed: %s", p, err)
	}
}

// UnmountMetadataFilesystems destroys the mountpoints for the bootstrap
func UnmountMetadataFilesystems(root string) {
	umount, err := exec.LookPath("umount")
	if err != nil {
		log.Errorf("Could not find umount: %s", err)
	}

	commands := []*exec.Cmd{
		exec.Command(umount, "-R", path.Join(root, "dev")),
		exec.Command(umount, "-R", path.Join(root, "proc")),
		exec.Command(umount, "-R", path.Join(root, "sys")),
		exec.Command(umount, "-R", path.Join(root, "run")),
	}

	for _, cmd := range commands {
		if err = cmd.Run(); err != nil {
			log.Errorf("Unmount %s failed: %s", cmd, err)
		}
	}
}
