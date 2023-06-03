package system

import (
	"context"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type LocalSystem struct {
	ctx context.Context
	log *logrus.Logger
}

func (s *LocalSystem) Context() context.Context {
	return s.ctx
}

func (s *LocalSystem) Logger() *logrus.Logger {
	return s.log
}

func (s *LocalSystem) RunCommand(cmd string, args ...string) error {
	outLog := s.log.WriterLevel(logrus.TraceLevel)
	return s.RunCommandWithIO(nil, outLog, cmd, args...)
}

func (s *LocalSystem) RunCommandWithInput(
	input io.Reader, cmd string, args ...string,
) error {
	outLog := s.log.WriterLevel(logrus.TraceLevel)
	return s.RunCommandWithIO(input, outLog, cmd, args...)
}

func (s *LocalSystem) RunCommandWithOutput(
	output io.Writer, cmd string, args ...string,
) error {
	return s.RunCommandWithIO(nil, output, cmd, args...)
}

func (s *LocalSystem) RunCommandWithIO(
	input io.Reader, output io.Writer, cmd string, args ...string,
) error {
	s.log.Debugf("Running command: %s %v", cmd, args)
	outLog := s.log.WriterLevel(logrus.TraceLevel)
	c := exec.CommandContext(s.ctx, cmd, args...)
	c.Stdin = input
	c.Stdout = output
	c.Stderr = outLog
	return c.Run()
}

func (s *LocalSystem) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *LocalSystem) SetLogger(log *logrus.Logger) {
	s.log = log
}

func (s *LocalSystem) WithContext(ctx context.Context) *LocalSystem {
	return &LocalSystem{ctx, s.log}
}

func (s *LocalSystem) WithLogger(log *logrus.Logger) *LocalSystem {
	return &LocalSystem{s.ctx, log}
}
