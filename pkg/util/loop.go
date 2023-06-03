package util

import (
	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
)

// LoopDevice represents a loop device
type LoopDevice struct {
	DevicePath string
}

// Detach detaches the loop device
func (l *LoopDevice) Detach() {
	// Detach loop device
	if err := system.Local.RunCommand("losetup", "-d", l.DevicePath); err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}
