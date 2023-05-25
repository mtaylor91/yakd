package bootstrap

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/chroot"
)

// TmpFS is a filesystem that is mounted as a tmpfs
type TmpFS struct {
	Path   string
	SizeMB int
}

// NewTmpFS initializes a new TmpFS struct
func NewTmpFS(path string, sizeMB int) *TmpFS {
	return &TmpFS{
		Path:   path,
		SizeMB: sizeMB,
	}
}

// Allocate creates the tmpfs
func (t *TmpFS) Allocate() error {
	if err := util.CreateMountpointAt(t.Path); err != nil {
		return err
	}

	return MountTmpFSAt(t.Path, t.SizeMB)
}

// Bootstrap runs filesystem bootstrapping
func (t *TmpFS) Bootstrap(operatingSystem os.OS) error {
	// Create mountpoint
	if err := t.Allocate(); err != nil {
		return err
	}

	defer t.Destroy()

	// Bootstrap OS
	installer := operatingSystem.Installer(t.Path)
	err := installer.Bootstrap()
	if err != nil {
		return err
	}

	// Setup chroot executor
	log.Infof("Setting up chroot at %s", t.Path)
	chrootExecutor := chroot.NewExecutor(t.Path)
	defer chrootExecutor.Teardown()

	// Run post-bootstrap step
	log.Infof("Running post-bootstrap step")
	if err := installer.PostBootstrap(chrootExecutor); err != nil {
		return err
	}

	return nil
}

// Destroy removes the tmpfs
func (t *TmpFS) Destroy() {
	util.UnmountFilesystems(t.Path)
	util.RemoveMountpointAt(t.Path)
}

// MountTmpFSAt mounts a tmpfs at the given path
func MountTmpFSAt(path string, sizeMB int) error {
	options := fmt.Sprintf("size=%dM", sizeMB)
	err := util.RunCmd("mount", "-t", "tmpfs", "-o", options, "tmpfs", path)
	if err != nil {
		return err
	}

	return nil
}
