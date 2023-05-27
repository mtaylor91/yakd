package chroot

import (
	"context"
	"fmt"
	"io"
	"os"
)

func CopyResolvConf(ctx context.Context, root string) error {
	hostResolvConf, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return fmt.Errorf("failed to open host resolv.conf: %s", err)
	}

	defer hostResolvConf.Close()

	chrootResolvConf, err := os.OpenFile(
		root+"/etc/resolv.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open chroot resolv.conf: %s", err)
	}

	defer chrootResolvConf.Close()

	_, err = io.Copy(chrootResolvConf, hostResolvConf)
	if err != nil {
		return fmt.Errorf("failed to copy resolv.conf: %s", err)
	}

	return nil
}
