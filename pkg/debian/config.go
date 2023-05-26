package debian

import (
	"github.com/mtaylor91/yakd/pkg/os"
	"github.com/mtaylor91/yakd/pkg/util/bootstrap"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const (
	DefaultDebootstrap = "debootstrap"
	DefaultSuite       = "bullseye"
	DefaultMirror      = "http://deb.debian.org/debian"
)

var DebianDefault = &Debian{
	DefaultSuite,
	DefaultMirror,
	DefaultDebootstrap,
}

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

func (d *Debian) BootloaderInstaller(
	device, target string, exec executor.Executor,
) os.OSBootloaderInstaller {
	return NewGrubInstaller(device, target, exec)
}

func (d *Debian) BootstrapInstaller(target string) os.OSBootstrapInstaller {
	return NewBootstrapConfig(target)
}

// NewBootstrapConfig initializes a new BootstrapConfig struct
func NewBootstrapConfig(target string) *BootstrapConfig {
	c := DefaultDebootstrapConfig
	c.Target = target
	return &c
}
