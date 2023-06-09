package image

type Config struct {
	// Force overwrite of existing image
	Force bool `mapstructure:"force"`
	// Image format
	Format string `mapstructure:"format"`
	// Path to the Gentoo binpkgs cache
	GentooBinPkgsCache string `mapstructure:"gentoo-binpkgs-cache"`
	// Path to mount the target filesystem
	Mountpoint string `mapstructure:"mountpoint"`
	// Operating system
	OS string `mapstructure:"os"`
	// Path template for the stage 1 archive
	Stage1Template string `mapstructure:"stage1-template"`
	// Image size in megabytes
	SizeMB int `mapstructure:"size-mb"`
	// Path template for the target image
	TargetTemplate string `mapstructure:"target-template"`
}
