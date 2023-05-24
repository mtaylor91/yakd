package util

import log "github.com/sirupsen/logrus"

func PartitionDisk(name string) error {
	// Create partition table
	log.Infof("Creating partition table on %s", name)
	if err := RunCmd("parted", "-s", name, "mklabel", "msdos"); err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", name)
	if err := RunCmd("parted", "-s", name,
		"mkpart", "primary", "fat32", "1MiB", "512MiB"); err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", name)
	if err := RunCmd("parted", "-s", name,
		"mkpart", "primary", "ext4", "512MiB", "100%"); err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", name)
	if err := RunCmd("parted", "-s", name,
		"set", "1", "esp", "on"); err != nil {
		return err
	}

	return nil
}
