package stage1

type Stage1 struct {
	DebianMirror       string `mapstructure:"debian-mirror"`
	DebianSuite        string `mapstructure:"debian-suite"`
	Force              bool   `mapstructure:"force"`
	GentooBinPkgsCache string `mapstructure:"gentoo-binpkgs-cache"`
	GentooStage3       string `mapstructure:"gentoo-stage3"`
	Mountpoint         string `mapstructure:"mountpoint"`
	OS                 string `mapstructure:"os"`
	TargetTemplate     string `mapstructure:"target-template"`
	TmpFSSize          int    `mapstructure:"tmpfs-size"`
}
