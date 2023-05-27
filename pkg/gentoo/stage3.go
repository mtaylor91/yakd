package gentoo

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util/http"
)

const (
	bouncer    = "https://bouncer.gentoo.org"
	releases   = bouncer + "/fetch/root/all/releases"
	autobuilds = releases + "/amd64/autobuilds"
	latest     = autobuilds + "/latest-stage3-amd64-systemd-mergedusr.txt"
)

func DownloadStage3(ctx context.Context, location string) error {
	// Identify latest stage3 tarball
	stage3url, err := identifyStage3(ctx)
	if err != nil {
		return err
	}

	log.Infof("Downloading %s to %s", stage3url, location)
	err = http.NewDownload(stage3url, location).Download(ctx)
	if err != nil {
		return err
	}

	return nil
}

func identifyStage3(ctx context.Context) (string, error) {
	// Identify latest stage3 tarball
	log.Info("Identifying latest gentoo stage3 tarball")
	resp, err := http.GetString(ctx, latest)
	if err != nil {
		return "", err
	}

	// Skip comments
	var line string
	for _, line = range strings.Split(resp, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		} else {
			break
		}
	}

	stage3info := strings.Split(line, " ")
	stage3url := autobuilds + "/" + stage3info[0]
	return stage3url, nil
}
