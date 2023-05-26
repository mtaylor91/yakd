package chroot

import (
	"context"
	"fmt"
	"sync"
)

type ChrootExecutor struct {
	// contains filtered or unexported fields
	isSetup  bool
	root     string
	runMutex sync.Mutex
}

// NewExecutor returns a new ChrootExecutor.
func NewExecutor(ctx context.Context, root string) *ChrootExecutor {
	chroot := &ChrootExecutor{false, root, sync.Mutex{}}
	chroot.Setup(ctx)
	return chroot
}

// Setup sets up the chroot.
func (c *ChrootExecutor) Setup(ctx context.Context) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if c.isSetup {
		return nil
	}

	if c.root == "" {
		return ErrNoRoot
	}

	if err := MountMetadataFilesystems(ctx, c.root); err != nil {
		return fmt.Errorf("chroot failed: %s", err)
	}

	c.isSetup = true
	return nil
}

// Teardown tears down the chroot.
func (c *ChrootExecutor) Teardown() {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if !c.isSetup {
		return
	}

	ctx := context.Background()
	UnmountMetadataFilesystems(ctx, c.root)
	c.isSetup = false
}
