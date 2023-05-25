package build

import (
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "build",
	Short: "Build commands",
}

func init() {
	Root.AddCommand(Disk)
	Root.AddCommand(Image)
	Root.AddCommand(Stage1)
}
