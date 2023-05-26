package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type LocalExecutor struct{}

// GetOutput runs a command and returns the output.
func (l *LocalExecutor) GetOutput(
	ctx context.Context, cmd string, args ...string,
) ([]byte, error) {
	return l.GetOutputWithStdin(ctx, cmd, nil, args...)
}

// GetOutputWithStdin runs a command and returns the output.
func (l *LocalExecutor) GetOutputWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) ([]byte, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	log.Debugf("Getting output of: %s %s", cmd, strings.Join(args, " "))
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdin = stdin
	return c.Output()
}

// RunCmd runs a command with output redirected to stderr.
func (l *LocalExecutor) RunCmd(
	ctx context.Context, cmd string, args ...string) error {
	return l.RunCmdWithStdin(ctx, cmd, nil, args...)
}

// RunCmdWithStdin runs a command with output redirected to stderr.
func (l *LocalExecutor) RunCmdWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) error {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	log.Debugf("Running command: %s %s", cmd, strings.Join(args, " "))
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdin = stdin
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	return c.Run()
}
