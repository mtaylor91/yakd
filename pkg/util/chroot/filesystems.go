package chroot

import (
	"context"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// MountMetadataFilesystems creates the mountpoints for the bootstrap
func MountMetadataFilesystems(ctx context.Context, root string) error {
	commands := [][]string{
		[]string{"mount", "-t", "proc", "/proc", path.Join(root, "proc")},
		[]string{"mount", "--rbind", "/dev", path.Join(root, "dev")},
		[]string{"mount", "--make-rslave", path.Join(root, "dev")},
		[]string{"mount", "--rbind", "/sys", path.Join(root, "sys")},
		[]string{"mount", "--make-rslave", path.Join(root, "sys")},
		[]string{"mount", "--bind", "/run", path.Join(root, "run")},
		[]string{"mount", "--make-slave", path.Join(root, "run")},
	}

	return executor.RunCmdList(ctx, executor.Default, commands...)
}

// UnmountMetadataFilesystems destroys the mountpoints for the bootstrap
func UnmountMetadataFilesystems(ctx context.Context, root string) {
	commands := [][]string{
		[]string{"umount", "-R", path.Join(root, "run")},
		[]string{"umount", "-R", path.Join(root, "sys")},
		[]string{"umount", "-R", path.Join(root, "dev")},
		[]string{"umount", "-R", path.Join(root, "proc")},
	}

	err := executor.RunCmdList(ctx, executor.Default, commands...)
	if err != nil {
		log.Errorf("Error running umount: %s", err)
	}
}
