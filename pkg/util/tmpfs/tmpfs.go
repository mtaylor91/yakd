package tmpfs

import (
	"context"
	"fmt"

	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
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
func (t *TmpFS) Allocate(ctx context.Context) error {
	if err := util.CreateMountpointAt(ctx, t.Path); err != nil {
		return err
	}

	return MountTmpFSAt(ctx, t.Path, t.SizeMB)
}

// Destroy removes the tmpfs
func (t *TmpFS) Destroy() {
	util.UnmountRecursive(t.Path)
	util.RemoveMountpointAt(t.Path)
}

// MountTmpFSAt mounts a tmpfs at the given path
func MountTmpFSAt(ctx context.Context, path string, sizeMB int) error {
	options := fmt.Sprintf("size=%dM", sizeMB)
	sys := system.Local.WithContext(ctx)
	err := sys.RunCommand("mount", "-t", "tmpfs", "-o", options, "tmpfs", path)
	if err != nil {
		return err
	}

	return nil
}
