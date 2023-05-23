package os

type OS interface {
	Installer(target string) OSInstaller
}
