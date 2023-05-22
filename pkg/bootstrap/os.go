package bootstrap

type OSFactory interface {
	NewOS(target string) OS
}

type OS interface {
	Bootstrap() error
	PostBootstrap() error
}
