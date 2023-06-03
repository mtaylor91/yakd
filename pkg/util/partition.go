package util

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/system"
)

func PartitionDisk(ctx context.Context, name string) error {
	sys := system.Local.WithContext(ctx)

	// Create partition table
	log.Infof("Creating partition table on %s", name)
	if err := sys.RunCommand("parted", "-s", name, "mklabel", "gpt"); err != nil {
		return err
	}

	// Create BIOS partition
	log.Infof("Creating BIOS partition on %s", name)
	if err := sys.RunCommand("parted", "-s", name,
		"mkpart", "primary", "1MiB", "10MiB"); err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", name)
	if err := sys.RunCommand("parted", "-s", name,
		"mkpart", "primary", "fat32", "10MiB", "512MiB"); err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", name)
	if err := sys.RunCommand("parted", "-s", name,
		"mkpart", "primary", "ext4", "512MiB", "100%"); err != nil {
		return err
	}

	// Set bios_grub flag on BIOS partition
	log.Infof("Setting bios_grub flag on BIOS partition on %s", name)
	if err := sys.RunCommand("parted", "-s", name,
		"set", "1", "bios_grub", "on"); err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", name)
	if err := sys.RunCommand("parted", "-s", name,
		"set", "2", "esp", "on"); err != nil {
		return err
	}

	return nil
}
