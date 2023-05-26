package bootstrap

import "github.com/mtaylor91/yakd/pkg/os"

const (
	DefaultTargetMountpoint = "/mnt/target"
)

// BootstrapConfig represents the configuration for a bootstrap
type BootstrapConfig struct {
	Cleanup              bool
	Disk                 string
	ESPPartition         string
	RootPartition        string
	Mount                string
	OSBootstrapInstaller os.OSBootstrapInstaller
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(disk, esp, root, mount string, os os.OS) *BootstrapConfig {
	return &BootstrapConfig{
		true, disk, esp, root, mount,
		os.BootstrapInstaller(mount),
	}
}
