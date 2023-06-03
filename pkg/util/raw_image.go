package util

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
)

// RawImage represents a disk image
type RawImage struct {
	ImagePath string
	sizeMB    int
	overwrite bool
}

// NewRawImage initializes a new RawImage struct
func NewRawImage(path string, sizeMB int, overwrite bool) *RawImage {
	return &RawImage{
		ImagePath: path,
		sizeMB:    sizeMB,
		overwrite: overwrite,
	}
}

// Alloc allocates a new image file
func (i *RawImage) Alloc(ctx context.Context) error {
	// Create image using dd
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("dd", "if=/dev/zero", "of="+i.ImagePath, "bs=1M",
		"count=1", "seek="+strconv.Itoa(i.sizeMB-1)); err != nil {
		return err
	}

	return nil
}

// Attach attaches the image to a loop device
func (i *RawImage) Attach(ctx context.Context) (*LoopDevice, error) {
	// Attach image to loop device
	log.Infof("Attaching image %s to loop device", i.ImagePath)
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("losetup", "-P", "-f", i.ImagePath); err != nil {
		return nil, err
	}

	// Get loop device info
	log.Infof("Getting loop device info for %s", i.ImagePath)
	var losetupOutput bytes.Buffer
	if err := sys.RunCommandWithOutput(
		&losetupOutput, "losetup", "-j", i.ImagePath); err != nil {
		return nil, err
	}

	// Get loop device path
	loopPath := strings.Split(losetupOutput.String(), ":")[0]
	log.Infof("Loop device path is %s", loopPath)

	return &LoopDevice{loopPath}, nil
}

// Convert converts the image
func (i *RawImage) Convert(ctx context.Context, output string) error {
	// Convert image to qcow2
	log.Infof("Converting image %s to qcow2", i.ImagePath)
	format := filepath.Ext(output)[1:]
	sys := system.Local.WithContext(ctx)
	if err := sys.RunCommand("qemu-img", "convert", "-f", "raw",
		"-O", format, i.ImagePath, output); err != nil {
		return err
	}

	return nil
}

// Free removes the raw image
func (i *RawImage) Free() {
	err := os.Remove(i.ImagePath)
	if err != nil {
		log.Warnf("Failed to remove %s image: %s", i.ImagePath, err)
	}
}
