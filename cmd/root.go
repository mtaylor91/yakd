package cmd

import (
	"github.com/mtaylor91/yakd/cmd/build"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use: "yakd",
}

func init() {
	Root.AddCommand(build.Root)
}
