package bootstrap

import (
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

// Bootstrap runs filesystem bootstrapping
func (c *BootstrapConfig) Bootstrap() error {
	// Create mountpoint
	log.Infof("Creating mountpoint %s", c.Mount)
	if err := CreateMountpointAt(c.Mount); err != nil {
		return err
	}

	if c.Cleanup {
		defer RemoveMountpointAt(c.Mount)
	}

	// Mount root partition
	log.Infof("Mounting root partition %s on %s", c.RootPartition, c.Mount)
	if err := MountPartitionAt(c.RootPartition, c.Mount); err != nil {
		return err
	}

	if c.Cleanup {
		defer UnmountFilesystems(c.Mount)
	}

	// Bootstrap OS root filesystem
	if err := c.OS.Bootstrap(); err != nil {
		return err
	}

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems on %s", c.Mount)
	if err := c.MountMetadataFilesystems(); err != nil {
		return err
	}

	if c.Cleanup {
		defer c.UnmountMetadataFilesystems()
	}

	// Create ESP mountpoint
	esp := path.Join(c.Mount, "boot", "efi")
	log.Infof("Creating ESP mountpoint at %s", esp)
	if err := CreateMountpointAt(esp); err != nil {
		return err
	}

	// Mount ESP
	log.Infof("Mounting ESP partition %s on %s", c.ESPPartition, esp)
	if err := MountPartitionAt(c.ESPPartition, esp); err != nil {
		return err
	}

	if c.Cleanup {
		defer UnmountFilesystems(esp)
	}

	// Run post-bootstrap step
	log.Infof("Running post-bootstrap step")
	if err := c.OS.PostBootstrap(); err != nil {
		return err
	}

	return nil
}

// CreateMountpoint creates the mountpoint for the bootstrap
func CreateMountpointAt(path string) error {
	mkdir, err := exec.LookPath("mkdir")
	if err != nil {
		return err
	}

	// Create mountpoint if it doesn't exist
	cmd := exec.Command(mkdir, "-p", path)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// MountPartitionAt mounts the specified disk at the specified location
func MountPartitionAt(partition, location string) error {
	mount, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	// Mount filesystem
	cmd := exec.Command(mount, partition, location)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// MountMetadataFilesystems creates the mountpoints for the bootstrap
func (c *BootstrapConfig) MountMetadataFilesystems() error {
	mount, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	commands := []*exec.Cmd{
		exec.Command(mount, "--rbind", "/dev", path.Join(c.Mount, "dev")),
		exec.Command(mount, "--make-rslave", path.Join(c.Mount, "dev")),
		exec.Command(mount, "-t", "proc", "/proc", path.Join(c.Mount, "proc")),
		exec.Command(mount, "--rbind", "/sys", path.Join(c.Mount, "sys")),
		exec.Command(mount, "--make-rslave", path.Join(c.Mount, "sys")),
		exec.Command(mount, "--bind", "/run", path.Join(c.Mount, "run")),
		exec.Command(mount, "--make-slave", path.Join(c.Mount, "run")),
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
	rmdir, err := exec.LookPath("rmdir")
	if err != nil {
		log.Errorf("Could not find rmdir: %s", err)
	}

	// Remove mountpoint
	cmd := exec.Command(rmdir, p)
	if err := cmd.Run(); err != nil {
		log.Errorf("Remove %s failed: %s", cmd, err)
	}
}

// UnmountFilesystems recursively unmounts the specified filesystem(s)
func UnmountFilesystems(p string) {
	umount, err := exec.LookPath("umount")
	if err != nil {
		log.Errorf("Could not find umount: %s", err)
	}

	// Unmount filesystem
	cmd := exec.Command(umount, "-R", p)
	if err := cmd.Run(); err != nil {
		log.Errorf("Unmount %s failed: %s", cmd, err)
	}
}

// UnmountMetadataFilesystems destroys the mountpoints for the bootstrap
func (c *BootstrapConfig) UnmountMetadataFilesystems() {
	umount, err := exec.LookPath("umount")
	if err != nil {
		log.Errorf("Could not find umount: %s", err)
	}

	commands := []*exec.Cmd{
		exec.Command(umount, "-l", path.Join(c.Mount, "dev", "pts")),
		exec.Command(umount, "-l", path.Join(c.Mount, "dev", "shm")),
		exec.Command(umount, "-l", path.Join(c.Mount, "dev")),
	}

	for _, cmd := range commands {
		if err = cmd.Run(); err != nil {
			log.Errorf("Unmount %s failed: %s", cmd, err)
		}
	}
}
