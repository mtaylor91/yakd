package bootstrap

const (
	DefaultTargetMountpoint = "/mnt/target"
)

// BootstrapConfig represents the configuration for a bootstrap
type BootstrapConfig struct {
	Cleanup       bool
	ESPPartition  string
	RootPartition string
	Mount         string
	OS            OS
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(esp, root, mount string, os OSFactory) *BootstrapConfig {
	return &BootstrapConfig{true, esp, root, mount, os.NewOS(mount)}
}
