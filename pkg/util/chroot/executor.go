package chroot

import (
	"context"
	"io"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// GetOutput implements Executor.GetOutput.
func (c *ChrootExecutor) GetOutput(
	ctx context.Context, cmd string, args ...string,
) ([]byte, error) {
	return c.GetOutputWithStdin(ctx, cmd, nil, args...)
}

// GetOutputWithStdin implements Executor.GetOutputWithStdin.
func (c *ChrootExecutor) GetOutputWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) ([]byte, error) {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if !c.isSetup {
		return nil, ErrNotSetup
	}

	return executor.Default.GetOutputWithStdin(
		ctx, "chroot", stdin, append([]string{c.root, cmd}, args...)...,
	)
}

// RunCmd implements Executor.RunCmd.
func (c *ChrootExecutor) RunCmd(ctx context.Context, cmd string, args ...string) error {
	return c.RunCmdWithStdin(ctx, cmd, nil, args...)
}

// RunCmdWithStdin implements Executor.RunCmdWithStdin.
func (c *ChrootExecutor) RunCmdWithStdin(
	ctx context.Context, cmd string, stdin io.Reader, args ...string,
) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if !c.isSetup {
		return ErrNotSetup
	}

	return executor.Default.RunCmdWithStdin(
		ctx, "chroot", stdin, append([]string{c.root, cmd}, args...)...,
	)
}
