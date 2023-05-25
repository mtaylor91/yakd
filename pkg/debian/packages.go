package debian

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

var basePackages = []string{
	"apt-transport-https",
	"ca-certificates",
	"curl",
	"gnupg2",
	"lvm2",
	"sudo",
}

var kubePackages = []string{
	"cri-o",
	"cri-o-runc",
	"kubeadm",
	"kubectl",
	"kubelet",
}

// installBasePackages installs the base packages
func installBasePackages(exec executor.Executor) error {
	// Install packages
	if err := installPackages(exec, basePackages...); err != nil {
		return err
	}

	return nil
}

// installKubePackages installs the Kubernetes packages
func installKubePackages(exec executor.Executor) error {
	// Install packages
	if err := installPackages(exec, kubePackages...); err != nil {
		return err
	}

	// Hold packages
	if err := holdPackages(exec, kubePackages...); err != nil {
		return err
	}

	return nil
}

// holdPackages is a helper function to hold packages at a specific version
func holdPackages(exec executor.Executor, packages ...string) error {
	// Hold packages
	log.Infof("Holding packages %v", packages)
	args := append([]string{"hold"}, packages...)
	if err := exec.RunCmd("apt-mark", args...); err != nil {
		return err
	}

	return nil
}

// installPackages is a helper function to install packages
func installPackages(exec executor.Executor, packages ...string) error {
	// Update apt indices
	log.Infof("Updating apt indices")
	if err := exec.RunCmd("apt-get", "update"); err != nil {
		return err
	}

	// Install packages
	log.Infof("Installing packages %v", packages)
	args := append([]string{"install", "-y"}, packages...)
	if err := exec.RunCmd("apt-get", args...); err != nil {
		return err
	}

	return nil
}
