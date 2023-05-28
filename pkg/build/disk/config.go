package disk

// Config for the disk builder
type Config struct {
	// Operating system
	OS string `mapstructure:"os"`
	// Path template for the stage 1 archive
	Stage1Template string `mapstructure:"stage1-template"`
	// Path to the target disk
	Target string `mapstructure:"target"`
	// Path to mount the target filesystem
	Mountpoint string `mapstructure:"mountpoint"`
}
