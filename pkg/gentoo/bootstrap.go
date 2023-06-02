package gentoo

import (
	"context"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/os/common"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const containersPolicyJSON = `{
	"default": [
		{
			"type": "insecureAcceptAnything"
		}
	],
	"transports": {
		"docker-daemon": {
			"": [{"type":"insecureAcceptAnything"}]
		}
	}
}
`

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

const gentooYAKDRepoConf = `[yakd-gentoo]
location = /var/db/repos/yakd
sync-type = git
sync-uri = https://github.com/mtaylor91/yakd-gentoo.git
auto-sync = yes
priority = 100
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
LC_MESSAGES=C.utf8
MAKEOPTS="-j{{.NumCores}}"
BINPKG_FORMAT="gpkg"
FEATURES="buildpkg"
`

const sudoersWheel = `%wheel ALL=(ALL:ALL) ALL
`

type GentooBootstrapInstaller struct {
	binPkgsCache string
	stage3       string
	target       string
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
	// Ensure binPkgsCache exists
	err := os.MkdirAll(g.binPkgsCache, 0755)
	if err != nil {
		return err
	}

	// Bind binPkgsCache to /var/cache/binpkgs
	if err = executor.Default.RunCmd(
		ctx, "mount", "--bind",
		g.binPkgsCache,
		path.Join(g.target, "var/cache/binpkgs"),
	); err != nil {
		return err
	}

	// Unmount /var/cache/binpkgs on exit
	defer func() {
		if err := executor.Default.RunCmd(
			ctx, "umount", path.Join(g.target, "var/cache/binpkgs"),
		); err != nil {
			log.Warnf("Failed to unmount /var/cache/binpkgs: %s", err)
		}
	}()

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

	// Write YAKD repo conf
	yakdRepoConfPath := path.Join(
		g.target, "etc/portage/repos.conf/yakd-gentoo.conf")
	err = util.WriteFile(yakdRepoConfPath, gentooYAKDRepoConf)
	if err != nil {
		return err
	}

	// Run emerge-webrsync
	log.Infof("Running emerge-webrsync")
	err = chroot.RunCmd(ctx, "emerge-webrsync")
	if err != nil {
		return err
	}

	// Install dev-vcs/git
	log.Infof("Installing dev-vcs/git")
	err = installPackages(ctx, chroot, "dev-vcs/git")
	if err != nil {
		return err
	}

	// Run emerge --sync
	log.Infof("Running emerge --sync")
	err = chroot.RunCmd(ctx, "emerge", "--sync", "yakd-gentoo")
	if err != nil {
		return err
	}

	// Emerge @world updates
	log.Infof("Emerging @world updates")
	err = chroot.RunCmd(
		ctx, "emerge", "--usepkg", "--update", "--deep", "--newuse", "@world")
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

	// Unmask cri-o
	log.Infof("Unmasking app-containers/cri-o")
	err = acceptKeywords(g.target, "app-containers", "cri-o", 99, "~amd64")
	if err != nil {
		return err
	}

	// Unmask kubeadm
	log.Infof("Unmasking sys-cluster/kubeadm")
	err = acceptKeywords(g.target, "sys-cluster", "kubeadm", 99, "~amd64")
	if err != nil {
		return err
	}

	// Unmask kubectl
	log.Infof("Unmasking sys-cluster/kubectl")
	err = acceptKeywords(g.target, "sys-cluster", "kubectl", 99, "~amd64")
	if err != nil {
		return err
	}

	// Unmask kubelet
	log.Infof("Unmasking sys-cluster/kubelet")
	err = acceptKeywords(g.target, "sys-cluster", "kubelet", 99, "~amd64")
	if err != nil {
		return err
	}

	// Install cri-o
	log.Infof("Installing kubernetes packages")
	if err := installPackages(ctx, chroot,
		"app-admin/sudo",
		"app-containers/cri-o",
		"net-firewall/ebtables",
		"sys-apps/ethtool",
		"sys-cluster/kubeadm",
		"sys-cluster/kubectl",
		"sys-cluster/kubelet",
	); err != nil {
		return err
	}

	// Install /etc/containers/policy.json
	log.Infof("Installing /etc/containers/policy.json")
	if err := util.WriteFile(
		path.Join(g.target, "etc/containers/policy.json"),
		containersPolicyJSON,
	); err != nil {
		return err
	}

	// Install gentoo-kernel
	log.Infof("Installing gentoo-kernel")
	if err := chroot.RunCmd(
		ctx, "emerge", "--usepkg", "sys-kernel/gentoo-kernel"); err != nil {
		return err
	}

	log.Infof("Creating admin user")
	err = chroot.RunCmd(ctx, "useradd", "-m", "-G", "wheel", "admin")
	if err != nil {
		return err
	}

	log.Infof("Removing admin password")
	if err := chroot.RunCmd(ctx, "passwd", "-d", "admin"); err != nil {
		return err
	}

	log.Infof("Configuring sudo")
	// Ensure sudoers.d exists
	err = os.MkdirAll(path.Join(g.target, "etc/sudoers.d"), 0755)
	if err != nil {
		return err
	}
	// Write sudoers.d/wheel
	err = os.WriteFile(
		path.Join(g.target, "etc/sudoers.d/wheel"), []byte(sudoersWheel), 0440)
	if err != nil {
		return err
	}

	log.Infof("Setting machine id")
	// Set machine id
	err = chroot.RunCmd(ctx, "systemd-machine-id-setup")
	if err != nil {
		return err
	}

	if err := common.ConfigureKubernetes(ctx, chroot, g.target); err != nil {
		return err
	}

	if err := common.ConfigureNetwork(ctx, chroot, g.target); err != nil {
		return err
	}

	return nil
}
