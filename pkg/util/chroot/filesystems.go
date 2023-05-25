package chroot

import (
	"os/exec"
	"path"

	log "github.com/sirupsen/logrus"
)

// MountMetadataFilesystems creates the mountpoints for the bootstrap
func MountMetadataFilesystems(root string) error {
	mount, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	commands := []*exec.Cmd{
		exec.Command(mount, "-t", "proc", "/proc", path.Join(root, "proc")),
		exec.Command(mount, "--rbind", "/dev", path.Join(root, "dev")),
		exec.Command(mount, "--make-rslave", path.Join(root, "dev")),
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

// UnmountMetadataFilesystems destroys the mountpoints for the bootstrap
func UnmountMetadataFilesystems(root string) {
	umount, err := exec.LookPath("umount")
	if err != nil {
		log.Errorf("Could not find umount: %s", err)
	}

	commands := []*exec.Cmd{
		exec.Command(umount, "-R", path.Join(root, "run")),
		exec.Command(umount, "-R", path.Join(root, "sys")),
		exec.Command(umount, "-R", path.Join(root, "dev")),
		exec.Command(umount, "-R", path.Join(root, "proc")),
	}

	for _, cmd := range commands {
		if err = cmd.Run(); err != nil {
			log.Errorf("Unmount %s failed: %s", cmd, err)
		}
	}
}
