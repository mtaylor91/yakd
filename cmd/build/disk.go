package build

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/disk"
)

func init() {
	f := Disk.Flags()
	f.String("stage1", "stage1.tar.gz", "Path to stage1 tarball")
	f.String("mountpoint", "/mnt/target", "Path to mount the target filesystem")
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

	v.BindPFlag("stage1", f.Lookup("stage1"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("no-cleanup", f.Lookup("no-cleanup"))

	target := args[0]
	stage1 := v.GetString("stage1")
	mountpoint := v.GetString("mountpoint")

	ctx := cmd.Context()
	if err := disk.BuildDisk(ctx, target, stage1, mountpoint); err != nil {
		log.Fatal(err)
	}
}
