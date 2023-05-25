package util

import (
	"io"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// GetOutput runs a command and returns the output.
func GetOutput(cmd string, args ...string) ([]byte, error) {
	return executor.Default.GetOutput(cmd, args...)
}

// GetOutputWithStdin runs a command and returns the output.
func GetOutputWithStdin(cmd string, stdin io.Reader, args ...string) ([]byte, error) {
	return executor.Default.GetOutputWithStdin(cmd, stdin, args...)
}

// RunCmd runs a command with output redirected to stderr.
func RunCmd(cmd string, args ...string) error {
	return executor.Default.RunCmd(cmd, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func RunCmdWithStdin(cmd string, stdin io.Reader, args ...string) error {
	return executor.Default.RunCmdWithStdin(cmd, stdin, args...)
}
