package bootstrap

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Stage2 represents the second stage of the bootstrap process
type Stage2 struct {
	// contains filtered or unexported fields
	stage1     string
	mountpoint string
	imagePath  string
	imageSize  int
}

// NewStage2 initializes a new Stage2 struct
func NewStage2(stage1, mountpoint, imagePath string, imageSize int) *Stage2 {
	return &Stage2{
		stage1:     stage1,
		mountpoint: mountpoint,
		imagePath:  imagePath,
		imageSize:  imageSize,
	}
}

// Run runs the second stage of the bootstrap process
func (s *Stage2) Run() error {
	// Check if image exists
	if _, err := os.Stat(s.imagePath); err == nil {
		return fmt.Errorf("image already exists")
	}

	// Initialize image
	i := NewImage(s.imagePath, s.imageSize, true, true)

	// Create image
	log.Infof("Creating image %s", i.imagePath)
	if err := i.Alloc(); err != nil {
		return err
	}

	// Create partition table
	log.Infof("Creating partition table on %s", i.imagePath)
	if err := i.Partition(); err != nil {
		return err
	}

	// Attach image
	log.Infof("Attaching image %s", i.imagePath)
	loop, err := i.Attach()
	if err != nil {
		return err
	}

	defer loop.Detach()

	// Format image
	log.Infof("Formatting image %s on %s", i.imagePath, loop.devicePath)
	if err := loop.Format(); err != nil {
		return err
	}

	// Identify partitions
	esp := loop.devicePath + "p1"
	root := loop.devicePath + "p2"

	// Initialize image disk
	d := NewDisk(loop.devicePath, esp, root, s.mountpoint, true)

	// Bootstrap disk image
	log.Infof("Bootstrapping disk image %s", d.devicePath)
	if err := d.Bootstrap(); err != nil {
		return err
	}

	return nil
}
