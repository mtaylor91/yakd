package system

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

type System interface {
	// Context returns the context
	Context() context.Context
	// Logger returns the logger
	Logger() *logrus.Logger
	// RunCommand runs the given command
	RunCommand(cmd string, args ...string) error
	// RunCommandWithInput runs the given command with the given input
	RunCommandWithInput(input io.Reader, cmd string, args ...string) error
	// RunCommandWithOutput runs the given command and returns the output
	RunCommandWithOutput(output io.Writer, cmd string, args ...string) error
	// RunCommandWithIO runs the given command with the given input/output streams
	RunCommandWithIO(
		input io.Reader, output io.Writer, cmd string, args ...string) error
	// SetContext sets the context
	SetContext(ctx context.Context)
	// SetLogger sets the logger
	SetLogger(logger *logrus.Logger)
}
