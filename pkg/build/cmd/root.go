package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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
	f := cmd.Flags()

	debug, _ := f.GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging enabled")
	}

	trace, _ := f.GetBool("trace")
	if trace {
		log.SetLevel(log.TraceLevel)
		log.Trace("Trace logging enabled")
	}
}
