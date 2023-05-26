package build

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mtaylor91/yakd/pkg/build/stage1"
	"github.com/mtaylor91/yakd/pkg/debian"
)

func init() {
	f := Stage1.Flags()
	f.BoolP("force", "f", false, "Overwrite existing stage 1")
	f.String("target", "stage1.tar.gz", "Target path for stage 1")
	f.String("suite", debian.DefaultSuite, "Debian suite")
	f.String("mirror", debian.DefaultMirror, "Debian mirror")
	f.String("mountpoint", "/mnt/target", "Mountpoint for stage 1 build")
	f.Int("tmpfs-size", 4096, "tmpfs size in MB")
}

var Stage1 = &cobra.Command{
	Use:   "stage1",
	Short: "Stage 1 of image build",
	Run:   BuildStage1,
}

func BuildStage1(cmd *cobra.Command, args []string) {
	log.Info("Building stage 1")

	f := cmd.Flags()
	v := viper.New()

	v.BindPFlag("force", f.Lookup("force"))
	v.BindPFlag("target", f.Lookup("target"))
	v.BindPFlag("suite", f.Lookup("suite"))
	v.BindPFlag("mirror", f.Lookup("mirror"))
	v.BindPFlag("mountpoint", f.Lookup("mountpoint"))
	v.BindPFlag("tmpfs-size", f.Lookup("tmpfs-size"))

	force := v.GetBool("force")
	target := v.GetString("target")
	suite := v.GetString("suite")
	mirror := v.GetString("mirror")
	mountpoint := v.GetString("mountpoint")
	tmpfsSize := v.GetInt("tmpfs-size")

	ctx := cmd.Context()
	err := stage1.BuildStage1(
		ctx, force, target, suite, mirror, mountpoint, tmpfsSize)
	if err != nil {
		log.Fatal(err)
	}
}
