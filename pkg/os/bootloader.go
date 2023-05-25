package os

type OSBootloaderInstaller interface {
	Install() error
}
