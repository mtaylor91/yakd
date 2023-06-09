package http

import (
	"context"
	"io"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Download struct {
	Source      string
	Destination string
}

// NewDownload initializes a new Download struct
func NewDownload(source, destination string) *Download {
	return &Download{source, destination}
}

// Download downloads a file from a URL to a destination
func (d *Download) Download(ctx context.Context) error {
	log.Infof("Downloading %s to %s", d.Source, d.Destination)

	// Send the request
	resp, err := Get(ctx, d.Source)
	if err != nil {
		return err
	}

	// Create the destination file
	f, err := os.Create(d.Destination)
	if err != nil {
		return err
	}

	// Copy the response body to the destination file
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return nil
}

// DownloadAndDearmorGPG downloads a file from a URL to a destination
// and performs a GPG dearmor
func (d *Download) DownloadAndDearmorGPG(ctx context.Context) error {
	log.Infof("Downloading %s to %s (and removing GPG armor)", d.Source, d.Destination)

	resp, err := Get(ctx, d.Source)

	// Locate gpg command
	gpg, err := exec.LookPath("gpg")
	if err != nil {
		return err
	}

	// Create the gpg dearmor command
	cmd := exec.Command(gpg, "--dearmor", "-o", d.Destination)
	cmd.Stdin = resp.Body
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
