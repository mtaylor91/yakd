package util

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RawImage represents a disk image
type RawImage struct {
	ImagePath string
	cleanup   bool
	sizeMB    int
	overwrite bool
}

// NewRawImage initializes a new RawImage struct
func NewRawImage(path string, sizeMB int, cleanup, overwrite bool) *RawImage {
	return &RawImage{
		ImagePath: path,
		cleanup:   true,
		sizeMB:    sizeMB,
		overwrite: overwrite,
	}
}

// Alloc allocates a new image file
func (i *RawImage) Alloc() error {
	// Create image using dd
	if err := RunCmd("dd", "if=/dev/zero", "of="+i.ImagePath, "bs=1M",
		"count=1", "seek="+strconv.Itoa(i.sizeMB-1)); err != nil {
		return err
	}

	return nil
}

// Attach attaches the image to a loop device
func (i *RawImage) Attach() (*LoopDevice, error) {
	// Attach image to loop device
	log.Infof("Attaching image %s to loop device", i.ImagePath)
	if err := RunCmd("losetup", "-P", "-f", i.ImagePath); err != nil {
		return nil, err
	}

	// Get loop device info
	log.Infof("Getting loop device info for %s", i.ImagePath)
	if out, err := GetOutput("losetup", "-j", i.ImagePath); err != nil {
		return nil, err
	} else {
		// Get loop device path
		loopPath := strings.Split(string(out), ":")[0]
		log.Infof("Loop device path is %s", loopPath)

		return &LoopDevice{loopPath}, nil
	}
}

// Convert converts the image
func (i *RawImage) Convert(output string) error {
	// Convert image to qcow2
	log.Infof("Converting image %s to qcow2", i.ImagePath)
	format := filepath.Ext(output)[1:]
	if err := RunCmd("qemu-img", "convert", "-f", "raw",
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

// Partition partitions the image
func (i *RawImage) Partition() error {
	// Create partition table
	log.Infof("Creating partition table on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath, "mklabel", "gpt"); err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath,
		"mkpart", "primary", "fat32", "1MiB", "512MiB"); err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath,
		"mkpart", "primary", "ext4", "512MiB", "100%"); err != nil {
		return err
	}

	// Set boot flag on EFI partition
	log.Infof("Setting boot flag on EFI partition on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath,
		"set", "1", "boot", "on"); err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath,
		"set", "1", "esp", "on"); err != nil {
		return err
	}

	// Set root partition label
	log.Infof("Setting root partition label on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath, "name", "2", "root"); err != nil {
		return err
	}

	// Set EFI partition label
	log.Infof("Setting EFI partition label on %s", i.ImagePath)
	if err := RunCmd("parted", i.ImagePath, "name", "1", "efi"); err != nil {
		return err
	}

	return nil
}
