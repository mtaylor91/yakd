package executor

import "io"

type Executor interface {
	GetOutput(cmd string, args ...string) ([]byte, error)
	GetOutputWithStdin(cmd string, stdin io.Reader, args ...string) ([]byte, error)
	RunCmd(cmd string, args ...string) error
	RunCmdWithStdin(cmd string, stdin io.Reader, args ...string) error
}
