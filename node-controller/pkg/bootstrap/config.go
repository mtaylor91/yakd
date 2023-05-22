package bootstrap

const (
	DefaultTargetMountpoint = "/mnt/target"
)

// BootstrapConfig represents the configuration for a bootstrap
type BootstrapConfig struct {
	Cleanup       bool
	Disk          string
	ESPPartition  string
	RootPartition string
	Mount         string
	OS            OS
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(disk, esp, root, mount string, os OSFactory) *BootstrapConfig {
	return &BootstrapConfig{true, disk, esp, root, mount, os.NewOS(mount)}
}
