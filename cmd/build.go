package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mtaylor91/yakd/cmd/build"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build management",
}

func init() {
	BuildCmd.AddCommand(build.ImageCmd)
	BuildCmd.AddCommand(build.Stage1Cmd)
}
