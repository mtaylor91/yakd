package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/release/debian"
	"github.com/mtaylor91/yakd/pkg/build/release/gentoo"
	"github.com/mtaylor91/yakd/pkg/build/stage1"
	"github.com/mtaylor91/yakd/pkg/util/log"
)

func init() {
	f := Stage1.Flags()
	f.String("os", "debian", "Operating system")
	f.BoolP("force", "f", false, "Overwrite existing stage 1")
	f.String("debian-mirror", debian.Default.Mirror, "Debian mirror")
	f.String("debian-suite", debian.Default.Suite, "Debian suite")
	f.String("gentoo-binpkgs-cache", gentoo.Default.BinPkgsCache,
		"Path to Gentoo binpkgs cache")
	f.String("gentoo-stage3", "build/gentoo/upstream-stage3.tar.xz",
		"Gentoo stage3 archive path")
	f.String("mountpoint", "build/mount", "Mountpoint for stage 1 build")
	f.String("target-template", "build/{{.OS}}/yakd-stage1-{{.Arch}}.tar.gz",
		"Target path template for stage 1 archive")
	f.Int("tmpfs-size", 10240, "tmpfs size in MB")
}

var Stage1 = &cobra.Command{
	Use:   "stage1",
	Short: "Build stage1 filesystem tarball",
	Run:   BuildStage1,
}

func BuildStage1(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	log := log.FromContext(ctx)

	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("os", f.Lookup("os"))
	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("debian-mirror", f.Lookup("debian-mirror"))
	v.BindPFlag("debian-suite", f.Lookup("debian-suite"))
	v.BindPFlag("gentoo-binpkgs-cache", f.Lookup("gentoo-binpkgs-cache"))
	v.BindPFlag("gentoo-stage3", f.Lookup("gentoo-stage3"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("target-template", f.Lookup("target-template"))
	v.BindPFlag("tmpfs-size", f.Lookup("tmpfs-size"))

	var stage1 stage1.Stage1
	err := v.Unmarshal(&stage1)
	if err != nil {
		log.Fatal(err)
	}

	if err := stage1.Build(ctx); err != nil {
		log.Fatal(err)
	}
}
