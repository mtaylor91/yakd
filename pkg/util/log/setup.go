package log

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
)

var DefaultLogger = logrus.StandardLogger()

const (
	TraceLevel = logrus.TraceLevel
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

func FromContext(ctx context.Context) *logrus.Entry {
	// Get logger from context
	return ctx.Value("log").(*logrus.Entry)
}

func Setup(ctx context.Context) (context.Context, *logrus.Entry) {
	DefaultLogger.SetOutput(os.Stderr)

	// Check if stdin is a terminal
	if isatty.IsTerminal(os.Stderr.Fd()) {
		// Log text to terminal
		DefaultLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		// Log JSON everywhere else
		DefaultLogger.SetFormatter(&logrus.JSONFormatter{})
	}

	entry := logrus.NewEntry(DefaultLogger)
	return context.WithValue(ctx, "log", entry), entry
}
