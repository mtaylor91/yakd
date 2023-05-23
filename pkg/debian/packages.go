package debian

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
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
func installBasePackages(target string) error {
	// Install packages
	if err := installPackages(target, basePackages...); err != nil {
		return err
	}

	return nil
}

// installKubePackages installs the Kubernetes packages
func installKubePackages(target string) error {
	// Install packages
	if err := installPackages(target, kubePackages...); err != nil {
		return err
	}

	// Hold packages
	if err := holdPackages(target, kubePackages...); err != nil {
		return err
	}

	return nil
}

// holdPackages is a helper function to hold packages at a specific version
func holdPackages(target string, packages ...string) error {
	// Look for chroot
	chroot, err := exec.LookPath("chroot")
	if err != nil {
		return err
	}

	// Hold packages
	log.Infof("Holding packages %v", packages)
	args := []string{target, "apt-mark", "hold"}
	args = append(args, packages...)
	cmd := exec.Command(chroot, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// installPackages is a helper function to install packages
func installPackages(target string, packages ...string) error {
	// Update apt indices
	log.Infof("Updating apt indices")
	args := []string{target, "apt-get", "update"}
	if err := util.RunCmd("chroot", args...); err != nil {
		return err
	}

	// Install packages
	log.Infof("Installing packages %v", packages)
	args = []string{target, "apt-get", "install", "-y"}
	args = append(args, packages...)
	if err := util.RunCmd("chroot", args...); err != nil {
		return err
	}

	return nil
}
