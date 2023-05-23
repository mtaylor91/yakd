package debian

import (
	"github.com/mtaylor91/yakd/pkg/bootstrap"
	"github.com/mtaylor91/yakd/pkg/os"
)

const (
	DefaultDebootstrap = "debootstrap"
	DefaultSuite       = "bullseye"
	DefaultMirror      = "http://deb.debian.org/debian"
)

var DefaultDebootstrapConfig = BootstrapConfig{
	DefaultSuite,
	DefaultMirror,
	bootstrap.DefaultTargetMountpoint,
	DefaultDebootstrap,
}

// BootstrapConfig represents the configuration for a bootstrap
type BootstrapConfig struct {
	Suite       string
	Mirror      string
	Target      string
	Debootstrap string
}

type Debian struct {
	Suite       string
	Mirror      string
	Debootstrap string
}

func (d *Debian) Installer(target string) os.OSInstaller {
	return NewBootstrapConfig(target)
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(target string) *BootstrapConfig {
	c := DefaultDebootstrapConfig
	c.Target = target
	return &c
}
