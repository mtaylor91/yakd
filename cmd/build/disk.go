package build

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/disk"
)

func init() {
	f := Disk.Flags()
	f.String("os", "debian", "Operating system")
	f.String("mountpoint", "build/mount", "Path to mount the target filesystem")
	f.String("stage1-template", "build/{{.OS}}/yakd-stage1-{{.Arch}}.tar.gz",
		"Path template for stage 1 archive")
}

var Disk = &cobra.Command{
	Use:   "disk TARGET",
	Short: "Build a disk image",
	Args:  cobra.ExactArgs(1),
	Run:   BuildDisk,
}

func BuildDisk(cmd *cobra.Command, args []string) {
	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("os", f.Lookup("os"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("stage1-template", f.Lookup("stage1-template"))

	var config disk.Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	ctx := cmd.Context()
	config.Target = args[0]
	if err := config.BuildDisk(ctx); err != nil {
		log.Fatal(err)
	}
}
