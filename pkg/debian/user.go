package debian

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// configureAdminUser creates the admin user.
func configureAdminUser(ctx context.Context, exec executor.Executor) error {
	if err := exec.RunCmd(ctx, "useradd", "-m", "-G", "sudo", "admin"); err != nil {
		return err
	}

	return nil
}
