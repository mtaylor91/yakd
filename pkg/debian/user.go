package debian

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// configureAdminUser creates the admin user.
func configureAdminUser(ctx context.Context, exec executor.Executor) error {
	log.Infof("Configuring admin user")

	if err := exec.RunCmd(ctx, "useradd", "-m", "-G", "sudo", "admin"); err != nil {
		return err
	}

	if err := exec.RunCmd(ctx, "passwd", "-d", "admin"); err != nil {
		return err
	}

	return nil
}
