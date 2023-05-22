package util

import (
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command with output redirected to stderr.
func RunCmd(cmd string, args ...string) error {
	return RunCmdWithStdin(cmd, nil, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func RunCmdWithStdin(cmd string, stdin io.Reader, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	return c.Run()
}
