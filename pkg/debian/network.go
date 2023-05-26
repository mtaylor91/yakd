package debian

import (
	"context"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const ethernetConfig = `[Match]
Name=enp*s*

[Network]
DHCP=yes
`

const defaultHostname = "yakd"

const hostsConfig = `127.0.0.1 localhost
::1 localhost ip6-localhost ip6-loopback
`

// configureNetworking configures networking for the target system.
func configureNetworking(
	ctx context.Context, exec executor.Executor, target string,
) error {
	log.Infof("Configuring networking")

	ethernet := path.Join(target, "etc", "systemd", "network", "10-ethernet.network")
	if err := util.WriteFile(ethernet, ethernetConfig); err != nil {
		return err
	}

	hostname := path.Join(target, "etc", "hostname")
	if err := util.WriteFile(hostname, defaultHostname); err != nil {
		return err
	}

	hosts := path.Join(target, "etc", "hosts")
	if err := util.WriteFile(hosts, hostsConfig); err != nil {
		return err
	}

	err := exec.RunCmd(ctx, "systemctl", "enable", "systemd-networkd")
	if err != nil {
		return err
	}

	return nil
}
