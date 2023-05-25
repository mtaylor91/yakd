package executor

func RunCmdList(executor Executor, cmds ...[]string) error {
	for _, cmd := range cmds {
		if err := executor.RunCmd(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}
