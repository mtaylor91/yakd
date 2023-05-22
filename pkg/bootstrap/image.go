package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Image represents a disk image
type Image struct {
	Cleanup   bool
	Path      string
	SizeMB    int
	Overwrite bool
}

// Loop represents a loop device
type Loop struct {
	Path string
}

// NewImage initializes a new Image struct
func NewImage(path string, sizeMB int, overwrite bool) *Image {
	return &Image{
		Path:      path,
		SizeMB:    sizeMB,
		Overwrite: overwrite,
	}
}

// Create creates a new image file
func (i *Image) Create(mountpoint string, osFactory OSFactory) error {
	// Check if image exists
	if _, err := os.Stat(i.Path); err == nil && !i.Overwrite {
		return fmt.Errorf("image already exists")
	} else if err == nil {
		log.Infof("Removing existing image %s", i.Path)
		if err := os.Remove(i.Path); err != nil {
			return err
		}
	}

	// Create image
	log.Infof("Creating image %s", i.Path)
	if err := i.Alloc(); err != nil {
		return err
	}

	// Create partition table
	log.Infof("Creating partition table on %s", i.Path)
	if err := i.Partition(); err != nil {
		return err
	}

	// Attach image
	log.Infof("Attaching image %s", i.Path)
	loop, err := i.Attach()
	if err != nil {
		return err
	}

	if i.Cleanup {
		defer loop.Detach()
	}

	// Format image
	log.Infof("Formatting image %s on %s", i.Path, loop.Path)
	if err := loop.Format(); err != nil {
		return err
	}

	// Initialize parameters for bootstrap
	esp := loop.Path + "p1"
	root := loop.Path + "p2"
	os := osFactory.NewOS(mountpoint)
	disk := &Disk{loop.Path, esp, root, mountpoint, i.Cleanup, os}

	// Bootstrap image
	err = disk.Bootstrap()
	if err != nil {
		return err
	}

	return nil
}

// Alloc allocates a new image file
func (i *Image) Alloc() error {
	// Get dd path
	dd, err := exec.LookPath("dd")
	if err != nil {
		return err
	}

	// Create image using dd
	cmd := exec.Command(dd, "if=/dev/zero", "of="+i.Path, "bs=1M",
		"count=1", "seek="+strconv.Itoa(i.SizeMB-1))
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Attach attaches the image to a loop device
func (i *Image) Attach() (*Loop, error) {
	// Get losetup path
	losetup, err := exec.LookPath("losetup")
	if err != nil {
		return nil, err
	}

	// Attach image to loop device
	log.Infof("Attaching image %s to loop device", i.Path)
	cmd := exec.Command(losetup, "-P", "-f", i.Path)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// Get loop device info
	log.Infof("Getting loop device info for %s", i.Path)
	cmd = exec.Command(losetup, "-j", i.Path)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Get loop device path
	loopPath := strings.Split(string(out), ":")[0]
	log.Infof("Loop device path is %s", loopPath)

	return &Loop{loopPath}, nil
}

// Partition partitions the image
func (i *Image) Partition() error {
	// Get parted path
	parted, err := exec.LookPath("parted")
	if err != nil {
		return err
	}

	// Create partition table
	log.Infof("Creating partition table on %s", i.Path)
	cmd := exec.Command(parted, i.Path, "mklabel", "gpt")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "mkpart", "primary", "fat32", "1MiB", "512MiB")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "mkpart", "primary", "ext4", "512MiB", "100%")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Set boot flag on EFI partition
	log.Infof("Setting boot flag on EFI partition on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "set", "1", "boot", "on")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "set", "1", "esp", "on")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Set root partition label
	log.Infof("Setting root partition label on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "name", "2", "root")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Set EFI partition label
	log.Infof("Setting EFI partition label on %s", i.Path)
	cmd = exec.Command(parted, i.Path, "name", "1", "efi")
	err = cmd.Run()
	if err != nil {
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
	cmd := exec.Command(losetup, "-d", l.Path)
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
	cmd := exec.Command(mkfsVfat, l.Path+"p1")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create ext4 filesystem on root partition
	cmd = exec.Command(mkfsExt4, l.Path+"p2")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
