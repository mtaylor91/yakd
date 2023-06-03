package debian

import (
	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/system"
)

const (
	DefaultDebootstrap      = "debootstrap"
	DefaultSuite            = "bullseye"
	DefaultMirror           = "http://deb.debian.org/debian"
	DefaultTargetMountpoint = "build/mount"
)

var DebianDefault = &Debian{
	DefaultSuite,
	DefaultMirror,
	DefaultDebootstrap,
}

var DefaultDebootstrapConfig = BootstrapConfig{
	DefaultSuite,
	DefaultMirror,
	DefaultTargetMountpoint,
	DefaultDebootstrap,
}

type Debian struct {
	Suite       string
	Mirror      string
	Debootstrap string
}

func (d *Debian) DiskInstaller(
	device, target string, sys system.System,
) os.OSBootloaderInstaller {
	return NewGrubDiskInstaller(device, target, sys)
}

func (d *Debian) BootstrapInstaller(target string) os.OSBootstrapInstaller {
	return NewBootstrapConfig(target)
}

func (d *Debian) HybridISOSourceBuilder(fsDir, isoDir string) os.HybridISOSourceBuilder {
	return &HybridISOSourceBuilder{FSDir: fsDir, ISODir: isoDir}
}

func (d *Debian) HybridISOBuilder(isoDir, target string) os.HybridISOBuilder {
	return &HybridISOBuilder{ISODir: isoDir, Target: target}
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(target string) *BootstrapConfig {
	c := DefaultDebootstrapConfig
	c.Target = target
	return &c
}
