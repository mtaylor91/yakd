package util

import (
	"io"
	"os"
	"os/exec"
)

// GetOutput runs a command and returns the output.
func GetOutput(cmd string, args ...string) ([]byte, error) {
	return GetOutputWithStdin(cmd, nil, args...)
}

// GetOutputWithStdin runs a command and returns the output.
func GetOutputWithStdin(cmd string, stdin io.Reader, args ...string) ([]byte, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	return c.Output()
}

// RunCmd runs a command with output redirected to stderr.
func RunCmd(cmd string, args ...string) error {
	return RunCmdWithStdin(cmd, nil, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func RunCmdWithStdin(cmd string, stdin io.Reader, args ...string) error {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	return c.Run()
}
