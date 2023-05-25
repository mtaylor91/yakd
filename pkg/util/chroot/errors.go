package chroot

import "errors"

var ErrNoRoot = errors.New("chroot failed: no root specified")

var ErrNotSetup = errors.New("chroot failed: not setup")
