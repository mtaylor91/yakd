package bootstrap

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mtaylor91/yakd/pkg/util"
	log "github.com/sirupsen/logrus"
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
func (t *TmpFS) Bootstrap(osFactory OSFactory) error {
	// Create mountpoint
	if err := t.Allocate(); err != nil {
		return err
	}

	// Bootstrap OS
	os := osFactory.NewOS(t.Path)
	err := os.Bootstrap()
	if err != nil {
		return err
	}

	// Mount metadata filesystems
	log.Infof("Mounting metadata filesystems on %s", t.Path)
	if err := util.MountMetadataFilesystems(t.Path); err != nil {
		return err
	}

	defer util.UnmountMetadataFilesystems(t.Path)

	// Run post-bootstrap step
	log.Infof("Running post-bootstrap step")
	if err := os.PostBootstrap(); err != nil {
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
	mount, err := exec.LookPath("mount")
	if err != nil {
		return err
	}

	options := fmt.Sprintf("size=%dM", sizeMB)
	cmd := exec.Command(mount, "-t", "tmpfs", "-o", options, "tmpfs", path)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
