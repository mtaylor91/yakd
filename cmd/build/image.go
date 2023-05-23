package build

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/image"
)

func init() {
	f := ImageCmd.Flags()
	f.BoolP("force", "f", false, "Overwrite existing image")
	f.Int("size", 4096, "Image size in MB")
	f.String("stage1", "stage1.tar.gz", "Path to stage 1 tarball")
	f.String("target", "yakd.qcow2", "Target path for image")
	f.String("mountpoint", "/mnt/target", "Mountpoint for image build")
}

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image from stage1 tarball",
	Run:   BuildImage,
}

func BuildImage(cmd *cobra.Command, args []string) {
	log.Info("Building image")

	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("size", f.Lookup("size"))
	v.BindPFlag("stage1", f.Lookup("stage1"))
	v.BindPFlag("target", f.Lookup("target"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))

	force := v.GetBool("force")
	size := v.GetInt("size")
	stage1 := v.GetString("stage1")
	target := v.GetString("target")
	mountpoint := v.GetString("mountpoint")

	err := image.BuildImage(force, size, stage1, target, mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
