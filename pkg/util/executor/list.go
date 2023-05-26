package executor

import "context"

func RunCmdList(ctx context.Context, executor Executor, cmds ...[]string) error {
	for _, cmd := range cmds {
		if err := executor.RunCmd(ctx, cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}
