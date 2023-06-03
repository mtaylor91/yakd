package util

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/mtaylor91/yakd/pkg/system"
)

const fstabTemplate = `# <filesystem> <mountpoint> <type> <options> <dump> <pass>
UUID={{.RootPartitionUUID}} / ext4 defaults 0 1
UUID={{.ESPPartitionUUID}} /boot/efi vfat defaults 0 1
`

// ConfigureFilesystems configures the filesystems on the specified disk
func ConfigureFilesystems(ctx context.Context, mountpoint, rootPartition, espPartition string) error {
	rootPartitionUUID, err := GetFilesystemUUID(ctx, rootPartition)
	if err != nil {
		return err
	}

	espPartitionUUID, err := GetFilesystemUUID(ctx, espPartition)
	if err != nil {
		return err
	}

	t, err := template.New("fstab").Parse(fstabTemplate)
	if err != nil {
		return err
	}

	fstabPath := path.Join(mountpoint, "etc", "fstab")
	fstabFile, err := os.Create(fstabPath)
	if err != nil {
		return err
	}

	defer fstabFile.Close()

	return t.Execute(fstabFile, struct {
		RootPartitionUUID string
		ESPPartitionUUID  string
	}{
		rootPartitionUUID,
		espPartitionUUID,
	})
}

// GetFilesystemUUID returns the UUID of the specified filesystem
func GetFilesystemUUID(ctx context.Context, devicePath string) (string, error) {
	var blkidOutput bytes.Buffer
	sys := system.Local.WithContext(ctx)
	err := sys.RunCommandWithOutput(
		&blkidOutput, "blkid", "-s", "UUID", "-o", "value", devicePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(blkidOutput.String()), nil
}
