package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mtaylor91/yakd/cmd/build"
)

func init() {
	Root.AddCommand(build.Root)
	f := Root.PersistentFlags()
	f.Bool("debug", false, "Enable debug logging")
	f.Bool("trace", false, "Enable trace logging")
}

var Root = &cobra.Command{
	Use:              "yakd",
	PersistentPreRun: ConfigureRoot,
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
