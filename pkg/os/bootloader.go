package os

type OSBootloader interface {
	Install(device string) error
}
