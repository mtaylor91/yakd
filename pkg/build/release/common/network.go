package common

import (
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
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

// ConfigureNetwork configures the network for the target system.
func ConfigureNetwork(sys system.System, target string) error {
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

	err := sys.RunCommand("systemctl", "enable", "systemd-networkd")
	if err != nil {
		return err
	}

	return nil
}
