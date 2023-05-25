package executor

import (
	"io"
	"os"
	"os/exec"
)

var Default = &LocalExecutor{}

type LocalExecutor struct{}

// GetOutput runs a command and returns the output.
func (l *LocalExecutor) GetOutput(cmd string, args ...string) ([]byte, error) {
	return l.GetOutputWithStdin(cmd, nil, args...)
}

// GetOutputWithStdin runs a command and returns the output.
func (l *LocalExecutor) GetOutputWithStdin(
	cmd string, stdin io.Reader, args ...string,
) ([]byte, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	return c.Output()
}

// RunCmd runs a command with output redirected to stderr.
func (l *LocalExecutor) RunCmd(cmd string, args ...string) error {
	return l.RunCmdWithStdin(cmd, nil, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func (l *LocalExecutor) RunCmdWithStdin(
	cmd string, stdin io.Reader, args ...string,
) error {
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
