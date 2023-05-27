package gentoo

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const env02locale = `LANG="en_CA.UTF-8"
LC_COLLATE="C.UTF-8"
`

const gentooRepoConf = `[DEFAULT]
main-repo = gentoo

[gentoo]
location = /var/db/repos/gentoo
sync-type = rsync
sync-uri = rsync://rsync.gentoo.org/gentoo-portage
auto-sync = yes
sync-rsync-verify-jobs = 1
sync-rsync-verify-metamanifest = yes
sync-rsync-verify-max-age = 24
sync-openpgp-key-path = /usr/share/openpgp-keys/gentoo-release.asc
sync-openpgp-key-refresh-retry-count = 40
sync-openpgp-key-refresh-retry-overall-timeout = 1200
sync-openpgp-key-refresh-retry-delay-exp-base = 2
sync-openpgp-key-refresh-retry-delay-max = 60
sync-openpgp-key-refresh-retry-delay-mult = 4
`

const localeGen = `
en_CA.UTF-8 UTF-8
en_US.UTF-8 UTF-8
`

const makeConfTemplate = `# Please consult /usr/share/portage/config/make.conf.example
# for a more detailed example.
COMMON_FLAGS="-O2 -pipe"
CFLAGS="${COMMON_FLAGS}"
CXXFLAGS="${COMMON_FLAGS}"
FCFLAGS="${COMMON_FLAGS}"
FFLAGS="${COMMON_FLAGS}"

# NOTE: This stage was built with the bindist Use flag enabled

# This sets the language of build output to English.
# Please keep this setting intact when reporting bugs.
LC_MESSAGES=C.utf8
MAKEOPTS="-j{{.NumCores}}"
`

type GentooBootstrapInstaller struct {
	stage3 string
	target string
}

func (g *GentooBootstrapInstaller) Bootstrap(ctx context.Context) error {
	if _, err := os.Stat(g.stage3); err == nil {
		log.Infof("Stage3 tarball already exists at %s", g.stage3)
	} else if os.IsNotExist(err) {
		// Download stage3 tarball
		err := DownloadStage3(ctx, g.stage3)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	log.Infof("Unpacking stage3 tarball to %s", g.target)
	err := util.UnpackTarball(ctx, g.stage3, g.target)
	if err != nil {
		return err
	}

	return nil
}

func (g *GentooBootstrapInstaller) PostBootstrap(
	ctx context.Context, chroot executor.Executor,
) error {
	defer debugChroot(ctx, chroot)

	// Render make.conf template
	makeConf, err := util.TemplateString(makeConfTemplate, map[string]interface{}{
		"NumCores": runtime.NumCPU(),
	})

	// Write target system make.conf
	err = util.WriteFile(path.Join(g.target, "etc/portage/make.conf"), makeConf)
	if err != nil {
		return err
	}

	log.Info("Configuring repositories")

	// Create repos.conf directory
	err = os.MkdirAll(path.Join(g.target, "etc/portage/repos.conf"), 0755)
	if err != nil {
		return err
	}

	// Write gentoo repo conf
	gentooRepoConfPath := path.Join(g.target, "etc/portage/repos.conf/gentoo.conf")
	err = util.WriteFile(gentooRepoConfPath, gentooRepoConf)
	if err != nil {
		return err
	}

	// Run emerge-webrsync
	log.Infof("Running emerge-webrsync")
	err = chroot.RunCmd(ctx, "emerge-webrsync")
	if err != nil {
		return err
	}

	// Emerge @world updates
	log.Infof("Emerging @world updates")
	log.Warning("Disabled for testing")
	err = chroot.RunCmd(ctx, "emerge", "--update", "--deep", "--newuse", "@world")
	if err != nil {
		return err
	}

	// Remove timezone symlink
	log.Infof("Removing /etc/localtime symlink")
	err = os.Remove(path.Join(g.target, "etc/localtime"))
	if err != nil {
		return err
	}

	// Link /etc/localtime to /usr/share/zoneinfo/UTC
	log.Infof("Linking /etc/localtime to /usr/share/zoneinfo/UTC")
	err = os.Symlink("/usr/share/zoneinfo/UTC", path.Join(g.target, "etc/localtime"))
	if err != nil {
		return err
	}

	// Write locale.gen
	log.Infof("Writing locale.gen")
	localeGenPath := path.Join(g.target, "etc", "locale.gen")
	if err := os.WriteFile(localeGenPath, []byte(localeGen), 0644); err != nil {
		return err
	}

	// Configure locales
	log.Infof("Configuring locales")
	if err := chroot.RunCmd(ctx, "locale-gen"); err != nil {
		return err
	}

	// Install gentoo-kernel
	log.Infof("Installing gentoo-kernel")
	log.Warning("Disabled for testing")
	if err := chroot.RunCmd(ctx, "emerge", "sys-kernel/gentoo-kernel"); err != nil {
		return err
	}

	// Unmask cri-o
	log.Infof("Unmasking app-containers/cri-o")
	err = acceptKeywords(g.target, "app-containers", "cri-o", 99, "~amd64")
	if err != nil {
		return err
	}

	// Install cri-o
	log.Infof("Installing kubernetes packages")
	if err := installPackages(ctx, chroot,
		"app-containers/cri-o",
		"sys-cluster/kubeadm",
		"sys-cluster/kubectl",
		"sys-cluster/kubelet",
	); err != nil {
		return err
	}

	return nil
}

func acceptKeywords(
	target, section, pkg string, priority int, keywords ...string,
) error {
	filename := fmt.Sprintf("%02d-%s", priority, pkg)
	return util.AppendFile(
		path.Join(target, "etc", "portage", "package.accept_keywords", filename),
		strings.Join(append(
			[]string{fmt.Sprintf("%s/%s", section, pkg)}, keywords...,
		), " "),
	)
}

func installPackages(
	ctx context.Context, chroot executor.Executor, pkgs ...string,
) error {
	return chroot.RunCmd(ctx, "emerge", strings.Join(pkgs, " "))
}

func debugChroot(ctx context.Context, chroot executor.Executor) error {
	return chroot.RunCmdWithStdin(ctx, "/bin/bash", os.Stdin)
}