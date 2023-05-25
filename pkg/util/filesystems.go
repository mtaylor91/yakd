package util

import (
	"html/template"
	"os"
	"path"
	"strings"
)

const fstabTemplate = `# <filesystem> <mountpoint> <type> <options> <dump> <pass>
UUID={{.RootPartitionUUID}} / ext4 defaults 0 1
UUID={{.ESPPartitionUUID}} /boot/efi vfat defaults 0 1
`

// ConfigureFilesystems configures the filesystems on the specified disk
func ConfigureFilesystems(mountpoint, rootPartition, espPartition string) error {
	rootPartitionUUID, err := GetFilesystemUUID(rootPartition)
	if err != nil {
		return err
	}

	espPartitionUUID, err := GetFilesystemUUID(espPartition)
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
func GetFilesystemUUID(devicePath string) (string, error) {
	blkidOutput, err := GetOutput("blkid", "-s", "UUID", "-o", "value", devicePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(blkidOutput)), nil
}
