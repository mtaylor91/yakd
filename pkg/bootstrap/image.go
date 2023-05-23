package bootstrap

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
)

// RawImage represents a disk image
type RawImage struct {
	cleanup   bool
	imagePath string
	sizeMB    int
	overwrite bool
}

// Loop represents a loop device
type Loop struct {
	devicePath string
}

// NewImage initializes a new RawImage struct
func NewImage(path string, sizeMB int, cleanup, overwrite bool) *RawImage {
	return &RawImage{
		cleanup:   true,
		imagePath: path,
		sizeMB:    sizeMB,
		overwrite: overwrite,
	}
}

// Alloc allocates a new image file
func (i *RawImage) Alloc() error {
	// Create image using dd
	if err := util.RunCmd("dd", "if=/dev/zero", "of="+i.imagePath, "bs=1M",
		"count=1", "seek="+strconv.Itoa(i.sizeMB-1)); err != nil {
		return err
	}

	return nil
}

// Attach attaches the image to a loop device
func (i *RawImage) Attach() (*Loop, error) {
	// Attach image to loop device
	log.Infof("Attaching image %s to loop device", i.imagePath)
	if err := util.RunCmd("losetup", "-P", "-f", i.imagePath); err != nil {
		return nil, err
	}

	// Get loop device info
	log.Infof("Getting loop device info for %s", i.imagePath)
	if out, err := util.GetOutput("losetup", "-j", i.imagePath); err != nil {
		return nil, err
	} else {
		// Get loop device path
		loopPath := strings.Split(string(out), ":")[0]
		log.Infof("Loop device path is %s", loopPath)

		return &Loop{loopPath}, nil
	}
}

// Partition partitions the image
func (i *RawImage) Partition() error {
	// Create partition table
	log.Infof("Creating partition table on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath, "mklabel", "gpt"); err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath,
		"mkpart", "primary", "fat32", "1MiB", "512MiB"); err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath,
		"mkpart", "primary", "ext4", "512MiB", "100%"); err != nil {
		return err
	}

	// Set boot flag on EFI partition
	log.Infof("Setting boot flag on EFI partition on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath,
		"set", "1", "boot", "on"); err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath,
		"set", "1", "esp", "on"); err != nil {
		return err
	}

	// Set root partition label
	log.Infof("Setting root partition label on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath, "name", "2", "root"); err != nil {
		return err
	}

	// Set EFI partition label
	log.Infof("Setting EFI partition label on %s", i.imagePath)
	if err := util.RunCmd("parted", i.imagePath, "name", "1", "efi"); err != nil {
		return err
	}

	return nil
}

// Detach detaches the loop device
func (l *Loop) Detach() {
	// Get losetup path
	losetup, err := exec.LookPath("losetup")
	if err != nil {
		log.Errorf("Failed to get losetup path: %s", err)
	}

	// Detach loop device
	cmd := exec.Command(losetup, "-d", l.devicePath)
	err = cmd.Run()
	if err != nil {
		log.Errorf("Failed to detach loop device: %s", err)
	}
}

// Format formats the image partitions via the loop device
func (l *Loop) Format() error {
	// Get mkfs.vfat path
	mkfsVfat, err := exec.LookPath("mkfs.vfat")
	if err != nil {
		return err
	}

	// Get mkfs.ext4 path
	mkfsExt4, err := exec.LookPath("mkfs.ext4")
	if err != nil {
		return err
	}

	// Create FAT32 filesystem on EFI partition
	cmd := exec.Command(mkfsVfat, l.devicePath+"p1")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	cmd = exec.Command(mkfsExt4, l.devicePath+"p2")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
