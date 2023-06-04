package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mtaylor91/yakd/pkg/util/log"
)

var Root = &cobra.Command{
	Use:              "yakd-build",
	PersistentPreRun: ConfigureRoot,
}

func init() {
	Root.AddCommand(Image)
	Root.AddCommand(Stage1)

	flags := Root.PersistentFlags()
	flags.Bool("debug", false, "Enable debug logging")
	flags.Bool("trace", false, "Enable trace logging")
}

func ConfigureRoot(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	f := cmd.Flags()

	debug, _ := f.GetBool("debug")
	if debug {
		log.DefaultLogger.SetLevel(log.DebugLevel)
		log.FromContext(ctx).Debug("Debug logging enabled")
	}

	trace, _ := f.GetBool("trace")
	if trace {
		log.DefaultLogger.SetLevel(log.TraceLevel)
		log.FromContext(ctx).Trace("Trace logging enabled")
	}
}
