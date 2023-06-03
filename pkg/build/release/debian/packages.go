package debian

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
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
func installBasePackages(sys system.System) error {
	// Install packages
	if err := installPackages(sys, basePackages...); err != nil {
		return err
	}

	return nil
}

// installKubePackages installs the Kubernetes packages
func installKubePackages(sys system.System) error {
	// Install packages
	if err := installPackages(sys, kubePackages...); err != nil {
		return err
	}

	// Hold packages
	if err := holdPackages(sys, kubePackages...); err != nil {
		return err
	}

	return nil
}

// holdPackages is a helper function to hold packages at a specific version
func holdPackages(sys system.System, packages ...string) error {
	// Hold packages
	log.Infof("Holding packages %v", packages)
	args := append([]string{"hold"}, packages...)
	if err := sys.RunCommand("apt-mark", args...); err != nil {
		return err
	}

	return nil
}

// installPackages is a helper function to install packages
func installPackages(sys system.System, packages ...string) error {
	// Update apt indices
	log.Infof("Updating apt indices")
	if err := sys.RunCommand("apt-get", "update"); err != nil {
		return err
	}

	// Install packages
	log.Infof("Installing packages %v", packages)
	args := append([]string{"install", "-y"}, packages...)
	if err := sys.RunCommand("apt-get", args...); err != nil {
		return err
	}

	return nil
}
