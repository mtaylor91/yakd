package bootstrap

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Stage1 is the result of the first stage of the bootstrap process
type Stage1 struct {
	Source string
	Target string
}

// BuildArchive builds the stage1 archive
func (s *Stage1) BuildArchive() error {
	// Locate tar
	tarPath, err := exec.LookPath("tar")
	if err != nil {
		return err
	}

	// Create archive
	log.Infof("Creating stage1 archive at %s", s.Target)
	cmd := exec.Command(tarPath, "-C", s.Source, "-cJf", s.Target, ".")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
