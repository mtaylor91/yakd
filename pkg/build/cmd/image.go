package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/image"
	"github.com/mtaylor91/yakd/pkg/build/release/gentoo"
	"github.com/mtaylor91/yakd/pkg/util/log"
)

func init() {
	f := Image.Flags()
	f.BoolP("force", "f", false, "Overwrite existing image")
	f.String("format", "iso", "Image format")
	f.String("gentoo-binpkgs-cache", gentoo.Default.BinPkgsCache,
		"Path to Gentoo binpkgs cache")
	f.String("mountpoint", "build/mount", "Mountpoint for image build")
	f.String("os", "debian", "Operating system")
	f.String("stage1-template", "build/{{.OS}}/yakd-stage1-{{.Arch}}.tar.gz",
		"Path template for stage 1 tarball")
	f.Int("size-mb", 10240, "Image size in MB")
	f.String("target-template", "build/{{.OS}}/yakd.{{.Format}}",
		"Target path template for image")
}

var Image = &cobra.Command{
	Use:   "image",
	Short: "Build image from stage1 tarball",
	Run:   BuildImage,
}

func BuildImage(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	log := log.FromContext(ctx)

	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("format", f.Lookup("format"))
	v.BindPFlag("gentoo-binpkgs-cache", f.Lookup("gentoo-binpkgs-cache"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("os", f.Lookup("os"))
	v.BindPFlag("stage1-template", f.Lookup("stage1-template"))
	v.BindPFlag("size-mb", f.Lookup("size-mb"))
	v.BindPFlag("target-template", f.Lookup("target-template"))

	var config image.Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	err := config.BuildImage(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
