package common

import (
	"context"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const sysctlConf = `net.ipv4.ip_forward=1
net.bridge.bridge-nf-call-iptables=1
`

// ConfigureKubernetes configures the target system to run Kubernetes.
func ConfigureKubernetes(
	ctx context.Context, exec executor.Executor, target string,
) error {
	log.Infof("Configuring system to run Kubernetes")

	modulesLoad := path.Join(target, "etc", "modules-load.d", "10-kubernetes.conf")
	if err := util.WriteFile(modulesLoad, "br_netfilter\n"); err != nil {
		return err
	}

	sysctl := path.Join(target, "etc", "sysctl.d", "10-kubernetes.conf")
	if err := util.WriteFile(sysctl, sysctlConf); err != nil {
		return err
	}

	err := exec.RunCmd(ctx, "systemctl", "enable", "crio")
	if err != nil {
		return err
	}

	err = exec.RunCmd(ctx, "systemctl", "enable", "kubelet")
	if err != nil {
		return err
	}

	return nil
}
