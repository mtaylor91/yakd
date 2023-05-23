package os

type OS interface {
	Installer(target string) OSInstaller
	Bootloader(target string) OSBootloader
}
