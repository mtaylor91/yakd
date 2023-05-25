package debian

import "github.com/mtaylor91/yakd/pkg/util/executor"

// configureAdminUser creates the admin user.
func configureAdminUser(exec executor.Executor) error {
	if err := exec.RunCmd("useradd", "-m", "-G", "sudo", "admin"); err != nil {
		return err
	}

	return nil
}
