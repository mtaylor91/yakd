package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/bootstrap"
	"github.com/mtaylor91/yakd/pkg/debian"
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image management",
}

var ImageStage1Cmd = &cobra.Command{
	Use:   "stage1",
	Short: "Stage 1 of image build",
	Run:   ImageStage1,
}

var ImageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an image",
	Run:   ImageBuild,
}

func init() {
	ImageCmd.AddCommand(ImageBuildCmd)
	ImageCmd.AddCommand(ImageStage1Cmd)

	stage1Flags := ImageStage1Cmd.Flags()
	stage1Flags.BoolP("force", "f", false, "Overwrite existing stage 1")
	stage1Flags.String("target", "stage1.tar.gz", "Target path for stage 1")
	stage1Flags.String("suite", debian.DefaultSuite, "Debian suite")
	stage1Flags.String("mirror", debian.DefaultMirror, "Debian mirror")
	stage1Flags.String("mountpoint", "/mnt/target", "Mountpoint for stage 1 build")
	stage1Flags.Int("tmpfs-size", 4096, "tmpfs size in MB")
	stage1Flags.Bool("no-cleanup", false, "Do not cleanup after build")

	imageBuildFlags := ImageBuildCmd.Flags()
	imageBuildFlags.String("target", "debian.img", "Target path for image")
	imageBuildFlags.String("suite", debian.DefaultSuite, "Debian suite")
	imageBuildFlags.String("mirror", debian.DefaultMirror, "Debian mirror")
	imageBuildFlags.String("mountpoint", "/mnt/target", "Mountpoint for image build")
	imageBuildFlags.Int("size", 8192, "Image size in MB")
	imageBuildFlags.BoolP("force", "f", false, "Overwrite existing image")
	imageBuildFlags.Bool("no-cleanup", false, "Do not cleanup after build")
}

func ImageStage1(cmd *cobra.Command, args []string) {
	log.Info("Building stage 1")

	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("target", f.Lookup("target"))
	v.BindPFlag("suite", f.Lookup("suite"))
	v.BindPFlag("mirror", f.Lookup("mirror"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("tmpfs-size", f.Lookup("tmpfs-size"))
	v.BindPFlag("no-cleanup", f.Lookup("no-cleanup"))

	cleanup := !v.GetBool("no-cleanup")

	debian := &debian.Debian{}
	debian.Suite = v.GetString("suite")
	debian.Mirror = v.GetString("mirror")

	stage1 := &bootstrap.Stage1{
		Source: v.GetString("mountpoint"),
		Target: v.GetString("target"),
	}

	tmpfs := &bootstrap.TmpFS{
		Path:   v.GetString("mountpoint"),
		SizeMB: v.GetInt("tmpfs-size"),
	}

	err := tmpfs.Bootstrap(debian, cleanup)
	if err != nil {
		log.Fatal(err)
	}

	if cleanup {
		defer tmpfs.Destroy()
	}

	err = stage1.BuildArchive()
	if err != nil {
		log.Fatal(err)
	}
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
