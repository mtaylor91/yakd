package stage1

type Stage1 struct {
	DebianMirror   string `mapstructure:"debian-mirror"`
	DebianSuite    string `mapstructure:"debian-suite"`
	Force          bool   `mapstructure:"force"`
	Mountpoint     string `mapstructure:"mountpoint"`
	OS             string `mapstructure:"os"`
	TargetTemplate string `mapstructure:"target-template"`
	TmpFSSize      int    `mapstructure:"tmpfs-size"`
}
