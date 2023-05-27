package http

import (
	"bytes"
	"context"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Get performs a GET request.
func Get(ctx context.Context, url string) (*http.Response, error) {
	// Construct the request
	log.Debugf("GET %s", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the user agent
	req.Header.Set("User-Agent", "yakd")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetString performs a GET request and returns the response body as a string.
func GetString(ctx context.Context, url string) (string, error) {
	// Send the request
	resp, err := Get(ctx, url)
	if err != nil {
		return "", err
	}

	// Read the response body
	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		return "", err
	}

	return b.String(), nil
}
