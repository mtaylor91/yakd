package disk

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// BuildDisk builds a disk from a stage1 tarball
func BuildDisk(target, stage1, mountpoint string) error {
	log.Infof("Building disk from %s", stage1)
	return fmt.Errorf("Not implemented")
}
