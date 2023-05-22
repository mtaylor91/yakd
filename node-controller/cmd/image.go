package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/node-controller/pkg/bootstrap"
	"github.com/mtaylor91/yakd/node-controller/pkg/debian"
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image management",
}

var ImageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an image",
	Run:   ImageBuild,
}

func init() {
	ImageCmd.AddCommand(ImageBuildCmd)
	imageBuildFlags := ImageBuildCmd.Flags()
	imageBuildFlags.String("target", "debian.img", "Target path for image")
	imageBuildFlags.String("suite", debian.DefaultSuite, "Debian suite")
	imageBuildFlags.String("mirror", debian.DefaultMirror, "Debian mirror")
	imageBuildFlags.String("mountpoint", "/mnt/target", "Mountpoint for image build")
	imageBuildFlags.Int("size", 8192, "Image size in MB")
	imageBuildFlags.BoolP("force", "f", false, "Overwrite existing image")
	imageBuildFlags.Bool("no-cleanup", false, "Do not cleanup after build")
}

func ImageBuild(cmd *cobra.Command, args []string) {
	log.Info("Building image")

	f := cmd.Flags()
	v := viper.New()
	v.BindPFlag("target", f.Lookup("target"))
	v.BindPFlag("suite", f.Lookup("suite"))
	v.BindPFlag("mirror", f.Lookup("mirror"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("size", f.Lookup("size"))
	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("no-cleanup", f.Lookup("no-cleanup"))

	debian := &debian.Debian{}
	debian.Suite = v.GetString("suite")
	debian.Mirror = v.GetString("mirror")

	img := &bootstrap.Image{}
	img.Cleanup = !v.GetBool("no-cleanup")
	img.Path = v.GetString("target")
	img.SizeMB = v.GetInt("size")
	img.Overwrite = v.GetBool("force")

	err := img.Create(v.GetString("mountpoint"), debian)
	if err != nil {
		log.Fatal(err)
	}
}
