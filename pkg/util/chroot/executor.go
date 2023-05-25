package chroot

import (
	"io"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// GetOutput implements Executor.GetOutput.
func (c *ChrootExecutor) GetOutput(cmd string, args ...string) ([]byte, error) {
	return c.GetOutputWithStdin(cmd, nil, args...)
}

// GetOutputWithStdin implements Executor.GetOutputWithStdin.
func (c *ChrootExecutor) GetOutputWithStdin(
	cmd string, stdin io.Reader, args ...string,
) ([]byte, error) {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if !c.isSetup {
		return nil, ErrNotSetup
	}

	return executor.Default.GetOutputWithStdin(
		"chroot", stdin, append([]string{c.root, cmd}, args...)...,
	)
}

// RunCmd implements Executor.RunCmd.
func (c *ChrootExecutor) RunCmd(cmd string, args ...string) error {
	return c.RunCmdWithStdin(cmd, nil, args...)
}

// RunCmdWithStdin implements Executor.RunCmdWithStdin.
func (c *ChrootExecutor) RunCmdWithStdin(
	cmd string, stdin io.Reader, args ...string,
) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if !c.isSetup {
		return ErrNotSetup
	}

	return executor.Default.RunCmdWithStdin(
		"chroot", stdin, append([]string{c.root, cmd}, args...)...,
	)
}
