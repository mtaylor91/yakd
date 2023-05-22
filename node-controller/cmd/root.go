package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "node-controller",
}

func init() {
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(ImageCmd)
}

func Main() {
	RootCmd.Execute()
}
