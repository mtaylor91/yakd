package util

import log "github.com/sirupsen/logrus"

func PartitionDisk(name string) error {
	// Create partition table
	log.Infof("Creating partition table on %s", name)
	if err := RunCmd("parted", name, "mklabel", "gpt"); err != nil {
		return err
	}

	// Create EFI partition
	log.Infof("Creating EFI partition on %s", name)
	if err := RunCmd("parted", name,
		"mkpart", "primary", "fat32", "1MiB", "512MiB"); err != nil {
		return err
	}

	// Create root partition
	log.Infof("Creating root partition on %s", name)
	if err := RunCmd("parted", name,
		"mkpart", "primary", "ext4", "512MiB", "100%"); err != nil {
		return err
	}

	// Set boot flag on EFI partition
	log.Infof("Setting boot flag on EFI partition on %s", name)
	if err := RunCmd("parted", name,
		"set", "1", "boot", "on"); err != nil {
		return err
	}

	// Set esp flag on EFI partition
	log.Infof("Setting esp flag on EFI partition on %s", name)
	if err := RunCmd("parted", name,
		"set", "1", "esp", "on"); err != nil {
		return err
	}

	// Set root partition label
	log.Infof("Setting root partition label on %s", name)
	if err := RunCmd("parted", name, "name", "2", "root"); err != nil {
		return err
	}

	// Set EFI partition label
	log.Infof("Setting EFI partition label on %s", name)
	if err := RunCmd("parted", name, "name", "1", "efi"); err != nil {
		return err
	}

	return nil
}
