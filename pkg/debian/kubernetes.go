package debian

import (
	"path"

	"github.com/mtaylor91/yakd/pkg/util"
)

const sysctlConf = `net.ipv4.ip_forward=1
net.bridge.bridge-nf-call-iptables=1
`

// configureKubernetes configures the target system to run Kubernetes.
func configureKubernetes(target string) error {
	modulesLoad := path.Join(target, "etc", "modules-load.d", "10-kubernetes.conf")
	if err := util.WriteFile(modulesLoad, "br_netfilter\n"); err != nil {
		return err
	}

	sysctl := path.Join(target, "etc", "sysctl.d", "10-kubernetes.conf")
	if err := util.WriteFile(sysctl, sysctlConf); err != nil {
		return err
	}

	return nil
}
