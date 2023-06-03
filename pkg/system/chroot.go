package system

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"

	"github.com/sirupsen/logrus"
)

type ChrootSystem struct {
	host System
	root string
	log  *logrus.Logger
}

func Chroot(host System, root string) *ChrootSystem {
	return &ChrootSystem{host, root, host.Logger()}
}

func (s *ChrootSystem) Context() context.Context {
	return s.host.Context()
}

func (s *ChrootSystem) Logger() *logrus.Logger {
	return s.log
}

func (s *ChrootSystem) RunCommand(cmd string, args ...string) error {
	outLog := s.log.WriterLevel(logrus.TraceLevel)
	return s.RunCommandWithIO(nil, outLog, cmd, args...)
}

func (s *ChrootSystem) RunCommandWithInput(
	input io.Reader, cmd string, args ...string,
) error {
	outLog := s.log.WriterLevel(logrus.TraceLevel)
	return s.RunCommandWithIO(input, outLog, cmd, args...)
}

func (s *ChrootSystem) RunCommandWithOutput(
	output io.Writer, cmd string, args ...string,
) error {
	return s.RunCommandWithIO(nil, output, cmd, args...)
}

func (s *ChrootSystem) RunCommandWithIO(
	input io.Reader, output io.Writer, cmd string, args ...string,
) error {
	// Prepend chroot command and root directory
	args = append([]string{s.root, cmd}, args...)
	// Run command
	return s.host.RunCommandWithIO(input, output, "chroot", args...)
}

func (s *ChrootSystem) SetContext(ctx context.Context) {
	s.host.SetContext(ctx)
}

func (s *ChrootSystem) SetLogger(log *logrus.Logger) {
	s.log = log
}

func (s *ChrootSystem) Setup() error {
	if err := s.MountMetadataFilesystems(); err != nil {
		return fmt.Errorf("failed to mount metadata filesystems: %w", err)
	}

	if err := s.CopyResolvConf(); err != nil {
		return fmt.Errorf("failed to copy resolv.conf: %w", err)
	}

	return nil
}

func (s *ChrootSystem) Teardown() {
	if err := s.UnmountMetadataFilesystems(); err != nil {
		s.log.Errorf("failed to unmount metadata filesystems: %v", err)
	}
}

// CopyResolvConf copies the host's resolv.conf to the bootstrap
func (s *ChrootSystem) CopyResolvConf() error {
	var hostResolvConf bytes.Buffer

	err := s.host.RunCommandWithOutput(&hostResolvConf, "cat", "/etc/resolv.conf")
	if err != nil {
		return fmt.Errorf("failed to read host resolv.conf: %w", err)
	}

	err = s.RunCommandWithInput(&hostResolvConf, "tee", "/etc/resolv.conf")
	if err != nil {
		return fmt.Errorf("failed to write bootstrap resolv.conf: %w", err)
	}
	return nil
}

// MountMetadataFilesystems creates the mountpoints for the bootstrap
func (s *ChrootSystem) MountMetadataFilesystems() error {
	commands := [][]string{
		[]string{"mount", "-t", "proc", "/proc", path.Join(s.root, "proc")},
		[]string{"mount", "--rbind", "/dev", path.Join(s.root, "dev")},
		[]string{"mount", "--make-rslave", path.Join(s.root, "dev")},
		[]string{"mount", "--rbind", "/sys", path.Join(s.root, "sys")},
		[]string{"mount", "--make-rslave", path.Join(s.root, "sys")},
		[]string{"mount", "--bind", "/run", path.Join(s.root, "run")},
		[]string{"mount", "--make-slave", path.Join(s.root, "run")},
	}

	for _, cmd := range commands {
		if err := s.host.RunCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

// UnmountMetadataFilesystems destroys the mountpoints for the bootstrap
func (s *ChrootSystem) UnmountMetadataFilesystems() error {
	commands := [][]string{
		[]string{"umount", "-R", path.Join(s.root, "proc")},
		[]string{"umount", "-R", path.Join(s.root, "dev", "pts")},
		[]string{"umount", "-R", path.Join(s.root, "dev", "shm")},
		[]string{"umount", "-R", path.Join(s.root, "dev")},
		[]string{"umount", "-R", path.Join(s.root, "sys")},
		[]string{"umount", "-R", path.Join(s.root, "run")},
	}

	var err error
	errCount := 0
	for _, cmd := range commands {
		if umountErr := s.host.RunCommand(cmd[0], cmd[1:]...); umountErr != nil {
			if umountErr != nil {
				errCount++
				if errCount == 1 {
					err = umountErr
				} else if errCount == 2 {
					desc := "multiple errors (see below)"
					err = fmt.Errorf("%s:\n%s\n%s",
						desc, err, umountErr)
				} else if errCount > 2 {
					err = fmt.Errorf("%s\n%s", err, umountErr)
				}
			}
		}
	}

	return err
}
