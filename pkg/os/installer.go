package os

type OSInstaller interface {
	Bootstrap() error
	PostBootstrap() error
}
