package debian

import (
	"github.com/mtaylor91/yakd/pkg/system"
)

// configureAdminUser creates the admin user.
func configureAdminUser(sys system.System) error {
	sys.Logger().Infof("Configuring admin user")

	if err := sys.RunCommand("useradd", "-m", "-G", "sudo", "admin"); err != nil {
		return err
	}

	if err := sys.RunCommand("passwd", "-d", "admin"); err != nil {
		return err
	}

	return nil
}
