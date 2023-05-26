package executor

import (
	"context"
	"io"
)

var Default = &LocalExecutor{}

// GetOutput runs a command and returns the output.
func GetOutput(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	return Default.GetOutput(ctx, cmd, args...)
}

// GetOutputWithStdin runs a command and returns the output.
func GetOutputWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) ([]byte, error) {
	return Default.GetOutputWithStdin(ctx, cmd, stdin, args...)
}

// RunCmd runs a command with output redirected to stderr.
func RunCmd(ctx context.Context, cmd string, args ...string) error {
	return Default.RunCmd(ctx, cmd, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func RunCmdWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) error {
	return Default.RunCmdWithStdin(ctx, cmd, stdin, args...)
}
