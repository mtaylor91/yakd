package util

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

// LoopDevice represents a loop device
type LoopDevice struct {
	DevicePath string
}

// Detach detaches the loop device
func (l *LoopDevice) Detach() {
	// Detach loop device
	ctx := context.Background()
	if err := executor.RunCmd(ctx, "losetup", "-d", l.DevicePath); err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}
