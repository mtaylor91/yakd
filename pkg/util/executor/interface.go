package executor

import (
	"context"
	"io"
)

type Executor interface {
	GetOutput(ctx context.Context, cmd string, args ...string) ([]byte, error)
	GetOutputWithStdin(
		ctx context.Context, cmd string, stdin io.Reader, args ...string,
	) ([]byte, error)
	RunCmd(ctx context.Context, cmd string, args ...string) error
	RunCmdWithStdin(
		ctx context.Context, cmd string, stdin io.Reader, args ...string,
	) error
}
